# ALB Checklist

- SG do cluster: 80 e 443
- R53: dns para EIPs dos manager
- ALB: swarm-ha-public-alb
- TG: swarm-ha-traefik-tg
- SG do ALB: Swarm HA ALB SG
- ACM: swarm.dslab.dev.br
- R53: record CNAME do desafio
- IAM: Grupo de Usuario anexada policies
- IAM: policies para ACM e ALB
- Traefik: Health Check e trustedips
- Script de Requests: https
