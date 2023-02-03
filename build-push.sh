#!/usr/bin/env bash

NAME=modfin/securitytxt

docker buildx create --name securitytxt --use
docker buildx build --platform linux/amd64,linux/arm64 --tag=${NAME}:latest --push .
docker buildx rm securitytxt