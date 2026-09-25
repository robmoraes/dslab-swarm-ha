# Padrão mínimo de tags para laboratórios AWS

Este documento define o conjunto mínimo de tags para organizar os recursos dos
laboratórios do Distributed Systems Lab na AWS. O padrão tem três objetivos:

- identificar recursos com facilidade no console;
- agrupar custos por laboratório;
- registrar quem é responsável por cada recurso.

## Tags obrigatórias

| Chave | Finalidade | Exemplo |
| --- | --- | --- |
| `Name` | Nome legível e específico do recurso | `swarm-ha-manager-01` |
| `dslab:lab-id` | Identificador comum a todos os recursos do laboratório | `swarm-ha` |
| `dslab:owner` | Usuário ou equipe responsável pelo recurso | `dslab` |

Exemplo completo:

```text
Name=swarm-ha-manager-01
dslab:lab-id=swarm-ha
dslab:owner=dslab
```

### `Name`

Use um nome que permita reconhecer o recurso sem precisar abrir seus detalhes.
O formato recomendado é:

```text
<lab-id>-<componente>-<função ou sequência>
```

Exemplos:

```text
swarm-ha-manager-01
swarm-ha-worker-02
swarm-ha-public-alb
swarm-ha-app-target-group
```

Mantenha a letra maiúscula em `Name`. Muitos serviços da AWS usam exatamente
essa chave para mostrar o nome amigável no console.

### `dslab:lab-id`

Essa tag relaciona todos os recursos que pertencem ao mesmo laboratório. Use um
identificador curto, estável e em `kebab-case`:

```text
dslab:lab-id=swarm-ha
```

O valor não deve conter o tipo do recurso. EC2, volumes, snapshots, security
groups, load balancers e demais componentes do mesmo laboratório devem receber
o mesmo `lab-id`.

Depois de ativada como tag de alocação de custos, essa chave pode ser usada no
Cost Explorer para filtrar ou agrupar os gastos de cada laboratório.

### `dslab:owner`

Essa tag identifica quem responde pelo recurso e deve decidir sobre sua
manutenção ou remoção:

```text
dslab:owner=dslab
```

Use um identificador estável, como o nome do usuário IAM, uma equipe ou um alias
interno. Não use nome completo, e-mail, senha, token ou outra informação pessoal
ou sensível.

`dslab:owner` não é preenchida automaticamente pela AWS. Antes de criar os
recursos, consulte a identidade autenticada:

```bash
aws sts get-caller-identity --query Arn --output text
```

Para uma identidade como `arn:aws:iam::123456789012:user/dslab`, use
`dslab:owner=dslab`. Ao trabalhar com roles temporárias, prefira um
alias estável da pessoa ou equipe em vez do nome da sessão.

`dslab:owner` representa responsabilidade, não uma trilha de auditoria. Para
descobrir qual identidade criou, alterou ou removeu um recurso, consulte os
eventos do AWS CloudTrail.

## Regras de preenchimento

- Aplique as três tags a todo recurso que ofereça suporte a tags.
- Escreva `Name` exatamente com essa capitalização.
- Escreva as chaves `dslab:*` em minúsculas e com hífens.
- Use valores em minúsculas e `kebab-case`, exceto quando o identificador de
  origem possuir outro formato obrigatório.
- Não altere o `lab-id` durante a vida do laboratório.
- Não armazene informações pessoais, credenciais ou outros segredos em tags.
- Inclua também recursos auxiliares que geram custos, como volumes EBS,
  snapshots, Elastic IPs, NAT Gateways e load balancers.

## Aplicação com AWS CLI

Para adicionar ou atualizar as tags de uma instância EC2:

```bash
aws ec2 create-tags \
  --resources i-0123456789abcdef0 \
  --tags \
    Key=Name,Value=swarm-ha-manager-01 \
    Key=dslab:lab-id,Value=swarm-ha \
    Key=dslab:owner,Value=dslab
```

Para consultar os recursos encontrados pela Resource Groups Tagging API para um
laboratório:

```bash
aws resourcegroupstaggingapi get-resources \
  --tag-filters Key=dslab:lab-id,Values=swarm-ha \
  --query 'ResourceTagMappingList[].ResourceARN' \
  --output table
```

Alguns tipos de recurso possuem APIs próprias ou não aparecem nessa consulta.
Confirme a cobertura de tags nos serviços usados pelo laboratório.

## Aplicação com Terraform

Use `default_tags` para propagar `lab-id` e `owner` a todos os recursos
compatíveis do provider:

```hcl
variable "lab_id" {
  type = string
}

variable "owner" {
  type = string
}

provider "aws" {
  region = "us-east-1"

  default_tags {
    tags = {
      "dslab:lab-id" = var.lab_id
      "dslab:owner"  = var.owner
    }
  }
}
```

Defina `Name` em cada recurso porque seu valor deve descrever aquele componente:

```hcl
resource "aws_instance" "manager" {
  # Demais configurações da instância.

  tags = {
    Name = "${var.lab_id}-manager-01"
  }
}
```

Antes de aplicar, confira o plano para garantir que as tags serão propagadas:

```bash
terraform plan
```

## Ativação para análise de custos

Criar a tag nos recursos não a habilita automaticamente nos relatórios de
custos. Depois que `dslab:lab-id` tiver sido aplicada a pelo menos um recurso:

1. Abra **Billing and Cost Management** no console da AWS.
2. Acesse **Cost Allocation Tags**.
3. Abra a lista de tags definidas pelo usuário.
4. Localize e selecione `dslab:lab-id`.
5. Escolha **Activate**.

Também é possível ativar `dslab:owner` para analisar custos por responsável.
A chave pode levar até 24 horas para aparecer na página e até mais 24 horas para
ficar ativa. A alocação não é retroativa: os custos anteriores à ativação não
serão classificados por essa tag.

Depois da ativação, no Cost Explorer:

1. Abra um relatório de custos.
2. Em **Group by**, selecione **Tag**.
3. Escolha `dslab:lab-id` ou `dslab:owner`.

## Checklist do laboratório

Antes de considerar um laboratório pronto, confirme:

- [ ] todos os recursos possuem `Name`;
- [ ] todos os recursos compartilham o mesmo `dslab:lab-id`;
- [ ] todos os recursos possuem um `dslab:owner` válido;
- [ ] volumes, snapshots e componentes de rede também foram verificados;
- [ ] `dslab:lab-id` está ativa como tag de alocação de custos;
- [ ] a consulta pela Tagging API retorna os recursos esperados.

## Referências

- [Boas práticas e estratégias para tags](https://docs.aws.amazon.com/tag-editor/latest/userguide/best-practices-and-strats.html)
- [Ativação de tags de alocação de custos](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/activating-tags.html)
- [Tags definidas pelo usuário para alocação de custos](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/custom-tags.html)
