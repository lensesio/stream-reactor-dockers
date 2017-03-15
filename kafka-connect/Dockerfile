FROM confluentinc/cp-kafka-connect:3.0.1

ARG STREAM_REACTOR_VERSION
ARG STREAM_REACTOR_COMPONENT
ARG CONFLUENT_VERSION
ARG DIR
ENV STREAM_REACTOR_VERSION ${STREAM_REACTOR_VERSION}
ENV STREAM_REACTOR_COMPONENT ${STREAM_REACTOR_COMPONENT}

#
# Liveliness probe
#
RUN apt-get update && apt-get install jq
COPY liveliness.sh /etc/confluent/liveliness.sh
#
# jmx exporter
#
COPY config.yaml /etc/jmx_exporter/config.yaml
RUN mkdir -p /opt/jmx_exporter && \
    wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.7/jmx_prometheus_javaagent-0.7.jar -O \
    /opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar
ENV KAFKA_JMX_OPTS="-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar=9102:/etc/jmx_exporter/config.yaml"
EXPOSE 9102

RUN cd /etc/${COMPONENT}/jars && \
    wget https://github.com/datamountaineer/stream-reactor/releases/download/v${STREAM_REACTOR_VERSION}/kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz && \
    tar -xvf kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz -C . && \
    rm -f -r kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz && \
    wget https://github.com/datamountaineer/kafka-connect-tools/releases/download/v1.0/kafka-connect-cli-1.0-all.jar
COPY entry-point.sh /etc/confluent/entry-point.sh
COPY templates/connector.properties.template /etc/confluent/docker/connector.properties.template
RUN apt-get update && apt-get install unzip
RUN chmod +x /etc/confluent/entry-point.sh
RUN wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64
RUN chmod +x /usr/local/bin/dumb-init
CMD ["dumb-init", "/etc/confluent/entry-point.sh"]
