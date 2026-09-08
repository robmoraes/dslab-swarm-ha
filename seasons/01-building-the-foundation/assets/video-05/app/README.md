# dslab-whoami

API em Go para observar execucoes de containers em EC2, Docker standalone, Docker Swarm e, posteriormente, ALB.

## Endpoints

- `GET /` ou `GET /metadata`: snapshot completo da requisicao, runtime, host, container, rede e EC2 IMDS.
- `GET /headers`: headers recebidos, com valores sensiveis mascarados.
- `GET /healthz` e `GET /readyz`: health checks simples.
- `GET /version`: versao e dados de build.
- `GET /slow?ms=1000`: resposta com atraso controlado para testar timeout e distribuicao.

Use `?pretty=1` para JSON indentado.

## Variaveis uteis

- `PORT`: porta HTTP, padrao `8080`.
- `LISTEN_ADDR`: endereco completo de bind, sobrescreve `PORT`.
- `APP_NAME`: nome exposto na resposta, padrao `dslab-whoami`.
- `APP_ENV`: ambiente logico, padrao `lab`.
- `AWS_EC2_METADATA`: `auto`, `enabled` ou `disabled`, padrao `auto`.
- `AWS_EC2_METADATA_DISABLED=true`: desativa coleta EC2 IMDS.
- `AWS_EC2_METADATA_TIMEOUT`: timeout da coleta IMDS, padrao `250ms`.

Para Docker Swarm, exponha metadados com templates de env no service, por exemplo:

```sh
docker service create \
  --name dslab-api \
  --replicas 3 \
  --publish published=8080,target=8080 \
  --env NODE_HOSTNAME='{{.Node.Hostname}}' \
  --env SERVICE_NAME='{{.Service.Name}}' \
  --env TASK_ID='{{.Task.ID}}' \
  --env TASK_NAME='{{.Task.Name}}' \
  --env TASK_SLOT='{{.Task.Slot}}' \
  robmoraes/dslab-whoami:latest
```

## Testes com curl

Suba um container standalone:

```sh
docker run --rm -p 8080:8080 robmoraes/dslab-whoami:latest
```

Consulte o snapshot completo:

```sh
curl -s http://localhost:8080/metadata?pretty=1
```

Teste o health check:

```sh
curl -i http://localhost:8080/healthz
```

Dispare varias requisicoes para perceber qual container respondeu. O exemplo abaixo gera uma linha por resposta, facilitando enxergar quando cada chamada foi resolvida por um container diferente.

```sh
for i in $(seq 1 10); do
  printf "call=%02d " "$i"
  curl -s --no-keepalive http://localhost:8080/metadata \
    | jq -r '(.host.hostname) as $host | ((.container.id // "") | if . == "" then $host else . end) as $container | "seq=\(.request.sequence) container=\($container[0:12]) host=\($host) task_slot=\(.container.swarm.task_slot // "-") task_id=\((.container.swarm.task_id // "-")[0:12]) remote=\(.request.remote_address)"'
done
```

Exemplo de saida quando o balanceamento distribui entre replicas:

```text
call=01 seq=1 container=8ffab8daba60 host=8ffab8daba60 task_slot=1 task_id=x5z2n8d1k4p9 remote=10.0.1.20:54122
call=02 seq=1 container=3bb2a31291de host=3bb2a31291de task_slot=2 task_id=f7h3m1q8r6w2 remote=10.0.1.20:54128
call=03 seq=1 container=91c0f4d76a11 host=91c0f4d76a11 task_slot=3 task_id=n9v6b2c5t1y8 remote=10.0.1.20:54134
```

Se nao tiver `jq`, use o JSON completo:

```sh
for i in $(seq 1 5); do
  curl -s "http://localhost:8080/metadata?pretty=1"
done
```

Simule headers que normalmente chegam via proxy ou ALB:

```sh
curl -s "http://localhost:8080/metadata?pretty=1" \
  -H "X-Forwarded-For: 203.0.113.10, 10.0.0.25" \
  -H "X-Forwarded-Proto: https" \
  -H "X-Forwarded-Port: 443" \
  -H "X-Amzn-Trace-Id: Root=1-65f00000-0123456789abcdef01234567"
```

Veja apenas os headers recebidos pela API:

```sh
curl -s "http://localhost:8080/headers?pretty=1" \
  -H "Authorization: Bearer segredo" \
  -H "X-Request-Id: teste-local-001"
```

Teste atraso controlado para observar timeout, concorrencia e distribuicao de carga:

```sh
curl -s "http://localhost:8080/slow?ms=2000&pretty=1"
```

Em EC2 acessando por IP e porta, troque `localhost` pelo endereco da instancia:

```sh
curl -s "http://EC2_PUBLIC_IP:8080/metadata?pretty=1"
```

## Build

```sh
make test
make docker-build
make push
```
