# SET BASE STAGE
ARG CONFLUENT_VERSION
FROM confluentinc/cp-kafka-connect-base:${CONFLUENT_VERSION}

# SET DOCKER ARGS
ARG SECRET_PROVIDER_VERSION
ARG CONNECT_CLI_VERSION
ARG CALCITE_LINQ4J_VERSION

# SET ENV
ENV CALCITE_LINQ4J_VERSION=${CALCITE_LINQ4J_VERSION:-"1.12.0"}
ENV CONNECT_CLI_VERSION=${CONNECT_CLI_VERSION:-"1.0.9"}
ENV SECRET_PROVIDER_VERSION=${SECRET_PROVIDER_VERSION:-"0.0.2"}

ENV CONNECT_PLUGIN_PATH=/opt/lenses/lib
ENV APP_PROPERTIES_FILE=/etc/kafka-connect/connector.properties
ENV CONNECT_CLI=/opt/lenses/bin/connect-cli

# UPDATE SYSTEM AND INSTALL JQ
RUN apt-get update -y && \
    apt-get -y --force-yes install jq

# SETUP DIRECTORIES
RUN mkdir -p /opt/lenses/lib && \
    mkdir -p /opt/lenses/bin && \
    mkdir -p /opt/lenses/kafka-connect/plugins \
    mkdir -p /opt/calcite \
    echo "progress = dot:giga" | tee /etc/wgetrc

# SETUP JMX EXPORTER
RUN mkdir -p /opt/jmx_exporter && \
    wget https://repo1.maven.org/maven2/io/prometheus/jmx/jmx_prometheus_javaagent/0.7/jmx_prometheus_javaagent-0.7.jar -O \
    /opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar

# SETUP CONNECT CLI
RUN wget https://github.com/lensesio/kafka-connect-tools/releases/download/v${CONNECT_CLI_VERSION}/connect-cli -O \
    /opt/lenses/bin/connect-cli \
    && chmod +x  /opt/lenses/bin/connect-cli

# SETUP SECRET PROVIER
RUN wget https://github.com/lensesio/secret-provider/releases/download/${SECRET_PROVIDER_VERSION}/secret-provider-${SECRET_PROVIDER_VERSION}-all.jar -O \
    /opt/lenses/kafka-connect/plugins/secret-provider-${SECRET_PROVIDER_VERSION}-all.jar

# SETUP CALCITE LING4J
RUN wget https://repo1.maven.org/maven2/org/apache/calcite/calcite-linq4j/${CALCITE_LINQ4J_VERSION}/calcite-linq4j-${CALCITE_LINQ4J_VERSION}.jar -O \
    /opt/calcite/calcite-linq4j-${CALCITE_LINQ4J_VERSION}.jar

# CLEANUP
RUN  rm -f -r /usr/share/java/kafka-connect-* \
             /usr/share/java/confluent-* \
             /usr/share/java/rest-utils-* \
        && unset CONNECT_CLI

# INSTALL DUMP INIT
RUN wget -O /usr/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64

# COPY LOCAL FILES
COPY liveliness /opt/lenses/bin/liveliness
COPY entry-point /opt/lenses/bin/entry-point
COPY config.yaml /etc/jmx_exporter/config.yaml
COPY kc-cli /opt/lenses/bin/kc-cli

RUN chmod +x /usr/bin/dumb-init /opt/lenses/bin/entry-point

ENV COMPONENT=kafka-connect

# CONFIGURE JMX AND EXPOSE JMX PORT
ENV KAFKA_JMX_OPTS="-javaagent:/opt/jmx_exporter/jmx_prometheus_javaagent-0.7.jar=9103:/etc/jmx_exporter/config.yaml"
EXPOSE 9103

# DEFINE ENTRYPOINT
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["/opt/lenses/bin/entry-point"]

