# Policies para provisionamento pelo console AWS

Pacotes complementares para os [guias de ACM](../acm/README.md) e
[ALB](../alb/README.md), usando a regiao `us-east-1`.

| Policy | Finalidade | Situacao na conferencia |
| --- | --- | --- |
| [ProvisionACMCertificates](ProvisionACMCertificates.json) | Solicitar certificados, consultar detalhes e adicionar tags | Ja associada ao grupo `DSLabProvisioners` |
| [ProvisionApplicationLoadBalancers](ProvisionApplicationLoadBalancers.json) | Criar ALB, Target Group e listeners; configurar health check, targets e tags | Preparada para associar ao grupo |

## Tags das policies no video 06

Ao criar cada policy no console IAM, use somente estas tres tags do
[tag-pack](../../tag-pack.md):

| Chave | Valor |
| --- | --- |
| `Name` | Nome amigavel da policy, conforme a tabela abaixo |
| `Owner` | `dslab` |
| `Project` | `swarm.dslab` |

| Nome tecnico da policy | Valor da tag `Name` |
| --- | --- |
| `ProvisionACMCertificates` | `Swarm HA ACM Provisioning Policy` |
| `ProvisionApplicationLoadBalancers` | `Swarm HA ALB Provisioning Policy` |

As tags sao metadados da policy: preencha-as na tela de criacao do IAM.
Os arquivos JSON definem as permissoes concedidas pela policy.

## Conferencia do dslab

Consulta realizada em 2026-10-28 com a identidade administrativa `cluster-admin`,
usando leitura do IAM e `SimulatePrincipalPolicy` para o usuario `dslab`.
O usuario recebe as permissoes pelo grupo `DSLabProvisioners` e nao possui
permissions boundary, policies inline ou policies diretamente associadas.

| Acoes verificadas | Resultado antes da nova policy |
| --- | --- |
| Criar ALB, Target Group e listener; modificar listener e health check; registrar targets; consultar ELB e adicionar tags | `implicitDeny` |
| Consultar EC2/rede, criar SG, adicionar/remover regras de entrada e saida e aplicar tags | `allowed` |
| Solicitar e consultar certificados; adicionar tags no ACM | `allowed` |
| Listar hosted zones e alterar registros no Route 53 | `allowed` |
| `iam:CreateServiceLinkedRole` para ELB | `allowed` |

A policy `DSLabProvisionerPolicy` vigente ja concede `ec2:*`, `route53:*` e
`iam:CreateServiceLinkedRole`. A `ProvisionACMCertificates` foi copiada da
versao `v1` efetivamente associada ao grupo. A role
`AWSServiceRoleForElasticLoadBalancing` tambem ja existe.

Na conferencia acima, faltava associar `ProvisionApplicationLoadBalancers`.
Os JSONs sao complementares: o pacote ALB depende das permissoes existentes
de EC2/Security Groups, Route 53 e ACM. Nao substitua a policy base por ele.

A simulacao verifica as permissoes IAM, sem criar um ALB para testar. A
emissao do certificado do laboratorio foi confirmada como `ISSUED`.

Com o JSON ALB adicionado apenas como entrada da simulacao, as 32 acoes
verificadas para o fluxo resultaram em `allowed`. Os dois JSONs passaram no
IAM Access Analyzer sem apontamentos. Esses resultados registram a conferencia
de 2026-10-28; para preparar a nova gravacao, siga os passos abaixo.

## Associar pelo console

Com uma identidade administradora:

1. Abra **IAM > Policies > Create policy > JSON**.
2. Para cada pacote, cole o conteudo do JSON correspondente e use o nome
   tecnico da tabela. Antes de criar, preencha `Name`, `Owner` e `Project`
   com os valores indicados acima.
3. Em **IAM > User groups > DSLabProvisioners > Permissions**, associe as duas
   policies. Todos os membros desse grupo recebem as concessoes.
4. Atualize o console usado pelo `dslab` e inicie o guia do ALB em `us-east-1`.

O JSON do ACM preserva a configuracao que funcionou na gravacao anterior.
As concessoes regionais abrangem os recursos da conta
nessa regiao, sem restricao por nome ou tag. O pacote ALB permite criar e
ajustar a entrada, sem conceder exclusao de load balancers ou certificados.

As acoes de ELB e a dependencia de `AddTags` na criacao estao na
[referencia IAM do ELBv2](https://docs.aws.amazon.com/service-authorization/latest/reference/list_elbv2.html).
A permissao para criar a role fica limitada ao
[servico Elastic Load Balancing](https://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/elb-service-linked-roles.html).
