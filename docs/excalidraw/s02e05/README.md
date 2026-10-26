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

O vídeo foi gravado com aproximadamente 20 minutos de explicação. A sequência
dos diagramas foi mantida, mas a configuração do ACM ficou para o vídeo 6.

1. Retrospectiva: Traefik único com HTTPS via HTTP-01, depois Traefik global
   em HTTP e cluster com três managers e três workers.
2. Justificativas para ALB + ACM: devolver HTTPS como requisito determinante
   e impedir acesso público direto às portas 80 e 443 dos nodes.
3. Apresentação de DNS-01, HTTP-01 e Traefik único, com os respectivos limites.
4. Encerramento dos desenhos com ACM + ALB + TG, fazendo o gancho para a
   implementação no vídeo 6.

Como bônus, foi feito o deploy da [stack do Excalidraw](../../../stacks/excalidraw.yml)
e a apresentação dos diagramas foi feita nela. Os arquivos preservam os quadros
usados na gravação; a previsão de ACM no vídeo 5 indicada no último quadro foi
transferida para o vídeo 6.

O conteúdo realizado está no
[roadmap do vídeo 5](../../../seasons/02-high-availability/docs/video-roadmap.md#vídeo-5--por-que-alb-e-acm-certificados-e-segurança).
O passo a passo do ACM e do ALB está no
[roteiro do vídeo 6](../../../seasons/02-high-availability/docs/video-roadmap.md#vídeo-6--acm-alb-e-target-groups-publicando-https).

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
- O certificado do ALB protege cliente-ALB, não o HTTP interno. No vídeo 6,
  a porta 80 dos managers ficará restrita ao SG do ALB, sem acesso público
  direto às portas 80/443 dos nodes. Retirar IPs públicos e migrar para
  subnets privadas são uma evolução posterior.
- A validação do ACM por CNAME e a emissão do certificado ficam para o
  vídeo 6. Ao final do vídeo 5, a aplicação permanece em HTTP e a arquitetura
  ACM + ALB + TG foi apenas apresentada nos desenhos.

## Referências

- [Tipos de desafio ACME, Let's Encrypt](https://letsencrypt.org/docs/challenge-types/).
- [Certificate resolvers, Traefik](https://doc.traefik.io/traefik/v3.6/reference/install-configuration/tls/certificate-resolvers/overview/).
- [Solicitar certificado público, ACM](https://docs.aws.amazon.com/acm/latest/userguide/acm-public-certificates.html).
- [Validação DNS, ACM](https://docs.aws.amazon.com/acm/latest/userguide/dns-validation.html).
- [Renovação com validação DNS, ACM](https://docs.aws.amazon.com/acm/latest/userguide/dns-renewal-validation.html).
- [Listeners HTTPS, ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/create-https-listener.html).
- [IPs privados e segurança dos targets, ELB](https://docs.aws.amazon.com/elasticloadbalancing/latest/userguide/infrastructure-security.html).
