# Deploy do traefik:

Antes do primeiro deploy deve ser criada uma rede overlay que vai
ser compartilhada com todos os serviços que precisem expor DNS;

```bash
docker network create --driver overlay --attachable traefik-public
```
