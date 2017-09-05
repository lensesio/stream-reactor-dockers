FROM confluentinc/cp-kafka-connect:3.2.0

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
    wget https://github.com/datamountaineer/kafka-connect-tools/releases/download/v1.0/kafka-connect-cli-1.0-all.jar && \
    mv kafka-connect-cli-1.0-all.jar /etc/datamountaineer/jars && \
    wget -O /usr/local/bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 && \
    chmod +x /usr/local/bin/dumb-init && \
    apt-get install -y libmysql-java && \
    cp /usr/share/java/mysql-connector-java.jar /etc/datamountaineer/jars && \
    wget http://central.maven.org/maven2/com/microsoft/sqlserver/mssql-jdbc/6.1.0.jre8/mssql-jdbc-6.1.0.jre8.jar && \
    cp mssql-jdbc-6.1.0.jre8.jar /etc/datamountaineer/jars

CMD ["dumb-init", "/etc/confluent/entry-point.sh"]
