# Provisionamento manual do ALB pelo console AWS

Publicar HTTPS no ALB com o certificado do ACM e encaminhar HTTP para os tres
Traefiks. Ao final, as portas 80/443 dos nodes deixam de aceitar acesso direto
da Internet.

## 1. Conferir os pre-requisitos

- Regiao **N. Virginia (`us-east-1`)**.
- Certificado de `swarm.dslab.dev.br` e `*.swarm.dslab.dev.br` em **Issued**,
  conforme o [guia do ACM](../acm/README.md).
- `dslab` com a [policy complementar do ALB](../policies/README.md) associada
  ao grupo `DSLabProvisioners`.
- Managers ligados, `Ready` e `Active`, com Traefik global publicando a porta
  80 em modo `host`; WAF e API nos workers.

Dados conferidos na conta em 2026-10-28; confirme os IDs caso recrie instancias:

| Recurso | Valor |
| --- | --- |
| VPC | `vpc-0d576bafa8d8b3f3d` |
| SG dos nodes | `DSLab Swarm SG` / `sg-025a75482a78a7974` |
| Nome do ALB | `swarm-ha-public-alb` |
| Nome do Target Group | `swarm-ha-traefik-tg` |

| Zona | Subnet publica | Manager a registrar | IP privado |
| --- | --- | --- | --- |
| `us-east-1a` | `subnet-09f3ff80031929574` | `i-082a6a856fb2f48cd` / Manager 2 | `172.31.4.63` |
| `us-east-1b` | `subnet-0eaff8c6b2e355ba4` | `i-0363d71d3af03c457` / Manager 1 | `172.31.87.11` |
| `us-east-1c` | `subnet-0459012bb11a4ad9f` | `i-0fb9525f811335887` / Manager 3 | `172.31.26.247` |

As tres subnets usam a rota para Internet Gateway e possuem enderecos livres.
O ALB usara essas tres zonas; os targets serao as instancias dos managers,
alcancadas pelos IPs privados. Nao cadastre EIPs ou workers no Target Group.

## 2. Preparar o health check do Traefik

No `command` de [stacks/traefik.yml](../../stacks/traefik.yml), acrescente:

```yaml
      - --ping=true
      - --ping.entrypoint=web
      - --entrypoints.web.forwardedheaders.trustedips=172.31.0.0/20,172.31.80.0/20,172.31.16.0/20
```

Os CIDRs sao das subnets do ALB. Eles permitem confiar nos headers encaminhados
pela borda, que sera a unica origem permitida na porta 80 dos nodes.

Em um manager, publique o arquivo atualizado e confirme a convergencia:

```bash
docker stack deploy -c stacks/traefik.yml --detach=false traefik
docker service ps traefik_traefik
```

Execute a partir da raiz do repositorio ou ajuste o caminho ao arquivo copiado.
Em cada manager, confirme `curl -fsS http://127.0.0.1/ping`: deve retornar `OK`.
O [ping do Traefik](https://doc.traefik.io/traefik/v3.6/reference/install-configuration/observability/healthcheck/)
responde sem exigir o dominio no `Host`, que no health check do ALB contem o
[IP privado do target](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-troubleshooting.html).
Esse teste verifica o Traefik; nao verifica a saude do WAF ou da API.

## 3. Criar o Security Group do ALB

Em **EC2 > Security Groups > Create security group**:

1. Nome: `swarm-ha-alb-sg`; descricao: `Entrada publica HTTP e HTTPS do ALB`.
2. VPC: `vpc-0d576bafa8d8b3f3d`.
3. Entrada: TCP **80** e TCP **443**, ambas de `0.0.0.0/0`.
4. Saida: TCP **80**, destino **Custom**, selecionando
   `sg-025a75482a78a7974`. Remova a saida padrao de todos os protocolos.
5. Tags: `Name=swarm-ha-alb-sg`, `dslab:lab-id=swarm-ha`, `dslab:owner=dslab`.

No SG existente dos nodes, adicione uma entrada TCP **80** com origem no
novo SG do ALB. Preserve por enquanto as regras publicas existentes para
permitir a transicao do DNS. O ALB deve usar somente seu SG dedicado.
Veja as [regras recomendadas para ALB e targets](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-update-security-groups.html).

## 4. Criar o Target Group

Em **EC2 > Target Groups > Create target group**, preencha:

| Campo | Valor |
| --- | --- |
| Target type | **Instances** |
| Name | `swarm-ha-traefik-tg` |
| Protocol / Port | **HTTP / 80** |
| IP address type | **IPv4** |
| VPC | `vpc-0d576bafa8d8b3f3d` |
| Protocol version | **HTTP1** |
| Health check protocol / Path | **HTTP / `/ping`** |
| Health check port | **Traffic port** |
| Success codes | **200** |
| Interval / Timeout | **30 / 5 segundos** |
| Healthy / Unhealthy threshold | **2 / 2** |

Adicione as tags `Name=swarm-ha-traefik-tg`, `dslab:lab-id=swarm-ha` e
`dslab:owner=dslab`. Avance, selecione os tres managers da tabela, mantenha a
porta **80**, clique em **Include as pending below** e crie o grupo.
O status pode ser `unused` ate existir um listener associado ao grupo.
[Configuracao dos health checks](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/target-group-health-checks.html).

## 5. Criar o ALB e os listeners

Em **EC2 > Load Balancers > Create load balancer > Application Load Balancer**:

1. Nome: `swarm-ha-public-alb`; **Scheme: Internet-facing**;
   **IP address type: IPv4**.
2. Em **Network mapping**, escolha a VPC e as tres zonas/subnets da tabela.
3. Em **Security groups**, selecione somente `swarm-ha-alb-sg`.
4. Configure um listener **HTTPS :443**, acao **Forward to target group**,
   destino `swarm-ha-traefik-tg`.
5. Em certificado, escolha **From ACM** e selecione o certificado emitido
   para `swarm.dslab.dev.br`, que tambem cobre o wildcard. Use uma security
   policy TLS 1.2/1.3, como `ELBSecurityPolicy-TLS13-1-2-2021-06`.
6. Configure o listener **HTTP :80** com acao **Redirect to URL**:
   protocolo **HTTPS**, porta **443**, codigo **HTTP_301**. Preserve host,
   caminho e query nos campos de redirecionamento.
7. Tags: `Name=swarm-ha-public-alb`, `dslab:lab-id=swarm-ha`,
   `dslab:owner=dslab`. Mantenha integracoes opcionais desativadas nesta etapa.
8. Revise e crie o ALB. Se o assistente nao oferecer redirect, crie inicialmente
   apenas o listener 443 e adicione o listener 80 com redirect na aba
   **Listeners and rules** apos a criacao.

Espere o ALB ficar **Active** e os tres targets ficarem **healthy**.
O [assistente de criacao](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-application-load-balancer.html)
permite associar o certificado ao listener HTTPS. O TLS termina no ALB;
o encaminhamento ao Traefik continua em HTTP :80.

## 6. Apontar o dominio para o ALB

Em **Route 53 > Hosted zones > dslab.dev.br**, edite o registro A existente
`swarm.dslab.dev.br`:

| Campo | Valor |
| --- | --- |
| Record type | **A** |
| Alias | **On** |
| Route traffic to | **Alias to Application and Classic Load Balancer** |
| Region / Load balancer | **us-east-1 / swarm-ha-public-alb** |
| Routing policy | **Simple** |
| Evaluate target health | **Yes** |

Substitua os tres EIPs pelo Alias e salve. Mantenha o wildcard dos servicos
apontando para `swarm.dslab.dev.br` e preserve os CNAMEs de validacao do ACM.
O [Alias do Route 53](https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/routing-to-elb-load-balancer.html)
acompanha o endereco do ALB; nao copie seus IPs para um registro A comum.

Aguarde a atualizacao e os caches do registro anterior. No navegador, abra
`https://whoami.swarm.dslab.dev.br/metadata` e
`https://excalidraw.swarm.dslab.dev.br`. Confira o certificado e teste o
redirecionamento ao abrir esses dominios em HTTP. O hostname `*.elb.amazonaws.com`
do ALB nao corresponde ao certificado nem aos routers do laboratorio.

## 7. Fechar o acesso publico direto aos nodes

Depois de confirmar HTTPS e DNS pelo ALB, em **EC2 > Security Groups >
DSLab Swarm SG > Edit inbound rules**:

1. Remova TCP **80** e TCP **443** com origem `0.0.0.0/0`; remova tambem
   eventuais permissoes publicas equivalentes em IPv6.
2. Mantenha TCP **80** com origem **somente no SG do ALB**.
3. Preserve SSH do seu IP e as regras internas do Swarm, incluindo as
   referencias do SG para ele mesmo nas portas 2377, 7946 e 4789.

Esse SG e compartilhado pelos managers e workers. Confira que nenhum outro
SG associado aos nodes libera 80/443 publicamente. A porta 443 publica fica
no ALB. Os EIPs e o acesso SSH existentes podem permanecer nesta etapa.

## 8. Conferir o resultado

- [ ] ALB `Active`, com os tres managers `healthy` no Target Group.
- [ ] HTTPS funciona pelo dominio e HTTP redireciona para HTTPS.
- [ ] Certificado ACM associado ao listener 443, com validacao DNS preservada.
- [ ] WAF e API continuam nos workers; Traefik continua global nos managers.
- [ ] Porta 80 dos nodes aceita o SG do ALB; 80/443 sem permissao publica direta.
- [ ] A API continua exibindo o IP do cliente nos headers encaminhados.

Ao usar o loop de [requests](../../scripts/requests-whoami.sh), atualize sua
URL para `https://whoami.swarm.dslab.dev.br/metadata`, pois a entrada HTTP
agora responde com redirecionamento. Novos managers precisam ser registrados
manualmente no Target Group; mudancas no Swarm nao atualizam essa lista.
