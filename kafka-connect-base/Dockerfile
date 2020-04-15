ARG CONFLUENT_VERSION
FROM confluentinc/cp-kafka-connect-base:${CONFLUENT_VERSION}

ARG LENSES_CLI_VERSION
ARG CALCITE_LINQ4J_VERSION=1.12.0
ENV LENSES_CLI_VERSION=${LENSES_CLI_VERSION}
ENV CALCITE_LINQ4J_VERSION=${CALCITE_LINQ4J_VERSION}
ENV CONNECT_PLUGIN_PATH=/opt/lenses/lib
ENV APP_PROPERTIES_FILE=/etc/kafka-connect/connector.properties
ENV CONNECT_CLI=/opt/lenses/bin/connect-cli

RUN apt-get update && \
    apt-get install jq && \
    mkdir -p /opt/lenses/lib && \
    mkdir -p /opt/lenses/bin && \
    echo "progress = dot:giga" | tee /etc/wgetrc

COPY liveliness /opt/lenses/bin/liveliness
COPY entry-point /opt/lenses/bin/entry-point
COPY configure /opt/lenses/bin/configure
COPY config.yaml /etc/jmx_exporter/config.yaml

RUN mkdir -p /opt/jmx_exporter && \
    wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.7/jmx_prometheus_javaagent-0.7.jar -O \
    /opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar && \
    wget https://github.com/lensesio/kafka-connect-tools/releases/download/v1.0.6/connect-cli && \
    chmod +x /opt/lenses/bin/entry-point && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 && \
    chmod +x /usr/local/bin/dumb-init && \
    mv connect-cli /opt/lenses/bin && \
    chmod +x  /opt/lenses/bin/connect-cli && \
    wget https://github.com/lensesio/lenses-go/releases/download/${LENSES_CLI_VERSION}/lenses-cli-linux-amd64.tar.gz && \
    tar -xf lenses-cli-linux-amd64.tar.gz -C /opt/lenses/bin && \
    chmod +x /opt/lenses/bin/lenses-cli-linux-amd64 && \
    mv /opt/lenses/bin/lenses-cli-linux-amd64 /opt/lenses/bin/lenses-cli && \
    wget https://repo1.maven.org/maven2/org/apache/calcite/calcite-linq4j/${CALCITE_LINQ4J_VERSION}/calcite-linq4j-${CALCITE_LINQ4J_VERSION}.jar && \
    mv calcite-linq4j-${CALCITE_LINQ4J_VERSION}.jar /etc/kafka-connect/jars/ && \
    rm -f -r /usr/share/java/kafka-connect-activemq \
             /usr/share/java/kafka-connect-activemq \
             /usr/share/java/kafka-connect-ibmmq \
             /usr/share/java/kafka-connect-jms \
             /usr/share/java/kafka-connect-elasticsearch \
             /usr/share/java/kafka-connect-storage-common \
             /usr/share/java/monitoring-interceptors \
             /usr/share/java/rest-utils \
             /usr/share/java/confluent-hub-client \
             /usr/share/java/confluent-control-center \
             /usr/share/java/confluent-rebalancer && \
    unset CONNECT_CLI

ENV COMPONENT=kafka-connect
ENV KAFKA_JMX_OPTS="-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar=9102:/etc/jmx_exporter/config.yaml"
EXPOSE 9102
