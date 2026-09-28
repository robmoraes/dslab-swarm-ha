#!/bin/bash

set -o pipefail

i=1
while true; do
  printf "call=%02d " "$i"
  curl -fs --no-keepalive --connect-timeout 2 --max-time 3 \
    http://whoami.swarm.dslab.dev.br/metadata \
    | jq -r '
      (.host.hostname) as $host
      | ((.container.id // "")
      | if . == "" then $host else . end) as $container
      | "container=\($container[0:12]) node_eip=\(.cloud.public_ipv4) lb=\(.request.load_balancer.forwarded_for)"
      ' \
    || printf 'status=unavailable\n'
  i=$((i + 1))
  sleep 1
done
