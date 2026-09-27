# Provisionamento manual dos managers EC2

## Referencia do laboratorio

Regiao `us-east-1`:

| Node      | Zona         | Subnet                     | Situacao                             |
| --------- | ------------ | -------------------------- | ------------------------------------ |
| Manager 1 | `us-east-1b` | `subnet-0eaff8c6b2e355ba4` | Existente; IP privado `172.31.87.11` |
| Manager 2 | `us-east-1a` | `subnet-09f3ff80031929574` | Criar                                |
| Manager 3 | `us-east-1c` | `subnet-0459012bb11a4ad9f` | Criar                                |

## 1. Preparar o acesso SSH

O key pair `dslab-swarm-key-pair` do manager 1 corresponde a uma chave privada
que ficou em outro computador. No Console EC2, abra **Network & Security > Key
Pairs > Import key pair** e importe a chave publica atual
`.ssh/swarm.dslab.pub` deste repositorio com o nome `dslab-swarm-current`.
Importe apenas o arquivo `.pub`; a chave privada `.ssh/swarm.dslab` permanece no
seu computador. Se o key pair ja tiver sido importado, apenas selecione-o no
lancamento. Nao escolha o key pair antigo sem ter a chave privada correspondente.

## 2. Atualizar o security group existente

No Console EC2, abra **Network & Security > Security Groups** e
selecione `DSLab Swarm SG` (`sg-025a75482a78a7974`), ja usado pelo
manager 1. Em **Edit inbound rules**, adicione as regras abaixo. Em
**Source**, escolha **Custom** e selecione o proprio `DSLab Swarm SG`
pelo ID.

| Protocolo | Porta | Uso                                     |
| --------- | ----: | --------------------------------------- |
| TCP       |  2377 | Controle e entrada de managers no Swarm |
| TCP       |  7946 | Descoberta entre nodes                  |
| UDP       |  7946 | Descoberta entre nodes                  |
| UDP       |  4789 | Trafego da rede overlay                 |

O manager 1 ja usa esse grupo; associe o mesmo grupo aos managers 2 e 3
no lancamento. As portas do Swarm devem aceitar trafego apenas dos nodes
com esse grupo; nao abra 2377, 7946 ou 4789 para `0.0.0.0/0`.

O grupo atual tambem libera 80/443 publicamente. Ao reutiliza-lo, os novos
managers terao a mesma permissao de entrada nessas portas. Se um servico
publicar essas portas pelo routing mesh do Swarm, os novos nodes poderao
receber conexoes pelos seus IPs, mesmo sem uma replica local. O DNS/EIP do
laboratorio continua apontando para o manager 1.

## 3. Lancar cada instancia

No Console EC2, em `us-east-1`, abra **Instances > Launch instances** e repita o
processo para o manager 2 e o manager 3:

1. **Application and OS Images:** Ubuntu Server 26.04 LTS x86_64. Para repetir
   exatamente a imagem do manager 1, use a AMI `ami-0b6d9d3d33ba97d99` e
   confirme que o proprietario e a Canonical (`099720109477`).
2. **Instance type:** `t3.small`. Mantenha a compra On-Demand; nao selecione
   Spot para um manager do quorum.
3. **Key pair (login):** `dslab-swarm-current`.
4. **Network settings:** escolha a VPC da tabela e a subnet da zona do node
   correspondente. Deixe **Auto-assign public IP = Enable** para SSH e download
   dos pacotes. Em **Firewall**, escolha **Select existing security group** e
   selecione `DSLab Swarm SG` (`sg-025a75482a78a7974`).
5. **Configure storage:** volume raiz gp3 de 30 GiB, com **Delete on
   termination**.

## Referencias

- [Campos de lancamento de uma instancia EC2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-launch-parameters.html)
- [Importar um key pair](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/create-key-pairs.html)
- [Regras de security groups para instancias EC2](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-security-groups.html)
- [Portas necessarias para o Docker Swarm](https://docs.docker.com/engine/swarm/swarm-tutorial/)
- [Routing mesh do Docker Swarm](https://docs.docker.com/engine/swarm/ingress/)
