# Policies para provisionamento pelo console AWS

Pacotes complementares para os [guias de ACM](../acm/README.md) e
[ALB](../alb/README.md), usando a regiao `us-east-1`.

| Policy | Finalidade | Situacao na conferencia |
| --- | --- | --- |
| [ProvisionACMCertificates](ProvisionACMCertificates.json) | Solicitar certificados, consultar detalhes e adicionar tags | Ja associada ao grupo `DSLabProvisioners` |
| [ProvisionApplicationLoadBalancers](ProvisionApplicationLoadBalancers.json) | Criar ALB, Target Group e listeners; configurar health check, targets e tags | Preparada para associar ao grupo |

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

O que falta para este fluxo e associar `ProvisionApplicationLoadBalancers`.
Os JSONs sao complementares: o pacote ALB depende das permissoes existentes
de EC2/Security Groups, Route 53 e ACM. Nao substitua a policy base por ele.

A simulacao verifica as permissoes IAM, sem criar um ALB para testar. A
emissao do certificado do laboratorio foi confirmada como `ISSUED`.

Com o JSON ALB adicionado apenas como entrada da simulacao, as 32 acoes
verificadas para o fluxo resultaram em `allowed`. Os dois JSONs passaram no
IAM Access Analyzer sem apontamentos. A nova policy ainda precisa ser criada
e associada ao grupo na conta.

## Associar pelo console

Com uma identidade administradora:

1. Abra **IAM > Policies > Create policy > JSON**.
2. Cole o conteudo de `ProvisionApplicationLoadBalancers.json` e crie a policy
   com o nome `ProvisionApplicationLoadBalancers`.
3. Em **IAM > User groups > DSLabProvisioners > Permissions**, associe a nova
   policy. Todos os membros desse grupo recebem a concessao.
4. Atualize o console usado pelo `dslab` e inicie o guia do ALB em `us-east-1`.

O JSON do ACM fica como referencia da configuracao que funcionou; a policy
ja esta associada. As concessoes regionais abrangem os recursos da conta
nessa regiao, sem restricao por nome ou tag. O pacote ALB permite criar e
ajustar a entrada, sem conceder exclusao de load balancers ou certificados.

As acoes de ELB e a dependencia de `AddTags` na criacao estao na
[referencia IAM do ELBv2](https://docs.aws.amazon.com/service-authorization/latest/reference/list_elbv2.html).
A permissao para criar a role fica limitada ao
[servico Elastic Load Balancing](https://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/elb-service-linked-roles.html).
