# Roteiro da Temporada 2 — Alta Disponibilidade

Este documento organiza a progressão editorial e técnica da temporada. Cada
vídeo resolve uma camada de disponibilidade, evidencia o próximo ponto único de
falha e prepara o laboratório para a etapa seguinte.

## Objetivo da temporada

Começar estabelecendo uma identidade própria para o laboratório:

- `dslab.dev.br` como zona DNS pública e autoritativa no Route 53;
- `swarm.dslab.dev.br` como domínio raiz do cluster;
- `*.swarm.dslab.dev.br` como padrão para os serviços.

Depois da delegação e da migração do nome, partir da arquitetura concluída na
temporada 1:

```text
*.swarm.dslab.dev.br → Route 53 → EIP → Traefik único com HTTP-01 → WAF → API
                    │
                    └── manager único
```

E chegar a uma arquitetura tolerante à perda de um manager, de uma instância do
Traefik ou de um worker:

```text
swarm.dslab.dev.br e *.swarm.dslab.dev.br
    ↓
Route 53 alias
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

| Vídeo | Camada conquistada                  | Problema que permanece                  |
| ----: | ----------------------------------- | --------------------------------------- |
|     1 | Autoridade DNS e nomes canônicos    | O cluster ainda possui um único manager |
|     2 | Compreensão do quorum               | O cluster ainda possui um único manager |
|     3 | Quorum e Traefik global em HTTP     | DNS sem health checks; HTTPS pendente   |
|     4 | Workloads isolados nos workers      | A entrada ainda depende de DNS direto   |
|     5 | Limites dos certificados discutidos | Falta uma borda com TLS e health checks |
|     6 | Entrada HA com ALB e ACM            | A arquitetura ainda precisa ser provada |
|     7 | HA validada por falhas controladas  | Limites da solução ficam documentados   |

## Decisões técnicas da temporada

- Criar a public hosted zone `dslab.dev.br` no Route 53 e delegar o domínio
  registrado no Registro.br para os quatro name servers fornecidos pela AWS.
- Usar `swarm.dslab.dev.br` como nome estável da entrada do cluster.
- Usar nomes no formato `<serviço>.swarm.dslab.dev.br` para os serviços.
- Manter `*.swarm.dslab.dev.br` apontando para o domínio raiz do cluster, para
  que a troca futura de EIP para ALB aconteça em um único destino.
- Incluir `swarm.dslab.dev.br` e `*.swarm.dslab.dev.br` no certificado do
  ACM, pois o wildcard não protege o próprio domínio raiz do cluster.
- Usar três managers distribuídos entre zonas de disponibilidade.
- Executar uma instância do Traefik em cada manager.
- Executar WAF e aplicações somente nos workers.
- Manter os managers ativos para que possam executar o Traefik; não colocá-los
  em `Drain`.
- Encerrar TLS no ALB quando ACM for introduzido.
- Desativar HTTPS, ACME, `certresolver` e o volume de certificados no experimento
  HTTP do vídeo 3. Reintroduzir HTTPS público com ALB e ACM, sem devolver a
  emissão de certificados às réplicas do Traefik.
- Usar health checks em todas as camadas relevantes.
- Preservar e validar o IP do cliente até a aplicação.
- Restaurar completamente o cluster entre dois testes de falha.
- Nunca publicar tokens, certificados, chaves privadas ou outros segredos.

## Vídeo 1 — Do Registro.br ao Route 53

**Status:** planejado.

**Formato:** fundamentos de DNS seguidos de prática.

**Pergunta central:** como delegar o domínio registrado para o Route 53 e criar
uma identidade estável para o cluster e seus serviços?

### Convenção de nomes

```text
dslab.dev.br                         zona pública
└── swarm.dslab.dev.br               entrada raiz do cluster
    ├── whoami.swarm.dslab.dev.br    serviço
    ├── traefik.swarm.dslab.dev.br   serviço
    └── *.swarm.dslab.dev.br         wildcard para novos serviços
```

O wildcard DNS organiza a resolução dos serviços, mas não deve ser confundido
com a cobertura do certificado. Um certificado para
`*.swarm.dslab.dev.br` não cobre `swarm.dslab.dev.br`; o ACM deverá receber
os dois nomes mais adiante.

### Roteiro

1. Mostrar o domínio `dslab.dev.br` no Registro.br.
2. Explicar a diferença entre registrar um domínio e hospedar sua zona DNS.
3. Criar primeiro uma public hosted zone chamada `dslab.dev.br` no Route 53.
4. Identificar os quatro name servers atribuídos pela AWS.
5. Verificar se existe uma configuração DNSSEC ou um registro DS anterior e
   evitar manter uma cadeia de confiança apontando para chaves antigas.
6. No Registro.br, substituir os servidores DNS atuais pelos quatro name
   servers da hosted zone.
7. Explicar propagação e cache sem prometer ativação instantânea.
8. Validar a delegação:

   ```bash
   dig NS dslab.dev.br +short
   dig SOA dslab.dev.br +short
   dig +trace dslab.dev.br
   ```

9. Criar um registro A para `swarm.dslab.dev.br` apontando inicialmente para
   o EIP do manager existente.
10. Criar `*.swarm.dslab.dev.br` como CNAME de
    `swarm.dslab.dev.br`.
11. Atualizar o router do serviço para
    `whoami.swarm.dslab.dev.br`.
12. Emitir o certificado específico necessário nessa fase usando o Traefik
    único e HTTP-01.
13. Validar o domínio raiz, o wildcard e o serviço:

    ```bash
    dig A swarm.dslab.dev.br +short
    dig whoami.swarm.dslab.dev.br
    curl -I https://whoami.swarm.dslab.dev.br
    ```

### Resultado

- Registro.br continua como registrador.
- Route 53 passa a ser autoritativo por `dslab.dev.br`.
- O cluster ganha um domínio próprio e independente de nomes pessoais.
- Serviços passam a seguir um padrão previsível.
- O wildcard acompanha automaticamente a futura troca do destino raiz para o
  ALB.

### Gancho

> Já sabemos como os usuários encontrarão o cluster. Agora precisamos entender
> quantos managers mantêm esse cluster disponível quando uma máquina falha.

## Vídeo 2 — Por que precisamos de três managers?

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

## Vídeo 3 — Três managers: quorum, falhas e entrada via DNS

**Status:** gravado, aproximadamente 45 minutos.

**Formato:** prática guiada e testes de falha, acompanhados em terminal dividido.

**Pergunta central:** como transformar o laboratório da temporada 1 em um
cluster com quorum e descobrir os limites da disponibilidade da entrada?

### Estado inicial

- Uma instância EC2.
- Um manager.
- Um EIP.
- Um Traefik gerenciando certificados por HTTP-01.
- WAF e API executando no mesmo node.

### Acompanhamento em terminal dividido

Manter a mesma visão durante a maior parte da demonstração:

- um painel por manager para acompanhar os containers e seus recursos com
  `docker stats`;
- um painel com `docker node ls` para observar `Ready`, `Leader`, `Reachable`
  e `Unreachable`;
- um painel externo com `scripts/requests-whoami.sh`, inicialmente em HTTPS e
  depois em HTTP, mostrando as respostas da API ou `status=unavailable`.

O loop aguarda um segundo depois de cada chamada. Na fase HTTP, os limites de
conexão e transferência impedem que uma tentativa sem resposta prenda o
terminal por muito tempo.

### Construção do cluster

1. Revisar a arquitetura herdada da temporada 1.
2. Apresentar a topologia com três managers em zonas de disponibilidade
   diferentes.
3. Criar as duas novas instâncias EC2.
4. Configurar a comunicação interna do Swarm, restringindo as portas do cluster
   ao security group dos próprios nodes:
   - TCP 2377 para gerenciamento;
   - TCP e UDP 7946 para descoberta;
   - UDP 4789 para a rede overlay.
5. Obter o token de manager.
6. Adicionar os dois managers ao cluster.
7. Usar `docker node ls` para identificar `Leader` e `Reachable`.
8. Distribuir réplicas dos serviços entre os managers e comprovar a distribuição
   com `docker service ps`.
9. Manter apenas um Traefik, no manager original que possui o EIP.
10. Tornar esse posicionamento determinístico com uma label de node e uma
    placement constraint.

### Testes com a entrada única

1. Registrar o estado saudável inicial.
2. Derrubar o Traefik único enquanto seu node continua ligado e mostrar a
   interrupção das requisições externas, sem confundir a queda do proxy com
   perda de quorum.
3. Restaurar o Traefik e confirmar que o acesso HTTPS voltou.
4. Derrubar um manager que não executa o Traefik.
5. Observar que o quorum e as requisições permanecem disponíveis.
6. Restaurar o node e esperar que volte a `Ready` e `Reachable`.
7. Repetir o teste com o outro manager sem Traefik.
8. Restaurá-lo e aguardar novamente a convergência do cluster.
9. Derrubar o manager original, que concentra o EIP e o único Traefik. Quando
   esse manager também é o líder, acompanhar a nova eleição.
10. Mostrar dois eventos independentes:
    - os managers restantes elegem um novo líder e mantêm o Swarm operacional;
    - as requisições externas param porque o único ponto de entrada desapareceu.
11. Confirmar que as tarefas continuam executando mesmo sem acesso público.

### Primeira conclusão

> O cluster sobreviveu à perda do líder, mas o serviço público não sobreviveu à
> perda da borda.

A indisponibilidade deve ser atribuída ao EIP e ao Traefik únicos, não à eleição
do Raft.

### Evolução experimental da entrada

1. Restaurar os managers e confirmar novamente o quorum e os serviços.
2. Desativar HTTPS, o redirecionamento HTTP para HTTPS e ACME no Traefik.
3. Alterar o router do WAF para o entrypoint HTTP `web`, sem TLS ou
   `certresolver`, e mudar o loop de requisições para HTTP.
4. Recriar o serviço Traefik em modo global, restrito aos managers, publicando
   a porta 80 em `mode: host`. A troca de `replicated` para `global` exige
   recriar o serviço; o Swarm não permite atualizar o modo do serviço existente.
5. Configurar o registro A da entrada no Route 53 com os três EIPs dos
   managers, mantendo o padrão de nomes dos serviços.
6. Confirmar uma tarefa do Traefik em cada manager e as réplicas de WAF e API
   distribuídas pelo cluster.
7. Repetir os testes de falha e acompanhar simultaneamente os containers, o
   estado dos managers e o loop externo.
8. Observar o que acontece de fato, incluindo respostas, demora e
   `status=unavailable`, sem pressupor uma proporção fixa de falhas.

O registro A simples não verifica a saúde dos nodes e continua anunciando um
EIP indisponível. O cliente pode tentar outro endereço e ainda concluir uma
chamada; isso não significa que o DNS tenha retirado o destino que falhou.

Na captura de apoio da demonstração, o terminal registra tarefas do Traefik nos
três nodes, um manager `Unreachable`, outro como `Leader` e chamadas externas
com `status=unavailable`. Essa composição permite acompanhar lado a lado o
estado do cluster e o efeito percebido pelo cliente.

### Resultado

- Três managers com quorum e eleição de líder observados durante as falhas.
- Diferença entre sobrevivência do cluster e disponibilidade pública
  demonstrada pela queda do Traefik único e pela queda do manager de entrada.
- Traefik global e distribuição da entrada pelos três EIPs experimentados em
  HTTP, sem a complexidade de certificados ACME concorrentes.
- Limites do DNS direto evidenciados; a entrada ainda não tem remoção ativa de
  targets indisponíveis e o HTTPS público precisa voltar.

### Gancho

> Replicamos a entrada, mas o DNS continua anunciando destinos mesmo quando
> falham. A próxima camada da evolução precisa saber quais targets estão
> saudáveis e devolver HTTPS sem certificados concorrentes: ALB e ACM.

## Vídeo 4 — Managers, workers e separação dos workloads

**Status:** planejado.

**Formato:** teoria curta seguida de prática.

**Pergunta central:** como tirar WAF e API dos managers mantendo o Traefik
global e verificando o caminho da requisição durante a migração?

### Ponto de partida

- Três managers com quorum; cada um executa uma tarefa do Traefik em HTTP.
- Route 53 anuncia os três EIPs dos managers, sem health checks.
- WAF com três réplicas e API com cinco, ainda elegíveis para executar nos
  managers.
- O loop externo de requisições mostra o serviço antes, durante e depois da
  migração. Não alterar DNS ou reintroduzir HTTPS neste vídeo.

### Conteúdo conceitual

- Managers mantêm o estado via Raft, agendam serviços e expõem a API do Swarm;
  também executam tarefas por padrão. Workers executam tarefas, mas não votam
  no quorum nem administram o cluster.
- O Traefik continua nos managers para consultar a API do Swarm pelo socket
  local. WAF e API podem executar nos workers e conversar com ele pelas redes
  overlay existentes.
- A separação é uma regra de placement dos serviços, não uma mudança de papel
  dos managers. Colocar um manager em `Drain` retiraria também seu Traefik
  global; os três managers devem permanecer `Active`.
- `node.role == worker` limita onde a tarefa pode executar, mas não garante
  uma réplica por worker nem uma por zona de disponibilidade. Mostrar a
  distribuição observada, sem tratá-la como garantia automática.

### Preparação

1. Registrar `docker node ls`, `docker service ls` e o posicionamento atual de
   `traefik_traefik`, `whoami_waf` e `whoami_api` com
   `docker service ps <servico>`.
2. Iniciar `scripts/requests-whoami.sh` em um terminal externo para comparar
   as respostas durante a transição. O DNS direto ainda pode produzir falhas
   independentes da migração dos workloads.
3. Revisar as três zonas de disponibilidade e o security group existente. Os
   workers precisam falar com os managers pelo IP privado na porta TCP 2377 e
   entre si e com os managers em TCP/UDP 7946 e UDP 4789. Não publicar essas
   portas para a Internet.

### Configuração dos workers

Provisionar as três instâncias pelo
[guia dos workers EC2](../../../provisioning/workers/ec2.md):

| Item | Configuração |
| --- | --- |
| Quantidade e zonas | Três workers, um por zona: `us-east-1a`, `us-east-1b` e `us-east-1c` |
| Sistema e arquitetura | Ubuntu Server 26.04 LTS, `x86_64` |
| AMI em `us-east-1` | `ami-0b6d9d3d33ba97d99`, da Canonical (`099720109477`) |
| Tipo e compra | `t3.small`, 2 vCPUs e 2 GiB de RAM, On-Demand |
| Disco | 30 GiB gp3, com exclusão ao encerrar a instância |
| Rede | VPC e subnets existentes; security group `DSLab Swarm SG` |
| Acesso | Importar `.ssh/swarm.dslab.pub` como `dslab-swarm-current`; habilitar IPv4 público automático |
| Metadados | IMDS habilitado, V2 only e response hop limit `2` |

Usar o [bootstrap do worker](../../../scripts/bootstrap-worker.sh) para
instalar o Docker em cada EC2. O guia inclui os comandos de SSH, a entrada com
o token de worker e a conferência dos nodes. Os IPs privados são usados na
comunicação do cluster; os workers não recebem EIP nem registro no DNS público.

### Roteiro prático

1. Criar as três EC2 workers conforme a configuração acima e o guia de
   provisionamento. Garantir conectividade de saída para instalar o Docker e
   baixar as imagens.
2. Executar o bootstrap do worker. Obter no manager o comando de
   entrada com `docker swarm join-token worker` e executá-lo nos novos nodes,
   usando o IP privado de um manager. Não exibir o token na gravação nem
   confundi-lo com o token de manager.
3. Confirmar seis nodes `Ready` e `Active` com `docker node ls`: três managers
   (`Leader`/`Reachable`) e três workers (sem `MANAGER STATUS`). A entrada dos
   workers, por si só, não deve ser apresentada como redistribuição imediata
   das tarefas existentes.
4. Editar `stacks/whoami.yml` durante a gravação: adicionar a constraint
   `node.role == worker` em `deploy.placement.constraints` tanto para `api`
   quanto para `waf`. Manter as cinco e três réplicas e as redes `appnet` e
   `traefik-public` como estão.
5. Confirmar que `stacks/traefik.yml` mantém `mode: global` e
   `node.role == manager`. Não colocar os managers em `Drain`.
6. Aplicar a stack atualizada com `docker stack deploy -c stacks/whoami.yml
   --detach=false whoami`. A mudança das constraints atualiza os serviços e
   move suas tarefas para os workers; acompanhar a convergência sem prometer
   zero requisições perdidas durante a atualização.
7. Verificar `docker service ps traefik_traefik`, `docker service ps
   whoami_waf` e `docker service ps whoami_api`: Traefik nos três managers,
   todas as tarefas atuais do WAF e da API nos workers. Confirmar as contagens
   `3/3` e `5/5` em `docker service ls` e acompanhar o loop externo.
8. Colocar temporariamente um worker em `Drain` e observar as tarefas da API
   e do WAF serem reagendadas nos outros workers, nunca nos managers. Voltar
   esse worker para `Active` e esperar a recuperação antes de encerrar. O
   Swarm não precisa redistribuir de volta as tarefas só porque o node voltou.

Se alguma tarefa ficar `Pending` ou falhar ao iniciar, examinar
`docker service ps --no-trunc <servico>` antes de atribuir o problema ao placement:
capacidade dos workers, acesso ao registry e comunicação overlay também entram
no caminho.

### Resultado

```text
Managers: Raft + API do Swarm + Traefik
Workers:  WAF + aplicações
```

Os workloads estão separados e o Traefik continua distribuído pelos managers.
Isso protege o control plane da carga da aplicação, mas não resolve sozinho a
distribuição das réplicas entre zonas, a saúde dos destinos anunciados no DNS
nem o HTTPS público.

### Gancho

> Replicamos o proxy em HTTP. Como devolver HTTPS sem fazer cada instância
> disputar a emissão e a renovação dos certificados?

## Vídeo 5 — Por que replicar o Traefik ainda não resolve tudo?

**Status:** planejado.

**Formato:** retrospectiva técnica e experimento controlado.

O vídeo 3 já incluiu o experimento com Traefik global em HTTP e os três EIPs no
DNS. Este bloco fica como aprofundamento dos certificados; o experimento TLS
abaixo é complementar e opcional, não uma repetição necessária daquela etapa.

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

### Experimento complementar com HTTPS (opcional)

1. Desativar ACME durante o experimento.
2. Carregar o mesmo certificado estático nas três instâncias do Traefik usando
   secrets e configuração dinâmica.
3. Nunca armazenar a chave privada no repositório.
4. Executar um Traefik em cada manager.
5. Associar um endereço público a cada manager e configurar os endereços no
   Route 53.
6. Validar cada destino individualmente preservando Host e SNI:

   ```bash
   curl --resolve whoami.swarm.dslab.dev.br:443:IP_DO_MANAGER \
     https://whoami.swarm.dslab.dev.br
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

1. Solicitar no ACM um certificado com os dois nomes:
   - `swarm.dslab.dev.br`;
   - `*.swarm.dslab.dev.br`.

   Explicar que o wildcard cobre os serviços, mas não o domínio raiz do
   cluster.

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
11. Confirmar que o Traefik continua em modo global, restrito aos managers e
    com publicação de porta em modo host, como no experimento do vídeo 3.
12. Verificar que o Traefik permanece sem os itens retirados na fase HTTP:
    - HTTP-01;
    - `certresolver`;
    - `acme.json`;
    - volume de certificados;
    - router HTTPS interno que deixou de existir.
13. Manter um router HTTP interno, pois o TLS termina no ALB.
14. Configurar health checks do Traefik para o serviço realmente exposto, usando
    endpoint leve, intervalo de 30 segundos e timeout entre 3 e 5 segundos.
15. Trocar o registro A de `swarm.dslab.dev.br` por um Alias para o ALB.
    Manter o wildcard apontando para esse nome raiz.
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

| Falha                 | Resultado esperado                                                 |
| --------------------- | ------------------------------------------------------------------ |
| Container Traefik     | Swarm recria a tarefa e o ALB usa os demais targets                |
| Manager follower      | Quorum permanece e o ALB remove o target indisponível              |
| Manager leader        | Novo líder é eleito e o tráfego usa os outros Traefiks             |
| Worker                | Réplicas restantes continuam atendendo e tarefas são reconciliadas |
| Réplica do WAF ou API | Balanceamento interno evita enviar tráfego à réplica incapaz       |

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

- [Alteração de servidores DNS no Registro.br](https://registro.br/ajuda/gerenciamento-de-conta)
- [Delegação de um domínio externo para o Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/dns-configuring-new-domain.html)
- [Cobertura de certificados wildcard no ACM](https://docs.aws.amazon.com/acm/latest/userguide/acm-certificate-characteristics.html)
- [Administração e quorum do Docker Swarm](https://docs.docker.com/engine/swarm/admin_guide/)
- [Managers, workers e disponibilidade dos nodes](https://docs.docker.com/engine/swarm/how-swarm-mode-works/nodes/)
- [Provider Docker Swarm do Traefik](https://doc.traefik.io/traefik/reference/install-configuration/providers/swarm/)
- [Certificate resolvers do Traefik](https://doc.traefik.io/traefik/v3.6/reference/install-configuration/tls/certificate-resolvers/overview/)
- [Desafios ACME do Let's Encrypt](https://letsencrypt.org/docs/challenge-types/)
- [Simple routing do Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-policy-simple.html)
- [Health checks de Target Groups do ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/target-group-health-checks.html)
- [Certificados HTTPS no ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/https-listener-certificates.html)
- [Headers encaminhados pelo ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html)
