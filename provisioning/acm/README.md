# Configuracao manual do certificado ACM

Emitir pelo console AWS o certificado publico do laboratorio para uso no ALB.
Esta etapa termina com o certificado em **Issued**; ainda nao ativa HTTPS na
aplicacao nem altera o Traefik.

Antes de comecar, a hosted zone publica `dslab.dev.br` deve estar delegada ao
Route 53. Use um usuario ou role com permissao para solicitar e adicionar tags
a certificados no ACM e criar registros nessa zona.

O pacote que funcionou para o `dslab` esta salvo em
[ProvisionACMCertificates.json](../policies/ProvisionACMCertificates.json).
Ele ja esta associado ao grupo `DSLabProvisioners`; consulte a
[conferencia de permissoes](../policies/README.md) para os demais pacotes.

## 1. Solicitar o certificado

1. Abra o [console do ACM](https://console.aws.amazon.com/acm/home?region=us-east-1)
   e confirme a regiao **N. Virginia (`us-east-1`)**, onde criaremos o ALB.
   Certificado e ALB precisam estar na mesma regiao.
   [Regioes e uso do ACM](https://aws.amazon.com/certificate-manager/faqs/).
2. Clique em **Request a certificate** e escolha um certificado publico.
3. Preencha os campos:

   | Campo | Valor |
   | --- | --- |
   | Primeiro dominio | `swarm.dslab.dev.br` |
   | Dominio adicional | `*.swarm.dslab.dev.br` |
   | Exportacao | Desabilitada; nao habilitar **Enable export** |
   | Metodo de validacao | **DNS validation** |
   | Algoritmo da chave | **RSA 2048** |

   O wildcard cobre nomes como `whoami.swarm.dslab.dev.br` e
   `excalidraw.swarm.dslab.dev.br`, mas nao cobre `swarm.dslab.dev.br` nem
   subdominios com mais niveis. Por isso, inclua os dois nomes.

4. Adicione as [tags do laboratorio](../../docs/aws-tags.md):

   | Chave | Valor |
   | --- | --- |
   | `Name` | `swarm-ha-public-certificate` |
   | `dslab:lab-id` | `swarm-ha` |
   | `dslab:owner` | `dslab` |

5. Clique em **Request** e abra os detalhes do certificado criado. O estado
   inicial normalmente sera **Pending validation**.
   [Solicitacao pelo console](https://docs.aws.amazon.com/acm/latest/userguide/acm-public-certificates.html).

## 2. Validar no Route 53

1. Nos detalhes do certificado, localize **Domains** e clique em
   **Create records in Route 53**.
2. Confira os registros de validacao para os dois nomes e a hosted zone
   publica `dslab.dev.br`. Confirme em **Create records**.
3. Volte aos detalhes do certificado, atualize a pagina e aguarde **Issued**.
   A mudanca de estado pode levar ate 30 minutos apos criar os registros.

A validacao do ACM usa **CNAME**, nao o TXT do desafio ACME DNS-01. O dominio
base e o wildcard podem mostrar o mesmo CNAME; nesse caso, um unico registro
atende aos dois. [Validacao DNS do ACM](https://docs.aws.amazon.com/acm/latest/userguide/dns-validation.html).

Se o botao nao estiver disponivel, confira as permissoes e se o dominio ja
foi validado. Para criar manualmente, copie o nome e o valor fornecidos pelo
ACM e acesse **Route 53 > Hosted zones > dslab.dev.br > Create record**,
escolhendo o tipo **CNAME**. Confira o nome completo para nao repetir o sufixo
`dslab.dev.br`.

**Nao altere o registro A com os tres EIPs nem o wildcard dos servicos.**
Adicione somente os CNAMEs de validacao fornecidos pelo ACM. Nao e necessario
abrir portas, mudar Security Groups ou ligar os nodes para validar pelo DNS.

## 3. Conferir e encerrar

- [ ] Certificado na regiao `us-east-1`, com status **Issued**.
- [ ] Nomes `swarm.dslab.dev.br` e `*.swarm.dslab.dev.br` presentes e validados.
- [ ] As tres tags aplicadas.
- [ ] CNAMEs de validacao mantidos na hosted zone publica.

Nao remova esses CNAMEs depois da emissao. Para este certificado nao
exportavel, a renovacao gerenciada depende dos registros acessiveis no DNS
publico e do certificado em uso num servico integrado, como o ALB.
[Renovacao com validacao DNS](https://docs.aws.amazon.com/acm/latest/userguide/dns-renewal-validation.html).

O certificado pode continuar sem recurso associado por enquanto. O proximo
passo e seleciona-lo no listener **HTTPS :443** do ALB, sem exportar chaves
ou instalar certificados nos nodes. A aplicacao permanece em HTTP ate essa
integracao ser concluida.

Continue pelo [guia manual do ALB](../alb/README.md).
