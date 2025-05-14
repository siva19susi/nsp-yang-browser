#!/bin/bash

set -e

# build script for frontend and backend containers

BASEDIR=$(realpath $(dirname "$0"))
PLATFORM_ARGS="--platform linux/amd64,linux/arm64"
PROXY_BUILD_ARGS="--build-arg http_proxy=${http_proxy} --build-arg HTTP_PROXY=${HTTP_PROXY} --build-arg https_proxy=${https_proxy} --build-arg HTTPS_PROXY=${HTTPS_PROXY}"

# build the automation panel container
docker build ${PROXY_BUILD_ARGS} ${PLATFORM_ARGS} -t sivasusi19/nsp-yang-browser:srx ${BASEDIR}
