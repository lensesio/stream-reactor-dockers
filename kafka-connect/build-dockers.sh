#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
STREAM_REACTOR_VERSION="0.2.3"
CONFLUENT_VERSION="3.0.1"
COMPONENT="kafka-connect"
export STREAM_REACTOR_VERSION
export CONFLUENT_VERSION

input="${DIR}/connectors.txt"

while IFS= read -r STREAM_REACTOR_COMPONENT
do
    echo "Building ${STREAM_REACTOR_COMPONENT} connector version ${STREAM_REACTOR_VERSION}"
    docker build \
        --build-arg STREAM_REACTOR_COMPONENT=${STREAM_REACTOR_COMPONENT} \
        --build-arg STREAM_REACTOR_VERSION=${STREAM_REACTOR_VERSION} \
        --build-arg CONFLUENT_VERSION=${CONFLUENT_VERSION} \
        --build-arg CONFLUENT_VERSION=${CONFLUENT_VERSION} \
        -t datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT}:${STREAM_REACTOR_VERSION} \
        -t datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT} \
        -f Dockerfile .
        docker push datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT}:${STREAM_REACTOR_VERSION}
        docker push datamountaineer/kafka-connect-${STREAM_REACTOR_COMPONENT}
done < "$input"


