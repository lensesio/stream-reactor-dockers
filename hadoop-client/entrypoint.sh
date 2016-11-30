#!/usr/bin/env bash

# Pull cloudera config
echo "Getting Hadoop client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/$YARN_SERVICE_NAME/clientConfig >> /etc/hadoop-conf.zip
(cd /etc && unzip hadoop-conf.zip)
(cd /etc/hadoop-conf && sed -i 's/.engd.local//g' *)

echo "Getting Hive client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/$HIVE_SERVICE_NAME/clientConfig >> /etc/hive-conf.zip
(cd /etc && unzip hive-conf.zip)
(cd /etc/hive-conf && sed -i 's/.engd.local//g' *)


echo "Getting HBase client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/$HBASE_SERVICE_NAME/clientConfig >> /etc/hbase-conf.zip
(cd /etc && unzip hbase-conf.zip)
(cd /etc/hbase-conf && sed -i 's/.engd.local//g' *)

# echo "Getting spark client config"
# curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/SPARK/clientConfig >> /etc/spark-conf.zip
# (cd /etc && unzip spark-conf.zip)
# (cd /etc/spark-conf && sed -i 's/.engd.local//g' *)