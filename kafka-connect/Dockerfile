FROM confluentinc/cp-kafka-connect:3.3.0

ARG STREAM_REACTOR_VERSION
ARG STREAM_REACTOR_COMPONENT
ARG CONFLUENT_VERSION
ARG DIR
ARG COMPONENT
ENV STREAM_REACTOR_VERSION ${STREAM_REACTOR_VERSION}
ENV STREAM_REACTOR_COMPONENT ${STREAM_REACTOR_COMPONENT}
ENV CONFLUENT_VERSION ${CONFLUENT_VERSION}
ENV COMPONENT ${COMPONENT}
#
# Liveliness probe
#
RUN apt-get update && apt-get install jq
COPY liveliness.sh /etc/confluent/liveliness.sh
COPY entry-point.sh /etc/confluent/entry-point.sh
COPY templates/connector.properties.template /etc/confluent/docker/connector.properties.template
#copy in our own launch script so confluents doesn't clobber the classpath'
COPY launch /etc/confluent/docker/

#
# jmx exporter
#
COPY config.yaml /etc/jmx_exporter/config.yaml
RUN mkdir -p /opt/jmx_exporter && \
    wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.7/jmx_prometheus_javaagent-0.7.jar -O \
    /opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar
ENV KAFKA_JMX_OPTS="-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar=9102:/etc/jmx_exporter/config.yaml"
EXPOSE 9102

RUN mkdir -p /etc/datamountaineer/jars && \
    mkdir -p /etc/datamountaineer/bin && \
    wget https://github.com/datamountaineer/stream-reactor/releases/download/${STREAM_REACTOR_VERSION}/kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz && \
    tar -xvf kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.tar.gz -C . && \
    wget https://github.com/datamountaineer/kafka-connect-tools/releases/download/v1.0.2/connect-cli && \
    sed -i "s/<STREAM_REACTOR_COMPONENT>/$STREAM_REACTOR_COMPONENT/g" /etc/confluent/entry-point.sh && \
    sed -i "s/<STREAM_REACTOR_VERSION>/$STREAM_REACTOR_VERSION/g" /etc/confluent/entry-point.sh && \
    sed -i "s/<CONFLUENT_VERSION>/$CONFLUENT_VERSION/g" /etc/confluent/entry-point.sh && \
    chmod +x /etc/confluent/entry-point.sh && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 && \
    chmod +x /usr/local/bin/dumb-init && \
    mv kafka-connect-${STREAM_REACTOR_COMPONENT}-${STREAM_REACTOR_VERSION}-${CONFLUENT_VERSION}-all.jar /etc/datamountaineer/jars && \
    mv connect-cli /etc/datamountaineer/bin && \
    chmod +x  /etc/datamountaineer/bin/connect-cli

CMD ["dumb-init", "/etc/confluent/entry-point.sh"]
