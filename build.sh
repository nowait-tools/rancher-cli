#!/usr/bin/env bash

set -ef -o pipefail

# BUILD
REVISION=$(git log --oneline -1)
VERSION=0.0.8
BUILD=$(date +%s)

DOCKER_REGISTRY="614339721584.dkr.ecr.ap-southeast-2.amazonaws.com"
DOCKER_IMAGE="rancher-cli"
REMOTE_TAG_PREFIX="${DOCKER_REGISTRY}/${DOCKER_IMAGE}"
VERSION_TAG=${DOCKER_IMAGE}:${VERSION}

DATE=$(date)

echo ${REVISION}
echo "**** BUILDING ${VERSION_TAG}+${BUILD} ****"

echo "TESTING"
echo "no tests"

echo "BUILDING IMAGE"
docker build --rm=true --disable-content-trust -t ${DOCKER_IMAGE} .
docker tag ${DOCKER_IMAGE} ${REMOTE_TAG_PREFIX}:latest
docker tag ${DOCKER_IMAGE} ${REMOTE_TAG_PREFIX}:${VERSION}

echo "PUSHING"
$(aws ecr get-login --region ap-southeast-2 --profile synergia --no-include-email)
docker push ${REMOTE_TAG_PREFIX}:${VERSION}

echo "**** COMPLETED ${DOCKER_IMAGE}:${VERSION}+${BUILD} ****"