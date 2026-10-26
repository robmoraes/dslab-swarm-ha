# Temporada 2 — Alta Disponibilidade

> Série: Cluster Docker Swarm HA na AWS

**Status:** em produção.

Esta temporada continua a partir da base construída na
[temporada 1](../01-building-the-foundation/) e evolui o laboratório para uma
arquitetura realmente tolerante a falhas.

## Vídeos

O [roteiro editorial e técnico](docs/video-roadmap.md) detalha a narrativa, os
experimentos, os resultados esperados e os critérios de aceite da temporada.

| Ordem | Vídeo | Status |
| ---: | --- | --- |
| 1 | Do Registro.br ao Route 53 | Planejado |
| 2 | Por que precisamos de três managers? | Pronto — aproximadamente 8 minutos |
| 3 | Três managers: quorum, falhas e entrada via DNS | Gravado — aproximadamente 45 minutos |
| 4 | Managers, workers e separação dos workloads | Planejado |
| 5 | Por que ALB e ACM? Certificados e segurança | Gravado: aproximadamente 20 minutos, com bônus do Excalidraw |
| 6 | ACM, ALB e Target Groups: publicando HTTPS | Planejado |
| 7 | Chaos Day: agora temos HA de verdade? | Planejado |

O episódio 3 reúne os antigos roteiros dos vídeos 3 e 4. A gravação também
inclui a transição para Traefik global em HTTP e o registro A com os três EIPs,
evidenciando os limites do DNS direto e preparando o gancho para ALB e ACM.

Para preparar os workers do episódio 4, consulte o
[guia de provisionamento EC2](../../provisioning/workers/ec2.md) e o
[script de bootstrap](../../scripts/bootstrap-worker.sh).

O episódio 5 retoma a evolução até três managers e três workers e justifica
ACM + ALB pelo retorno obrigatório do HTTPS e pela proteção das portas 80/443
dos nodes contra acesso direto da Internet. O certificado foi o motivo
determinante para a escolha.
Os [quatro desenhos editáveis do Excalidraw](../../docs/excalidraw/s02e05/README.md)
foram apresentados na ordem DNS-01, HTTP-01, Traefik único e ACM + ALB + TG,
com o último fazendo o gancho para o vídeo 6. Como bônus, foi feito o deploy
da stack do Excalidraw, usada na apresentação. A emissão do certificado e
a configuração de ALB e Target Group ficam para o episódio 6.

## Escopo

- Delegar o domínio `dslab.dev.br` do Registro.br para uma public hosted zone
  no Route 53.
- Adotar `swarm.dslab.dev.br` como domínio raiz do cluster e
  `*.swarm.dslab.dev.br` para os serviços.
- Adicionar nodes managers e workers ao cluster.
- Explorar quorum, distribuição de réplicas e regras de placement.
- Demonstrar separadamente a disponibilidade do control plane, dos workloads e
  da borda.
- Introduzir Application Load Balancer (ALB) e target groups.
- Usar certificados do AWS Certificate Manager (ACM) na borda.
- Preservar e validar o endereço IP do cliente em todo o caminho da requisição.
- Definir health checks e executar testes de falha e recuperação.

Os artefatos serão adicionados incrementalmente conforme os episódios forem
preparados.

[Voltar ao índice da série](../../README.md)
