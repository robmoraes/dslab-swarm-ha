# YouTube - Temporada 2, episódio 05

## Título

Docker Swarm na AWS T2 #05 | Teoria: Por que ALB e ACM? Certificados e Segurança

## Descrição

```text
#DockerSwarm #AltaDisponibilidade #AWS #Traefik #DevOps

Este é um episódio de teoria e arquitetura, com diagramas no Excalidraw e um bônus prático: o deploy da própria ferramenta no Docker Swarm. Vamos entender por que ALB e ACM são as próximas peças do laboratório.

Retomamos o caminho até aqui: primeiro, um Traefik com HTTPS e certificado emitido pelo desafio HTTP-01. Depois, distribuímos o Traefik pelos três managers em modo global e desativamos o HTTPS. Agora temos três managers e três workers, com WAF e API separados do control plane, mas a entrada continua em HTTP.

Devolver HTTPS é o requisito obrigatório e o motivo determinante para a escolha de ALB e ACM neste laboratório. Queremos tirar das réplicas do Traefik a responsabilidade pela emissão e renovação do certificado público, usando o AWS Certificate Manager e encerrando TLS no Application Load Balancer.

O segundo objetivo é ampliar a segurança: impedir o acesso direto da Internet às portas 80 e 443 dos nodes. Na arquitetura proposta, o ALB será a entrada pública e os managers receberão tráfego na porta do Traefik somente a partir do Security Group do ALB. Isso exige configurar as permissões de acesso, não apenas mudar o DNS.

Apresentamos quatro cenários, nesta ordem:

- DNS-01 com três Traefiks: a validação por DNS funciona, mas cada instância mantém seu próprio processo de emissão e renovação;
- HTTP-01 com três Traefiks: a validação pode chegar a uma réplica diferente daquela que possui o desafio;
- Traefik único: a gestão do certificado fica simples, mas a entrada depende de uma única instância;
- ACM + ALB + Target Group: certificado fora das réplicas, HTTPS no ALB e encaminhamento para os Traefiks por IPs privados.

Como bônus, subimos a stack do Excalidraw no próprio cluster e usamos a ferramenta para apresentar os diagramas.

O último desenho mostra a arquitetura que vamos implementar no vídeo 6. Neste episódio não configuramos ACM, ALB nem Target Group: a emissão do certificado, a publicação de HTTPS e a restrição do acesso direto aos nodes ficam para a próxima etapa.

Arquivos do projeto: https://github.com/robmoraes/dslab-swarm-ha
Temporadas e playlists: https://www.youtube.com/@DistributedSystemsLab/playlists
```

## Tags do vídeo

```text
Docker Swarm, Docker Swarm na AWS, AWS, alta disponibilidade, Traefik, Traefik global, ALB, Application Load Balancer, ACM, AWS Certificate Manager, HTTPS, TLS, certificados SSL, Let's Encrypt, ACME, HTTP-01, DNS-01, Route 53, Target Group, Security Group, managers e workers, segurança na AWS, Excalidraw, docker stack deploy, sistemas distribuídos, DevOps
```

## Publicação

- **Arquivo:** [s02e05.mp4](s02e05.mp4).
- **Duração do arquivo:** 18min51s.
- **Playlist:** Docker Swarm HA na AWS | Temporada 2: Alta Disponibilidade.
- **Categoria:** Ciência e tecnologia.
- **Idioma do vídeo:** Português (Brasil).
- **Miniatura:** [cover.png](cover.png).

## Conferência Antes de Publicar

- Manter o caráter explicativo no título e na descrição; a prática de ACM, ALB e Target Group pertence ao vídeo 6.
- A capa e o último diagrama representam a arquitetura de destino, não uma configuração concluída neste episódio.
- Revisar a gravação para ocultar tokens de entrada no Swarm, chaves privadas e outras credenciais.
- A descrição usa a página de playlists do canal; substituir pelo link direto da temporada somente quando o URL completo estiver confirmado.
- Não incluir capítulos com horários estimados; adicionar apenas após conferir os pontos de corte no vídeo final.
