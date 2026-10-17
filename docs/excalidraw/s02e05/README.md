# Desenhos do S02E05

Abra o [quadro completo](../s02e05.excalidraw) no Excalidraw. São quatro
cenários em duas linhas, da esquerda para a direita, também disponíveis
individualmente abaixo.

Textos, formas e conectores são editáveis. Cada quadro tem proporção 16:9,
área de 1920 x 1080, fundo branco e fonte Excalifont. O mural original da
temporada, `s02.excalidraw`, foi preservado.

## Ordem de apresentação

| Quadro | Arquivo | Ponto principal |
| --- | --- | --- |
| 01 | [DNS-01](01-dns-01.excalidraw) | A validação por TXT funciona, mas três processos ACME continuam emitindo e renovando de forma independente. |
| 02 | [HTTP-01](02-http-01.excalidraw) | A CA pode validar em um Traefik diferente daquele que possui o desafio. |
| 03 | [Um Traefik](03-traefik-unico.excalidraw) | Gestão simples do certificado, mas uma única instância sustenta a entrada. |
| 04 | [ACM + ALB + Target Group](04-acm-alb-tg.excalidraw) | Certificado fora das réplicas, TLS no ALB e encaminhamento para três Traefiks por IPs privados. |

## Sequência da gravação

1. Relembrar que o laboratório passou a usar Traefik global em HTTP no vídeo 3.
2. Explicar os três primeiros cenários e seus respectivos limites, sem
   implementá-los no cluster.
3. Apresentar o quarto cenário como a arquitetura escolhida. O ACM gerencia
   o certificado; o ALB recebe HTTPS; o Target Group reúne os destinos HTTP.
4. Sair do desenho para o console da AWS: solicitar o certificado, criar
   os CNAMEs de validação no Route 53 e confirmar o estado `Issued`.
5. Encerrar com o certificado pronto. ALB, Target Group, associação do
   certificado e publicação de HTTPS ficam para o vídeo 6.

O passo a passo do console está no
[roteiro do vídeo 5](../../../seasons/02-high-availability/docs/video-roadmap.md#vídeo-5--certificados-no-swarm-e-configuração-do-acm).

## Notas de narração

- DNS-01 não é incompatível com múltiplos Traefiks. O limite mostrado é a
  gestão independente do mesmo domínio por três processos ACME, que pode
  duplicar operações e consumir limites da autoridade certificadora.
- Em HTTP-01, o IP de origem do pedido de emissão não determina o destino
  da validação. O desafio precisa estar acessível em qualquer entrada usada
  pela CA. DNS fornece endereços, não encaminha o tráfego HTTP.
- Um Traefik pode emitir e renovar normalmente, mas sua falha interrompe
  a entrada até a recuperação. Existem outras soluções de emissão centralizada
  e distribuição de certificados; aqui escolhemos a integração com AWS.
- No quarto quadro, o ACM fornece o certificado ao listener; não é um salto
  no caminho da requisição. O Target Group é um agrupamento lógico de destinos,
  não outro proxy entre o ALB e os managers.
- O certificado do ALB protege cliente-ALB, não o HTTP interno. Restringir
  a porta 80 dos managers ao SG do ALB faz parte do vídeo 6. Retirar IPs
  públicos e migrar para subnets privadas são uma evolução posterior.
- No ACM usamos validação por CNAME, diferente do TXT do ACME DNS-01.
  Preservar os CNAMEs; a renovação gerenciada depende deles e do certificado
  em uso num serviço integrado. Ao final deste vídeo ele ainda não está
  associado ao ALB e a aplicação permanece em HTTP.

## Referências

- [Tipos de desafio ACME, Let's Encrypt](https://letsencrypt.org/docs/challenge-types/).
- [Certificate resolvers, Traefik](https://doc.traefik.io/traefik/v3.6/reference/install-configuration/tls/certificate-resolvers/overview/).
- [Solicitar certificado público, ACM](https://docs.aws.amazon.com/acm/latest/userguide/acm-public-certificates.html).
- [Validação DNS, ACM](https://docs.aws.amazon.com/acm/latest/userguide/dns-validation.html).
- [Renovação com validação DNS, ACM](https://docs.aws.amazon.com/acm/latest/userguide/dns-renewal-validation.html).
- [Listeners HTTPS, ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html).
- [IPs privados e segurança dos targets, ELB](https://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/infrastructure-security.html).
