# DataMountaineer Dockers

The images are based of Eneco's base image which in turn in base of Confluents Kafka Connect image, however the Eneco
image contains liveliness checks on the Connect worker running withn the Docker, this is meant for Kubernetes.

These images use Confluents docker utility belt to take all enviroment variables prefixed with `CONNECTOR` to create a
connector properties file. This properties file is then pushed via DataMountaineers Connect [CLI](https://github.com/datamountaineer/kafka-connect-tools) 
to the Connect API once it's up to start the connector.

The expected usecase is that the Connect Worker joins with other pods deployed via [Helm](https://helm.sh/) to form a 
Connect Cluster for a specific instance of one connector only. It can only post in a configuration for one type based on 
the enviroment variables.

For an awesome deployment app to deploy your landscape checkout Eneco's [Landscaper](https://github.com/Eneco/landscaper).

For example:
```bash
docker run \
    -e CONNECTOR_CONNECTOR_CLASS="com.datamountaineer.streamreactor.connect.cassandra.sink.CassandraSinkConnector" \
    -e CONNECTOR_NAME="cassandra-sink-orders" \
    -e CONNECTOR_TASKS_MAX="1" \
    -e CONNECTOR_TOPICS="orders-topic" \
    -e CONNECTOR_CONNECT_CASSANDRA_SINK_KCQL="INSERT INTO orders SELECT * FROM orders-topic" \
    -e CONNECTOR_CONNECT_CASSANDRA_CONTACT_POINTS="localhost" \
    -e CONNECTOR_CONNECT_CASSANDRA_PORT="9042" \
    -e CONNECTOR_CONNECT_CASSANDRA_KEY_SPACE="demo" \
    -e CONNECTOR_CONNECT_CASSANDRA_USERNAME="cassandra" \
    -e CONNECTOR_CONNECT_CASSANDRA_PASSWORD="cassandra" \
    --name datamountaineer\kafka-connect-cassandra:LATEST
```
