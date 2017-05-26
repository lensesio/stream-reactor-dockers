#!/usr/bin/env bash

function push_config {
    # wait until local connect's REST API comes up
    until $CLI_CMD ps >>/dev/null
    do
        echo "Waiting for connect's rest API at $KAFKA_CONNECT_REST"
        sleep 1
    done
    echo "Pushing connector config..."
    $CLI_CMD run $CONNECTOR_NAME < $APP_PROPERTIES_FILE
    echo "done."
}
APP_PROPERTIES_FILE=/etc/config/connector.properties
CLI_JAR=/etc/datamountaineer/jars/kafka-connect-cli-1.0-all.jar

# cli expects this env var
export KAFKA_CONNECT_REST="http://$CONNECT_REST_ADVERTISED_HOST_NAME:$CONNECT_REST_PORT"
CLI_CMD="java -jar $CLI_JAR"
# Create
echo "Creating connector properties file"
mkdir /etc/config
dub template "/etc/confluent/docker/connector.properties.template" "$APP_PROPERTIES_FILE"
push_config &
# start connect using dumb-init to handle signals
exec /etc/confluent/docker/run
# EOF
