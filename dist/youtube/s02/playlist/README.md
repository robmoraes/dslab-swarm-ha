# Playlist - Temporada 2

## Título

Docker Swarm HA na AWS | Temporada 2: Alta Disponibilidade

## Descrição

```text
Na primeira temporada, construímos a base do laboratório em uma única instância EC2: Docker Swarm, API, Traefik, WAF, domínio e HTTPS. Nesta segunda temporada, vamos distribuir essa arquitetura e testar o que realmente acontece quando um componente falha.

Começaremos com um domínio próprio para o laboratório, delegando dslab.dev.br do Registro.br para o Amazon Route 53. Depois, exploraremos quorum e Raft para entender por que o Docker Swarm precisa de três managers para tolerar a perda de um deles.

Com a base preparada, adicionaremos managers e workers em diferentes zonas de disponibilidade, distribuiremos os serviços e investigaremos separadamente a disponibilidade do cluster, das aplicações e da entrada pública. Também evoluiremos a borda com Traefik replicado, Application Load Balancer, target groups e certificados do AWS Certificate Manager.

Ao longo da temporada, você verá:

- delegação de domínio, zona hospedada e registros DNS no Route 53;
- quorum, Raft e recuperação de managers no Docker Swarm;
- distribuição de managers, workers e réplicas;
- testes de falha para identificar pontos únicos de falha;
- Traefik replicado e os limites do DNS como mecanismo de failover;
- ALB, target groups, health checks e certificados ACM;
- validação do caminho da requisição e do IP do cliente;
- testes finais de falha e recuperação da arquitetura.

O objetivo é sair de um laboratório funcional em um único node e chegar a uma plataforma capaz de continuar atendendo mesmo quando um componente importante estiver indisponível.

📁 Repositório do projeto: https://github.com/robmoraes/dslab-swarm-ha
🎬 Temporada 1 - Construindo a Base: https://www.youtube.com/playlist?list=PLd8XIvm0GJ7U
🎬 Série completa: @DistributedSystemsLab
👨‍💻 Sobre o autor: https://about.robmoraes.dev.br

#DockerSwarm #AWS #AltaDisponibilidade #DevOps #SistemasDistribuidos
```
