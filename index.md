# DSLab | Cluster Docker Swarm on AWS

## Video 01 - Intro

- Apresentação
- Blueprint (ilustracao):
  - Fluxo desde o browser até o cluster
  - Registro.BR
  - AWS Route 53
  - AWS ACM
  - AWS ALB
  - Target Goups
  - Cluster
    - Traefik
    - Managers
    - Workers
    - Workloads

## Video 02 - IAM - Um usuário para o DSLab

- Requisitos
  - Conta free tier na AWS
  - Familiaridade com o Console da AWS
- Login com `admin`
  - Acesso ao `Identity and Access Management` (`IAM`)
  - Criação da Policy `DSLabProvisionerPolicy`
    - Permissão para IAM:MFA
  - Criação do Grupo `DSLabProvisioners`
  - Criação do User `dslab`
- Login com `dslab`
  - Ativação de MFA
- Login com `dslab` e MFA

## Primeira instancia e fundamentos

### Video 03 - EC2

- Intro
- Login com `admin`
  - IAM: Permissão para `EC2` na Policy
    - Policy: `DSLabProvisionerPolicy`
    - EC2:Full
- Blueprint (excalidraw):
  - EC2, EIP e SG
  - Acesso SSH (22) e HTTP (8080)
- EC2
  - SG: Criar Security Group
    - Inbound Rule SSH 22 MyIP
  - Criar Instance
  - EIP: Associar Elastic IP
- Acessar Instance via SSH

### Video 04 - Docker Bootstrap

- Acessar Instance via SSH
- Instalar e configurar Docker Community Edition
- Rodar Container de teste
- Parar o container e verificar

### Video 05 - Criar uma `Docker Image` para testar os Containers

Nesse vídeo vamos criar uma imagem docker e publicar no Docker Hub.

A imagem deve conter um programa para testar os containers que serão
executados em nossa EC2 Instance, desde os containers standalone até
os containers orquestrados em múltiplos nodes no futuro.

Vamos usar o Codex.

- Criar mais containers nas portas:
  - 8081
  - 8082
  - 8083

### Video 06 - Docker Swarm

- Iniciar Docker Swarm
- Prepara a stack para a API whoami
- Testar no browser
- Testar com script de carga

## Um domínio para os serviços

Ainda nessa fase de introdução e fortalecimento de fundamentos, vamos
usar um domínio já delegado em minha conta AWS para resolver o `EIP`
da `EC2 Instance` que estamos usando no laboratório.
`carlosmoraesrodrigues.dev.br` será o domínio.

E também vamos resolver o problema de ficar publicando portas para os
containers para conseguir acessar de fora.

### Video 07 Sendo amigável com humanos

- Blueprint: Problema com acesso usando IP
- Configuração do domínio carlosmoraesrodrigues.dev.br;
- Teste de acesso aos containers da instancia usando o `http://domain:port`
  no lugar de `IP:port`.

### Video 08 Um proxy para a borda

- Blueprint: Problema de gerenciamento de muitas portas publicadas;
- Route 53: criar wildcard
- EC2 SG: abrir porta 80, e derrubar o resto;
- Traefik: Edge Reverse Proxy para todo tráfego controlar.

### Video 09 Uma Stack, dois Services; API, WAF

- Blueprint: Dois Services com dependencia na mesma stack.
- Ajuste na Stack `whoami` incluindo um Service `Web Firewall Application`
  (WAF).

### Video 10 Certificado SSL/TLS com Traefik e Let's Encrypt

- EC2 SG: abrir porta 443;
- Traefik Stack: Configurar para geração de certificado com Let's
  Encrypt HTTP Challenge;

### Video 10 Encerramento dos Fundamentos
