# Provisionamento manual dos workers EC2

Preparacao dos tres workers do video 04 da segunda temporada.

## Referencia do laboratorio

Regiao `us-east-1`, VPC `vpc-0d576bafa8d8b3f3d`, a mesma dos managers:

| Node     | Name e hostname      | Zona         | Subnet                     |
| -------- | -------------------- | ------------ | -------------------------- |
| Worker 1 | `swarm-ha-worker-01` | `us-east-1a` | `subnet-09f3ff80031929574` |
| Worker 2 | `swarm-ha-worker-02` | `us-east-1b` | `subnet-0eaff8c6b2e355ba4` |
| Worker 3 | `swarm-ha-worker-03` | `us-east-1c` | `subnet-0459012bb11a4ad9f` |

| Configuracao        | Valor                                        |
| ------------------- | -------------------------------------------- |
| Sistema             | Ubuntu Server 26.04 LTS, `x86_64`            |
| AMI                 | `ami-0b6d9d3d33ba97d99`                      |
| Proprietario da AMI | Canonical, `099720109477`                    |
| Instancia           | `t3.small`: 2 vCPUs e 2 GiB de RAM           |
| Compra              | On-Demand                                    |
| Volume raiz         | 30 GiB, `gp3`, com **Delete on termination** |
| Security group      | `DSLab Swarm SG` (`sg-025a75482a78a7974`)    |
| Key pair            | `dslab-swarm-current`, a importar no passo 1 |
| IPv4 publico        | Automatico, para SSH e downloads             |
| IMDS                | Enabled, V2 only, response hop limit `2`     |

Os workers usam os IPs privados para participar do cluster. Nao e necessario
alocar EIPs nem criar registros DNS para eles. O IPv4 publico automatico pode
mudar depois de parar e iniciar a instancia; consulte o valor atual para SSH.

## 1. Preparar a chave e conferir o security group

Usar a mesma chave dos managers: `.ssh/swarm.dslab.pub`

Em **Network & Security > Security Groups**, confira o grupo existente
`DSLab Swarm SG`. Ele sera reutilizado nos tres workers e ja possui as regras
internas abaixo, com o proprio grupo como origem:

| Protocolo | Porta | Uso                                   |
| --------- | ----: | ------------------------------------- |
| TCP       |  2377 | Entrada e comunicacao com os managers |
| TCP       |  7946 | Descoberta entre nodes                |
| UDP       |  7946 | Descoberta entre nodes                |
| UDP       |  4789 | Trafego da rede overlay               |

## 2. Lancar as tres instancias

Em **Instances > Launch instances**, lance uma instancia por vez para escolher
a subnet de cada zona conforme a tabela:

1. **Names:** use `Swarm Cluster Worker 1`, `Swarm Cluster Worker 2` e
   `Swarm Cluster Worker 3`.
2. **Application and OS Images:** selecione Ubuntu Server 26.04 LTS `x86_64`,
   AMI `ami-0b6d9d3d33ba97d99`, e confirme o proprietario Canonical
   (`099720109477`).
3. **Instance type:** selecione `t3.small` e mantenha a compra On-Demand.
4. **Key pair (login):** selecione `dslab-swarm-current`.
5. **Network settings > Edit:** selecione a VPC `vpc-0d576bafa8d8b3f3d` e a
   subnet do worker. Deixe **Auto-assign public IP = Enable**. Em **Firewall**,
   escolha **Select existing security group** e selecione `DSLab Swarm SG`
   (`sg-025a75482a78a7974`).
6. **Configure storage:** configure o volume raiz como `gp3`, 30 GiB, e
   mantenha **Delete on termination** habilitado.
7. **Advanced details:** configure **Metadata accessible = Enabled**,
   **Metadata version = V2 only** e **Metadata response hop limit = 2**.
8. Lance a instancia e anote os IPs publico e privado exibidos nos detalhes.

O hop limit `2` permite que a API no container consulte o IMDSv2 da EC2 para
obter os metadados exibidos pelo laboratorio.

## 3. Instalar o Docker e definir o hostname

Execute os comandos locais abaixo a partir da raiz deste repositorio,
substituindo `IP_PUBLICO_DO_WORKER` pelo IPv4 publico da instancia:

```bash
scp -i .ssh/swarm.dslab -o IdentitiesOnly=yes \
  scripts/bootstrap-worker.sh \
  ubuntu@IP_PUBLICO_DO_WORKER:~/bootstrap-worker.sh

ssh -i .ssh/swarm.dslab -o IdentitiesOnly=yes ubuntu@IP_PUBLICO_DO_WORKER
```

```bash
sudo hostnamectl set-hostname swarm-ha-worker-01
sudo bash ~/bootstrap-worker.sh
exit
```

## 4. Vincular cada worker ao cluster

Em um manager existente, obtenha o comando de entrada para **worker**:

```bash
docker swarm join-token worker
```
