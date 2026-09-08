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
| 1 | Por que precisamos de três managers? | Pronto — aproximadamente 8 minutos |
| 2 | De um manager para um cluster HA | Planejado |
| 3 | O cluster sobrevive, mas o site cai | Planejado |
| 4 | Managers, workers e separação dos workloads | Planejado |
| 5 | Por que replicar o Traefik ainda não resolve tudo? | Planejado |
| 6 | ALB, Target Groups e ACM | Planejado |
| 7 | Chaos Day: agora temos HA de verdade? | Planejado |

## Escopo

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
