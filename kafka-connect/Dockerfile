FROM eneco/connector-base:0.1.0

ARG STREAM_REACTOR_VERSION
ARG STREAM_REACTOR_COMPONENT
ARG CONFLUENT_VERSION
ARG DIR
ENV STREAM_REACTOR_VERSION ${STREAM_REACTOR_VERSION}
ENV STREAM_REACTOR_COMPONENT ${STREAM_REACTOR_COMPONENT}

RUN cd /etc/${COMPONENT}/jars && \
    wget https://github.com/datamountaineer/stream-reactor/releases/download/v${STREAM_REACTOR_VERSION}/kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz && \
    tar -xvf kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz -C . && \
    rm -f -r kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz
COPY entry-point.sh /etc/confluent/entry-point.sh
COPY templates/connector.properties.template /etc/confluent/docker/connector.properties.template
RUN apt-get update && apt-get install unzip
RUN chmod +x /etc/confluent/entry-point.sh
RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64
RUN chmod +x /usr/local/bin/dumb-init
CMD ["dumb-init", "/etc/confluent/entry-point.sh"]
