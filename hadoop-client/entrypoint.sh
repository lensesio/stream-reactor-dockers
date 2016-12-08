#!/usr/bin/env bash

echo "Getting yarn client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/YARN/clientConfig >> /etc/yarn-conf.zip || exit 1
(cd /etc && unzip yarn-conf.zip)

echo "Getting hive client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/HIVE/clientConfig >> /etc/hive-conf.zip || exit 1
(cd /etc && unzip hive-conf.zip)
