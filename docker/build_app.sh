#!/usr/bin/env bash
set -e

APP_PATH=github.com/akaumov/cube-http-gateway

DIR_RELATIVE_PATH="$( dirname "$( which "$0" )" )"
DIR_ABSOLUTE_PATH="$(pwd)/$(basename "$DIR_RELATIVE_PATH")"
mkdir "$DIR_ABSOLUTE_PATH"/build -p

docker run --rm -v "$PWD":/go/src/"$APP_PATH" -v "$DIR_ABSOLUTE_PATH"/build:/build -w /go/src/"$APP_PATH"/cmd/http-gateway golang:1.10.3-alpine go build -x -v -o /build/app