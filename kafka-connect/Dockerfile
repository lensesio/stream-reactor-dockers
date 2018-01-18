FROM confluentinc/cp-kafka-connect:3.3.1

ARG STREAM_REACTOR_VERSION
ARG STREAM_REACTOR_COMPONENT
ARG KAFKA_VERSION
ARG DIR
ARG COMPONENT
ENV STREAM_REACTOR_VERSION ${STREAM_REACTOR_VERSION}
ENV STREAM_REACTOR_COMPONENT ${STREAM_REACTOR_COMPONENT}
ENV KAFKA_VERSION ${KAFKA_VERSION}
ENV COMPONENT ${COMPONENT}
#
# Liveliness probe
#
RUN apt-get update && \
    apt-get install jq && \
    mkdir -p /etc/landoop/jars && \
    mkdir -p /etc/landoop/bin
COPY liveliness.sh /etc/landoop/bin/liveliness.sh
COPY entry-point.sh /etc/landoop/bin/entry-point.sh
COPY templates/connector.properties.template /etc/confluent/docker/connector.properties.template

#
# jmx exporter
#
COPY config.yaml /etc/jmx_exporter/config.yaml
RUN mkdir -p /opt/jmx_exporter && \
    wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.7/jmx_prometheus_javaagent-0.7.jar -O \
    /opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar
ENV KAFKA_JMX_OPTS="-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar=9102:/etc/jmx_exporter/config.yaml"
EXPOSE 9102

RUN wget https://github.com/landoop/stream-reactor/releases/download/${STREAM_REACTOR_VERSION}/kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${KAFKA_VERSION}-all.tar.gz && \
    tar -xvf kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${KAFKA_VERSION}-all.tar.gz -C . && \
    wget https://github.com/landoop/kafka-connect-tools/releases/download/v1.0.6/connect-cli && \
    chmod +x /etc/landoop/bin/entry-point.sh && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 && \
    chmod +x /usr/local/bin/dumb-init && \
    mv kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${KAFKA_VERSION}-all.jar /etc/landoop/jars && \
    mv connect-cli /etc/landoop/bin && \
    chmod +x  /etc/landoop/bin/connect-cli

CMD ["dumb-init", "/etc/landoop/entry-point.sh"]
