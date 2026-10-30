# Configuracao manual do certificado ACM

Emitir pelo console AWS o certificado publico do laboratorio para uso no ALB.
Esta etapa termina com o certificado em **Issued**; ainda nao ativa HTTPS na
aplicacao nem altera o Traefik.

## Permissão para ACM

- [Policy para ACM](../policies/ProvisionACMCertificates.json)
- [Guia de Permissões](../policies/README.md)

## Solicitar o certificado

Na região onde vai estar o ALB acessar o `ACM`;

- Clicar em **Request a certificate**;
- Escolher **Certificado público**;

Preencher os campos:

| Campo               | Valor                                         |
| ------------------- | --------------------------------------------- |
| Primeiro dominio    | `swarm.dslab.dev.br`                          |
| Dominio adicional   | `*.swarm.dslab.dev.br`                        |
| Exportacao          | Desabilitada; nao habilitar **Enable export** |
| Metodo de validacao | **DNS validation**                            |
| Algoritmo da chave  | **RSA 2048**                                  |

O wildcard cobre nomes como `whoami.swarm.dslab.dev.br` e
`excalidraw.swarm.dslab.dev.br`.

Tags

| Chave     | Valor                         |
| --------- | ----------------------------- |
| `Name`    | `Swarm HA Public Certificate` |
| `Owner`   | `dslab`                       |
| `Project` | `swarm.dslab`                 |

Clicar em **Request** e abrir os detalhes do certificado criado.
O estado inicial normalmente sera **Pending validation**.

## Validar no Route 53

Nos detalhes do certificado, localizar **Domains** e clicar em
**Create records in Route 53**.

Confirmar em **Create records**.

Aguardar status **Issued**. A mudanca de estado pode levar ate 30
minutos apos criar os registros.

- A Autoridade Certificadora (CA) é a própria AWS.
- A validacao do ACM usa **CNAME** para o Challenge.
