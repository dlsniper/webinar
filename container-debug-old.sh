#!/usr/bin/env bash

ARTIFACT_NAME="webinar"
CONTAINER_NAME="webinar"
CONTAINER_TAG="debug"

PROJECT_NAME='github.com/dlsniper/webinar'
PROJECT_DIR=${PWD}

CONTAINER_GOPATH='/go'
CONTAINER_PROJECT_DIR="${CONTAINER_GOPATH}/src/${PROJECT_NAME}"

docker run --rm \
        --net="host" \
        -v ${PROJECT_DIR}:${CONTAINER_PROJECT_DIR} \
        -e GOPATH=${CONTAINER_GOPATH} \
        -e CGO_ENABLED=0 \
        -w "${CONTAINER_PROJECT_DIR}" \
        golang:1.8.3 \
        go build -tags netgo -installsuffix netgo -ldflags "-X main.botVersion=${CONTAINER_TAG}" -o ${ARTIFACT_NAME} ${PROJECT_NAME}

docker build -f ${PROJECT_DIR}/Old-Dockerfile \
    -t ${CONTAINER_NAME}:${CONTAINER_TAG} \
    "${PROJECT_DIR}"
