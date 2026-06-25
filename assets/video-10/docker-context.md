# Docker Context

Para facilitar execução de comandos Docker no servidor, vou configurar
o Docker Context em meu computador:

```bash
docker context create manager1 \
  --docker "host=ssh://ubuntu@carlosmoraesrodrigues.dev.br"
```

Para usar o contexto na sessão do terminal:

```bash
docker context use manager1
docker ps

# Retornar a sessão para o contexto padrão
docker context use default
docker ps
```

Para usar inline:

```bash
docker --context manager1 ps
```
