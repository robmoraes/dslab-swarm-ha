# Roteiro da Temporada 2 — Alta Disponibilidade

Este documento organiza a progressão editorial e técnica da temporada. Cada
vídeo resolve uma camada de disponibilidade, evidencia o próximo ponto único de
falha e prepara o laboratório para a etapa seguinte.

## Objetivo da temporada

Partir da arquitetura concluída na temporada 1:

```text
Route 53 → EIP → Traefik único com HTTP-01 → WAF → API
                    │
                    └── manager único
```

E chegar a uma arquitetura tolerante à perda de um manager, de uma instância do
Traefik ou de um worker:

```text
Route 53
    ↓
ALB em múltiplas zonas + certificado ACM
    ↓
Target Group
    ├── manager-01 → Traefik
    ├── manager-02 → Traefik
    └── manager-03 → Traefik
                         ↓
                 WAF e API nos workers
```

## Fio narrativo

| Vídeo | Camada conquistada | Problema que permanece |
| ---: | --- | --- |
| 1 | Compreensão do quorum | O cluster ainda possui um único manager |
| 2 | Control plane com três managers | A entrada continua presa a um Traefik |
| 3 | Resiliência do Raft comprovada | O site cai quando a borda falha |
| 4 | Workloads isolados nos workers | Certificados e entrada continuam frágeis |
| 5 | Traefik replicado e problema isolado | DNS direto não remove falhas rapidamente |
| 6 | Entrada HA com ALB e ACM | A arquitetura ainda precisa ser provada |
| 7 | HA validada por falhas controladas | Limites da solução ficam documentados |

## Decisões técnicas da temporada

- Usar três managers distribuídos entre zonas de disponibilidade.
- Executar uma instância do Traefik em cada manager.
- Executar WAF e aplicações somente nos workers.
- Manter os managers ativos para que possam executar o Traefik; não colocá-los
  em `Drain`.
- Encerrar TLS no ALB quando ACM for introduzido.
- Remover ACME, `certresolver` e armazenamento de certificados do Traefik
  depois da migração para ACM.
- Usar health checks em todas as camadas relevantes.
- Preservar e validar o IP do cliente até a aplicação.
- Restaurar completamente o cluster entre dois testes de falha.
- Nunca publicar tokens, certificados, chaves privadas ou outros segredos.

## Vídeo 1 — Por que precisamos de três managers?

**Status:** pronto e gravado, aproximadamente 8 minutos.

**Formato:** teoria.

**Pergunta central:** por que dois managers não oferecem quorum tolerante a
falhas, mas três managers toleram a perda de um?

### Conteúdo

1. Apresentar Raft como o mecanismo que replica o estado do Swarm.
2. Explicar maioria e quorum.
3. Comparar clusters com um, dois, três, quatro e cinco managers.
4. Mostrar por que quantidades ímpares aproveitam melhor os nodes.
5. Diferenciar continuidade das tarefas existentes de capacidade de
   administrar e reagendar serviços.
6. Delimitar o resultado: três managers oferecem HA ao control plane, não
   necessariamente ao caminho completo da requisição.

### Resultado

O público entende o motivo técnico para construir um cluster com três
managers.

### Gancho

> Agora sabemos por que precisamos de três managers. No próximo vídeo vamos
> transformar nosso único node em um control plane tolerante a falhas.

## Vídeo 2 — De um manager para um cluster HA

**Status:** planejado.

**Formato:** prática guiada.

**Pergunta central:** como transformar o laboratório da temporada 1 em um
cluster com quorum?

### Estado inicial

- Uma instância EC2.
- Um manager.
- Um EIP.
- Um Traefik gerenciando certificados por HTTP-01.
- WAF e API executando no mesmo node.

### Roteiro

1. Revisar a arquitetura herdada da temporada 1.
2. Apresentar a topologia com três managers em zonas de disponibilidade
   diferentes.
3. Criar as duas novas instâncias EC2.
4. Configurar a comunicação interna do Swarm, restringindo as portas do cluster
   ao security group dos próprios nodes:
   - TCP 2377 para gerenciamento;
   - TCP e UDP 7946 para descoberta;
   - UDP 4789 para a rede overlay.
5. Obter o token de manager sem exibi-lo na gravação.
6. Adicionar os dois managers ao cluster.
7. Usar `docker node ls` para identificar `Leader` e `Reachable`.
8. Distribuir réplicas dos serviços entre os managers e comprovar a distribuição
   com `docker service ps`.
9. Manter apenas um Traefik, no manager original que possui o EIP.
10. Tornar esse posicionamento determinístico com uma label de node e uma
    placement constraint.

### Resultado

O Swarm termina com três managers e quorum, mas a entrada pública permanece
deliberadamente concentrada em um único manager.

### Gancho

> Temos quorum e réplicas distribuídas. Isso significa que podemos desligar
> qualquer máquina sem afetar o usuário?

## Vídeo 3 — O cluster sobrevive, mas o site cai

**Status:** planejado.

**Formato:** teste de desastre.

**Pergunta central:** ter quorum torna todo o sistema altamente disponível?

### Preparação da demonstração

Manter simultaneamente:

- um loop externo de requisições com horário, status e node respondente;
- `docker node ls`;
- `docker service ls`;
- `docker service ps` para os serviços principais.

### Roteiro

1. Registrar o estado saudável inicial.
2. Derrubar um manager que não executa o Traefik.
3. Observar que o quorum e as requisições permanecem disponíveis.
4. Restaurar o node e esperar que volte a `Ready` e `Reachable`.
5. Repetir o teste com o outro manager sem Traefik.
6. Restaurá-lo e aguardar novamente a convergência do cluster.
7. Derrubar o manager original, que acumula naquele momento os papéis de líder,
   EIP e Traefik.
8. Mostrar dois eventos independentes:
   - os managers restantes elegem um novo líder e mantêm o Swarm operacional;
   - as requisições externas param porque o único ponto de entrada desapareceu.
9. Confirmar que as tarefas continuam executando mesmo sem acesso público.

### Mensagem principal

> O cluster sobreviveu à perda do líder, mas o serviço público não sobreviveu à
> perda da borda.

A indisponibilidade deve ser atribuída ao EIP e ao Traefik únicos, não à eleição
do Raft.

### Resultado

HA do control plane comprovada; ausência de HA de ponta a ponta demonstrada.

### Gancho

> Antes de corrigir a entrada, precisamos separar quem administra o cluster de
> quem executa nossas aplicações.

## Vídeo 4 — Managers, workers e separação dos workloads

**Status:** planejado.

**Formato:** teoria curta seguida de prática.

**Pergunta central:** quais responsabilidades pertencem aos managers e quais
pertencem aos workers?

### Conteúdo conceitual

- Managers participam do Raft, armazenam o estado e fazem a orquestração.
- Workers executam as tarefas dos serviços.
- Managers também recebem tarefas por padrão.
- O Traefik permanece nos managers porque consulta a API do Swarm disponível
  nesses nodes.
- Separar workloads de aplicação não significa colocar os managers em
  `Drain`.

### Roteiro prático

1. Criar três workers distribuídos entre as zonas de disponibilidade.
2. Adicioná-los com o token de worker, sem expor o token.
3. Identificar claramente os dois grupos em `docker node ls`.
4. Aplicar placement constraints para que API e WAF executem somente em
   `node.role == worker`.
5. Manter o Traefik restrito a `node.role == manager`.
6. Redistribuir as réplicas e mostrar o resultado com `docker service ps`.
7. Explicar que os managers executam o control plane e o serviço de plataforma
   responsável pela entrada, enquanto os workloads da aplicação ficam nos
   workers.

### Resultado

```text
Managers: Raft + API do Swarm + Traefik
Workers:  WAF + aplicações
```

Os workloads estão separados, mas o Traefik e os certificados ainda formam um
ponto único de falha.

### Gancho

> Se o Traefik pode ser replicado, por que não simplesmente executar uma
> instância em cada manager?

## Vídeo 5 — Por que replicar o Traefik ainda não resolve tudo?

**Status:** planejado.

**Formato:** retrospectiva técnica e experimento controlado.

**Pergunta central:** quais problemas aparecem quando três instâncias do
Traefik tentam administrar a mesma entrada e os mesmos certificados?

### Histórico das tentativas

1. **DNS-01 com Route 53 em cada Traefik**
   - cada instância dispara sua própria emissão ou renovação;
   - surgem emissões duplicadas;
   - os limites da autoridade certificadora são consumidos;
   - o domínio pode sofrer bloqueio temporário por rate limit.
2. **HTTP-01 com o arquivo ACME em EFS**
   - o arquivo passa a ser compartilhado;
   - o filesystem compartilhado não fornece a coordenação exigida entre os
     processos;
   - gravações concorrentes e renovação continuam sendo riscos.
3. **Traefik único**
   - emissão e renovação ficam previsíveis;
   - o componente permanece como ponto único de falha.

### Experimento controlado

1. Desativar ACME durante o experimento.
2. Carregar o mesmo certificado estático nas três instâncias do Traefik usando
   secrets e configuração dinâmica.
3. Nunca armazenar a chave privada no repositório.
4. Executar um Traefik em cada manager.
5. Associar um endereço público a cada manager e configurar os endereços no
   Route 53.
6. Validar cada destino individualmente preservando Host e SNI:

   ```bash
   curl --resolve whoami.example.com:443:IP_DO_MANAGER \
     https://whoami.example.com
   ```

7. Realizar requisições pela resolução DNS normal.
8. Derrubar uma instância do Traefik.
9. Demonstrar que:
   - o certificado continua válido nos outros targets;
   - duas instâncias continuam saudáveis;
   - uma parcela dos clientes ainda tenta ou mantém em cache o IP indisponível.

Não prometer exatamente um terço de falhas. Route 53 retorna endereços, mas a
escolha, a reutilização e o cache também dependem do resolver e do cliente.

### Requisitos extraídos do problema

A solução final precisa oferecer:

- um endpoint público estável;
- balanceamento ativo;
- remoção automática de targets indisponíveis;
- certificado independente das réplicas do Traefik;
- Traefik sem estado ACME compartilhado;
- preservação segura do IP do cliente.

### Resultado

Os certificados deixam de ser a variável do teste e a limitação do DNS direto
fica isolada.

### Gancho

> Precisamos colocar uma borda altamente disponível e consciente da saúde dos
> targets antes do nosso cluster.

## Vídeo 6 — ALB, Target Groups e ACM

**Status:** planejado.

**Formato:** arquitetura seguida de implementação.

**Pergunta central:** como remover da aplicação a responsabilidade por
certificados e failover da borda?

### Arquitetura desejada

```text
Internet
   ↓
Route 53 alias
   ↓
ALB em múltiplas zonas
   ├── listener HTTP  :80  → redirect
   └── listener HTTPS :443 → certificado ACM
                              ↓
                         Target Group HTTP
                         ├── manager-01:porta-traefik
                         ├── manager-02:porta-traefik
                         └── manager-03:porta-traefik
                                      ↓
                              Traefik → WAF → API
```

### Roteiro

1. Solicitar o certificado no ACM.
2. Validar o domínio por DNS.
3. Criar o ALB nas zonas que contêm os targets.
4. Criar um security group público para os listeners do ALB.
5. Permitir que os managers recebam tráfego na porta do Traefik somente a partir
   do security group do ALB.
6. Criar o Target Group com os três managers.
7. Criar um endpoint ou router dedicado ao health check que não dependa do
   `Host` público enviado pelo cliente.
8. Configurar health check do Target Group com intervalo de 30 segundos e
   timeout de 5 segundos.
9. Criar o listener HTTP com redirecionamento para HTTPS.
10. Criar o listener HTTPS usando o certificado do ACM.
11. Executar o Traefik em modo global, restrito aos managers e com publicação de
    porta em modo host.
12. Remover do Traefik:
    - HTTP-01;
    - `certresolver`;
    - `acme.json`;
    - volume de certificados;
    - router HTTPS interno que deixou de existir.
13. Manter um router HTTP interno, pois o TLS termina no ALB.
14. Configurar health checks do Traefik para o serviço realmente exposto, usando
    endpoint leve, intervalo de 30 segundos e timeout entre 3 e 5 segundos.
15. Alterar o Route 53 para apontar para o ALB.
16. Validar a cadeia de encaminhamento:
    - ALB adiciona os headers `X-Forwarded-*`;
    - Traefik confia somente nos proxies e redes previstos;
    - WAF propaga os headers;
    - aplicação interpreta e exibe o IP correto do cliente.
17. Garantir que comunicações entre stacks usem nomes Swarm qualificados no
    formato `<stack>_<serviço>`.

### Resultado

- Certificado gerenciado pelo ACM.
- Três instâncias stateless do Traefik em relação ao TLS público.
- ALB removendo targets incapazes de receber tráfego.
- Managers sem exposição pública direta da porta da aplicação.
- IP do cliente preservado e validado na própria API.

### Gancho

> Construímos a arquitetura. Agora precisamos provar que ela realmente
> sobrevive às falhas que derrubavam o laboratório anterior.

> Se a implementação ficar longa, este vídeo pode ser dividido em “Desenhando a
> borda HA” e “Implementando ALB + ACM”, sem alterar a ordem conceitual.

## Vídeo 7 — Chaos Day: agora temos HA de verdade?

**Status:** planejado.

**Formato:** validação de arquitetura.

**Pergunta central:** a aplicação continua útil quando cada componente falha
isoladamente?

### Preparação

Criar um loop externo que registre:

- horário;
- status HTTP;
- latência;
- manager/Traefik utilizado;
- worker e réplica da aplicação;
- IP do cliente observado pela aplicação.

Acompanhar paralelamente:

- `docker node ls`;
- `docker service ls`;
- `docker service ps`;
- estado dos targets no ALB;
- logs do Traefik, WAF e aplicação.

### Matriz de falhas

| Falha | Resultado esperado |
| --- | --- |
| Container Traefik | Swarm recria a tarefa e o ALB usa os demais targets |
| Manager follower | Quorum permanece e o ALB remove o target indisponível |
| Manager leader | Novo líder é eleito e o tráfego usa os outros Traefiks |
| Worker | Réplicas restantes continuam atendendo e tarefas são reconciliadas |
| Réplica do WAF ou API | Balanceamento interno evita enviar tráfego à réplica incapaz |

### Roteiro

1. Registrar a linha de base saudável.
2. Executar uma falha por vez.
3. Medir erros transitórios e tempo de detecção.
4. Observar a mudança do target para `unhealthy`.
5. Confirmar que novas requisições seguem por targets saudáveis.
6. Restaurar o componente.
7. Esperar a recuperação completa antes do próximo teste.
8. Confirmar que o certificado continua sendo entregue pelo ALB.
9. Confirmar que nenhuma alteração manual de DNS é necessária.
10. Validar novamente o IP do cliente na aplicação.

Não prometer zero erros durante uma falha abrupta. Entre a queda e a detecção do
health check podem existir falhas transitórias. O objetivo é medir e explicar
essa janela.

### Critérios de aceite

- A perda de um manager mantém o quorum.
- A perda do líder produz uma nova eleição.
- A perda de um Traefik não remove o endereço público.
- A perda de um worker mantém a aplicação atendendo.
- Targets incapazes são retirados de circulação.
- O certificado não depende de uma réplica do Traefik.
- O IP do cliente chega corretamente à aplicação.
- A recuperação não exige mudança manual de DNS ou certificado.

### Limites que devem ser declarados

- Duas falhas simultâneas de managers eliminam o quorum de um cluster com três.
- Uma falha regional não é coberta.
- O laboratório valida workloads stateless; bancos e outros serviços stateful
  exigem estratégias próprias.
- Alta disponibilidade reduz impacto e tempo de recuperação, mas não significa
  ausência absoluta de erros.

### Encerramento

> Alta disponibilidade não é a ausência de falhas. É a capacidade de
> detectá-las, isolá-las e continuar entregando o serviço dentro de limites
> conhecidos.

## Checklist comum de gravação

Antes de cada teste de falha:

- Confirmar que todos os nodes estão `Ready`.
- Confirmar que os managers estão `Leader` ou `Reachable`.
- Confirmar que todos os targets do ALB estão saudáveis.
- Registrar a quantidade e a distribuição das réplicas.
- Iniciar o loop de requisições antes da falha.
- Exibir relógio ou timestamp para medir convergência.
- Ocultar tokens, IDs sensíveis, certificados e chaves.
- Restaurar completamente o ambiente antes do cenário seguinte.

## Referências técnicas

- [Administração e quorum do Docker Swarm](https://docs.docker.com/engine/swarm/admin_guide/)
- [Managers, workers e disponibilidade dos nodes](https://docs.docker.com/engine/swarm/how-swarm-mode-works/nodes/)
- [Provider Docker Swarm do Traefik](https://doc.traefik.io/traefik/reference/install-configuration/providers/swarm/)
- [Certificate resolvers do Traefik](https://doc.traefik.io/traefik/v3.6/reference/install-configuration/tls/certificate-resolvers/overview/)
- [Desafios ACME do Let's Encrypt](https://letsencrypt.org/docs/challenge-types/)
- [Simple routing do Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-simple.html)
- [Health checks de Target Groups do ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/target-group-health-checks.html)
- [Certificados HTTPS no ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/https-listener-certificates.html)
- [Headers encaminhados pelo ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html)
