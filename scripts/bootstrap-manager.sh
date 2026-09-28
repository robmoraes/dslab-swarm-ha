#!/bin/bash

# Prepara uma instancia Ubuntu para entrar como manager no Swarm existente.
# Uso na nova instancia: sudo ./scripts/bootstrap-manager.sh
# Instala o Docker, inicia o servico e adiciona ubuntu ao grupo docker.
#
# Depois do bootstrap:
# 1. No manager 1, execute: docker swarm join-token manager
# 2. Na nova instancia, execute com sudo o comando docker swarm join retornado.

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
newgrp docker
