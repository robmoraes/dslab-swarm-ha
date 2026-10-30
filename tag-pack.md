# Tag Pack (8 tags)

1.  Name = Friendly Resource Name
2.  Owner = dslab
3.  Group = DS Lab | Personal Projects
4.  Project = swarm.dslab
5.  Category =
6.  Purpose = learn | work | ...
7.  Environment = prd | stg | dev | lab
8.  ManagedBy = manual | terraform

| Tag       | Valor                                                         |
| --------- | ------------------------------------------------------------- |
| `Name`    | Nome amigável específico do recurso, conforme a tabela abaixo |
| `Owner`   | `dslab`                                                       |
| `Project` | `swarm.dslab`                                                 |

| Recurso               | Valor de `Name`                    |
| --------------------- | ---------------------------------- |
| Certificado ACM       | `Swarm HA Public Certificate`      |
| Security Group do ALB | `Swarm HA ALB SG`                  |
| Target Group          | `Swarm HA Traefik TG`              |
| ALB                   | `Swarm HA Public ALB`              |
| Listener HTTP :80     | `Swarm HA HTTP Listener`           |
| Listener HTTPS :443   | `Swarm HA HTTPS Listener`          |
| Policy IAM do ACM     | `Swarm HA ACM Provisioning Policy` |
| Policy IAM do ALB     | `Swarm HA ALB Provisioning Policy` |
