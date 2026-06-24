#!/bin/bash

IMAGE_NAME="nginx:alpine"
CONTAINER_NAME="nginx-port"
HOST_PORT="8080"
HOST_VOLUME="/home/ubuntu/docker"

docker run \
  --detach \
  --name "$CONTAINER_NAME" \
  --publish "${HOST_PORT}:80" \
  --volume "${HOST_VOLUME}:/usr/share/nginx/html:ro" \
  "$IMAGE_NAME"
