#!/usr/bin/env bash
set -e

IMAGE_NAME=azatk/cube-http-gateway
IMAGE_VERSION=0.0.4

DIR_RELATIVE_PATH="$( dirname "$( which "$0" )" )"
DIR_ABSOLUTE_PATH="$(pwd)/$(basename "$DIR_RELATIVE_PATH")"

source "$DIR_ABSOLUTE_PATH"/build_app.sh
sudo docker build -t "$IMAGE_NAME":"$IMAGE_VERSION" "$DIR_ABSOLUTE_PATH"
sudo docker push "$IMAGE_NAME":"$IMAGE_VERSION"