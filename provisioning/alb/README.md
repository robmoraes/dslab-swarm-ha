# Provisionamento manual do ALB pelo console AWS

Publicar HTTPS no ALB com o certificado do ACM e encaminhar HTTP para os tres
Traefiks. Ao final, as portas 80/443 dos nodes deixam de aceitar acesso direto
da Internet.

## Permissões para ALB

- [policy complementar do ALB](../policies/README.md)

## Pre-requisitos

- Regiao **N. Virginia (`us-east-1`)**.
- Certificado de `swarm.dslab.dev.br` e `*.swarm.dslab.dev.br` em **Issued** no ACM.
- Managers ligados, `Ready` e `Active`, com Traefik global publicando a porta
  80 em modo `host`; WAF e API nos workers.

## Rede

### VPC

A VPC deve ser a mesma usada pelos nodes do cluster;

- `vpc-0d576bafa8d8b3f3d`

### Subnets

Abaixo a zona dos managers que devem ser usadas no Target Group.

| Zona         | Subnet publica             | Manager a registrar               | IP privado      |
| ------------ | -------------------------- | --------------------------------- | --------------- |
| `us-east-1a` | `subnet-09f3ff80031929574` | `i-082a6a856fb2f48cd` / Manager 2 | `172.31.4.63`   |
| `us-east-1b` | `subnet-0eaff8c6b2e355ba4` | `i-0363d71d3af03c457` / Manager 1 | `172.31.87.11`  |
| `us-east-1c` | `subnet-0459012bb11a4ad9f` | `i-0fb9525f811335887` / Manager 3 | `172.31.26.247` |

## Criar o Security Group do ALB

Em **EC2 > Security Groups > Create security group**:

- Nome: `swarm-ha-alb-sg`;
- descricao: `Entrada publica HTTP e HTTPS do ALB`.
- VPC: `vpc-0d576bafa8d8b3f3d`.
- Inbound: TCP **80** e TCP **443**, ambas de `0.0.0.0/0`.
- Outbound: TCP **80**, destino **Custom**, selecionando `sg-025a75482a78a7974`.
  Remover a saida padrao de todos os protocolos.

### Tags

- `Name=Swarm HA ALB SG`
- `Owner=dslab`
- `Project=swarm.dslab`

## SG do Cluster (existente)

No SG existente dos nodes, adicionar uma entrada TCP **80** com origem
no novo SG do ALB.

Não precisa remover agora a 80 e 443 desse SG.

## Criar o Target Group

Em **EC2 > Target Groups > Create target group**, preencher:

| Campo                         | Valor                   |
| ----------------------------- | ----------------------- |
| Target type                   | **Instances**           |
| Name                          | `swarm-ha-traefik-tg`   |
| Protocol / Port               | **HTTP / 80**           |
| IP address type               | **IPv4**                |
| VPC                           | `vpc-0d576bafa8d8b3f3d` |
| Protocol version              | **HTTP1**               |
| Health check protocol / Path  | **HTTP / `/ping`**      |
| Health check port             | **Traffic port**        |
| Success codes                 | **200**                 |
| Interval / Timeout            | **30 / 5 segundos**     |
| Healthy / Unhealthy threshold | **2 / 2**               |

### Tags

- `Name=Swarm HA Traefik TG`
- `Owner=dslab`
- `Project=swarm.dslab`

## Criar o ALB e os listeners

### ALB

Em **EC2 > Load Balancers > Create load balancer > Application Load Balancer**:

- Nome: `swarm-ha-public-alb`; **Scheme: Internet-facing**; **IP address type: IPv4**.
- Em **Network mapping**, escolher a VPC e as tres zonas/subnets equivalente aos managers.
- Em **Security groups**, selecionar somente `swarm-ha-alb-sg`.

#### Tags do ALB

- `Name=Swarm HA Public ALB`
- `Owner=dslab`
- `Project=swarm.dslab`

### Listeners

#### HTTPS :443

- Configurar um listener **HTTPS :443**
- Acao **Forward to target group**
- Destino `swarm-ha-traefik-tg`.
- Em certificado: **From ACM**;
- Usar uma security policy TLS 1.2/1.3, como `ELBSecurityPolicy-TLS13-1-2-2021-06`;

##### Tags

- `Name=Swarm HA HTTPS Listener`
- `Owner=dslab`
- `Project=swarm.dslab`

#### HTTP :80

- Configurar o listener **HTTP :80**
- Action: **Redirect to URL**
- Protocolo **HTTPS**
- Porta **443**
- Codigo **HTTP_301**.
- Preserve host, caminho e query nos campos de redirecionamento.

##### Tags

- `Name=Swarm HA HTTP Listener`
- `Owner=dslab`
- `Project=swarm.dslab`

### Conferir

Esperar o ALB ficar **Active** e os tres targets ficarem **healthy**.

## `Route 53` Apontar o dominio para o ALB

Em **Route 53 > Hosted zones > dslab.dev.br**,

Editar o registro A existente `swarm.dslab.dev.br`:

| Campo                  | Valor                                              |
| ---------------------- | -------------------------------------------------- |
| Record type            | **A**                                              |
| Alias                  | **On**                                             |
| Route traffic to       | **Alias to Application and Classic Load Balancer** |
| Region / Load balancer | **us-east-1 / swarm-ha-public-alb**                |
| Routing policy         | **Simple**                                         |
| Evaluate target health | **Yes**                                            |

Substituir os tres EIPs pelo Alias e salvar.

Aguardar a atualizacao e os caches do registro anterior. No navegador, abrir
`https://whoami.swarm.dslab.dev.br/metadata` e
`https://excalidraw.swarm.dslab.dev.br`.

## `SG do Cluster` Fechar o acesso publico direto aos nodes

Depois de confirmar HTTPS e DNS pelo ALB, em **EC2 > Security Groups >
DSLab Swarm SG > Edit inbound rules**:

- Remover TCP **80** e TCP **443** com origem `0.0.0.0/0`;
