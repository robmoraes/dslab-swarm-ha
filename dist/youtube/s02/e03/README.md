# YouTube - Temporada 2, episódio 03

## Título

Docker Swarm na AWS T2 #03 | Três Managers: Quorum, Falhas e Entrada via DNS

## Descrição

```text
#DockerSwarm #AltaDisponibilidade #Traefik #Route53 #AWS

Ter três managers mantém o cluster disponível, mas isso basta para manter a aplicação acessível? Neste episódio prático, ampliamos o laboratório na AWS e testamos o que continua funcionando quando o Traefik ou um manager cai.

Criamos duas novas instâncias EC2 em zonas de disponibilidade distintas, instalamos o Docker e adicionamos os managers ao Swarm. Com três managers, acompanhamos o quorum, a eleição de líder e a distribuição das réplicas do WAF e da API.

Primeiro, mantemos um único Traefik com HTTPS. Demonstramos a queda desse proxy e as falhas dos managers para separar a disponibilidade do cluster da disponibilidade da entrada pública. O terminal dividido acompanha os containers de cada node, o estado dos managers e um loop externo de requisições.

Depois, evoluímos a entrada em tempo real: desativamos HTTPS e ACME, recriamos o Traefik em modo global nos managers e configuramos o registro A do Route 53 com os três Elastic IPs. O laboratório passa a operar em HTTP, permitindo explorar a distribuição por DNS sem a complexidade de certificados concorrentes.

Você verá na prática:

- provisionamento das duas novas EC2 e bootstrap do Docker;
- entrada dos managers no Swarm e identificação de Leader e Reachable;
- redistribuição dos serviços, com três réplicas do WAF e cinco da API;
- testes de falha do Traefik único e dos managers;
- a diferença entre quorum do control plane e disponibilidade da aplicação;
- mudança do Traefik de replicated para global, com porta 80 em modo host;
- configuração de um registro A com os três EIPs no Route 53;
- observação das respostas e da indisponibilidade pelo loop de requisições.

O DNS simples continua anunciando os IPs mesmo quando um destino fica indisponível. O cliente pode tentar outro endereço e ainda concluir a chamada, mas isso não significa que exista um health check removendo o destino que falhou.

HTTP é uma etapa experimental do laboratório, não a arquitetura final. O gancho desta demonstração é introduzir ALB e ACM: uma entrada que verifica a saúde dos targets e devolve HTTPS sem deixar a emissão dos certificados nas réplicas do Traefik.

Arquivos do projeto: https://github.com/robmoraes/dslab-swarm-ha
Temporadas e playlists: https://www.youtube.com/@DistributedSystemsLab/playlists
```

## Tags do vídeo

```text
Docker Swarm, Docker Swarm na AWS, AWS, Amazon EC2, alta disponibilidade, quorum, Raft, Swarm managers, três managers, eleição de líder, Traefik, Traefik global, Docker Swarm global mode, Route 53, Amazon Route 53, Elastic IP, EIP, DNS, DNS round robin, health check, failover, WAF, whoami, sistemas distribuídos, DevOps, ALB, ACM
```

## Publicação

- **Arquivo:** [s02e03.mp4](s02e03.mp4).
- **Duração do arquivo:** 49min55s.
- **Playlist:** Docker Swarm HA na AWS | Temporada 2: Alta Disponibilidade.
- **Categoria:** Ciência e tecnologia.
- **Idioma do vídeo:** Português (Brasil).
- **Miniatura:** [capa.png](../../../kdenlive/s02/e03/capa.png).

## Conferência Antes de Publicar

- Revisar a gravação para ocultar tokens de entrada no Swarm, chaves privadas e conteúdo do `acme.json`.
- A descrição usa a página de playlists do canal. Os links individuais usados nos episódios anteriores não foram confirmados; substituir somente quando houver os URLs completos validados.
- Não incluir capítulos com horários estimados; adicionar apenas após conferir os pontos de corte no vídeo final.
