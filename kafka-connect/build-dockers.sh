#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STREAM_REACTOR_VERSION="0.4.0"
KAFKA_VERSION="0.11.0.1"
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
        docker push datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT}:${STREAM_REACTOR_VERSION}
        docker push datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT}
done < "$input"


