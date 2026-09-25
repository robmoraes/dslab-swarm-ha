# Cluster Docker Swarm HA na AWS

Este repositório reúne o material de apoio da série **Cluster Docker Swarm HA
na AWS**, publicada no canal do YouTube
[@DistributedSystemsLab](https://www.youtube.com/@DistributedSystemsLab).

A série acompanha a evolução incremental de um laboratório de Docker Swarm na
AWS: primeiro construindo uma base funcional e, depois, adicionando os
componentes e os testes necessários para alta disponibilidade.

## Temporadas

| Ordem | Temporada | Status | Escopo |
| --- | --- | --- | --- |
| 1 | [Construindo a Base](seasons/01-building-the-foundation/) | Concluída e publicada | AWS, EC2, Docker, Swarm, DNS, Traefik, WAF e TLS |
| 2 | [Alta Disponibilidade](seasons/02-high-availability/) | Em produção | Delegação DNS, managers, workers, ALB, ACM e testes de falha |

Cada temporada mantém seu próprio índice, documentação e artefatos. O
[changelog](CHANGELOG.md) registra as versões publicadas e a evolução do
material.

## Documentação

- [Padrão mínimo de tags para laboratórios AWS](docs/aws-tags.md)

## Organização

```text
.
├── docs/
│   └── aws-tags.md
└── seasons/
    ├── 01-building-the-foundation/
    │   ├── README.md
    │   ├── assets/
    │   └── docs/
    └── 02-high-availability/
        ├── README.md
        └── docs/
            └── video-roadmap.md
```
