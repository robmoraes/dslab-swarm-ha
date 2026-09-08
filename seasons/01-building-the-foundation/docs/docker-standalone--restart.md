# docker `--restart`

| `--restart`      | Descrição curta                                                                            |
| ---------------- | ------------------------------------------------------------------------------------------ |
| `no`             | Nunca reinicia automaticamente (padrão).                                                   |
| `on-failure`     | Reinicia apenas se o processo terminar com erro (`exit code != 0`).                        |
| `on-failure:N`   | Igual ao `on-failure`, mas limita a `N` tentativas de reinício.                            |
| `always`         | Sempre tenta manter o container em execução, inclusive após reinício do Docker ou do host. |
| `unless-stopped` | Reinicia automaticamente, exceto se foi parado manualmente pelo administrador.             |

```bash
docker inspect -f '{{ .HostConfig.RestartPolicy.Name }}' nginx-port
docker update --restart always nginx-port
```
