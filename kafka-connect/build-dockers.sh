#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STREAM_REACTOR_VERSION="1.2.0"
KAFKA_VERSION="2.0.0"
COMPONENT="kafka-connect"
REPO="datamountaineer"
export STREAM_REACTOR_VERSION
export KAFKA_VERSION

input="${DIR}/connectors.txt"

while IFS= read -r STREAM_REACTOR_COMPONENT
do
    echo "Building ${STREAM_REACTOR_COMPONENT} connector version ${STREAM_REACTOR_VERSION}"
    docker build \
        --build-arg STREAM_REACTOR_COMPONENT=${STREAM_REACTOR_COMPONENT} \
        --build-arg STREAM_REACTOR_VERSION=${STREAM_REACTOR_VERSION} \
        --build-arg KAFKA_VERSION=${KAFKA_VERSION} \
        --build-arg KAFKA_VERSION=${KAFKA_VERSION} \
        --build-arg COMPONENT=${COMPONENT} \
        -t ${REPO}/kafka-connect-${STREAM_REACTOR_COMPONENT}:${STREAM_REACTOR_VERSION} \
        -t ${REPO}/kafka-connect-${STREAM_REACTOR_COMPONENT} \
        -f Dockerfile .
        docker push ${REPO}/kafka-connect-${STREAM_REACTOR_COMPONENT}:${STREAM_REACTOR_VERSION}
        docker push ${REPO}/kafka-connect-${STREAM_REACTOR_COMPONENT}
done < "$input"


