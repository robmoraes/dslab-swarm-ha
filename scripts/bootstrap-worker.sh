#!/bin/bash

# Prepara uma instancia Ubuntu para entrar como worker no Swarm existente.
# Copie este arquivo para a instancia e execute: sudo bash bootstrap-worker.sh
# Instala o Docker, inicia o servico e adiciona ubuntu ao grupo docker.
#
# Depois do bootstrap:
# 1. Saia do SSH e reconecte para usar docker sem sudo.
# 2. Em um manager existente, execute: docker swarm join-token worker
# 3. No worker, execute o comando docker swarm join retornado pelo manager.
#    Use o IP privado do manager e --advertise-addr com o IP privado do worker.
# O script apenas instala o Docker; a entrada no cluster e feita manualmente.

set -e

apt-get update
apt-get install -y ca-certificates curl

install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

. /etc/os-release
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu ${UBUNTU_CODENAME:-$VERSION_CODENAME} stable" \
  > /etc/apt/sources.list.d/docker.list

apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
systemctl enable --now docker
usermod -aG docker ubuntu
