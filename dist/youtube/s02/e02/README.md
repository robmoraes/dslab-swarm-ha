# YouTube - Temporada 2, episódio 02

## Título

Docker Swarm na AWS T2 #02 | Teoria: Réplicas, Falhas e Quorum

## Descrição

```text
#DockerSwarm #AltaDisponibilidade #Traefik #Quorum #AWS

Este é um episódio de teoria e arquitetura, com demonstrações curtas no terminal. Antes de adicionar novos nodes, vamos seguir uma requisição do browser até a API e descobrir por que ter várias réplicas ainda não torna todo o caminho altamente disponível.

Na resolução do domínio, o DNS resolver segue a delegação configurada no Registro.br até os servidores autoritativos do Route 53. Depois de obter o IP do manager 1, o browser envia a requisição: o Security Group permite a entrada, o Traefik encaminha o tráfego, o WAF inspeciona e a API responde.

O WAF tem três réplicas e a API tem cinco. Ao parar um container, as requisições continuam. Mas o que acontece quando o único Traefik cai? E quando o próprio node fica indisponível? Acompanhamos um loop de requisições a cada segundo, paramos e restauramos o Traefik e, em seguida, testamos a queda do node.

Neste episódio, você verá:

- o caminho da requisição, do DNS até a aplicação;
- a diferença entre replicar serviços e eliminar pontos únicos de falha;
- o efeito de parar uma réplica, o Traefik e o node;
- como Raft e quorum mantêm o estado do Docker Swarm;
- por que a próxima etapa do laboratório adotará três managers.

Três managers dão tolerância à perda de um manager no control plane. A entrada pública também precisará de redundância, que construiremos nas próximas etapas da temporada.

📁 Arquivos do projeto: https://github.com/robmoraes/dslab-swarm-ha
🎬 Temporada 2 - Alta Disponibilidade: https://www.youtube.com/playlist?list=PLWXeaGSAACi8
🎬 Temporada 1 - Construindo a Base: https://www.youtube.com/playlist?list=PLd8XIvm0GJ7U
```

## Publicação

- **Playlist:** Docker Swarm HA na AWS | Temporada 2: Alta Disponibilidade.
- **Miniatura:** [capa.png](../../../kdenlive/s02/e02/capa.png).
