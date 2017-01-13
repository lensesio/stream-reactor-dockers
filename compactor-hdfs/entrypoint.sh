#!/usr/bin/env bash

# Pull cloudera config
echo "Getting hadoop client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/YARN/clientConfig >> /etc/yarn-conf.zip
(cd /etc && unzip yarn-conf.zip)
(cd /etc/hadoop-conf && sed -i 's/.engd.local//g' *)

echo "Getting hive client config"
curl $CLOUDERA_MANAGER_URL/api/v12/clusters/$CLOUDERA_CLUSTER_NAME/services/HIVE/clientConfig >> /etc/hive-conf.zip
(cd /etc && unzip hive-conf.zip)
(cd /etc/hive-conf && sed -i 's/.engd.local//g' *)

#Get the current day and adjust for the offset
DAY=$(date +"%d" -d "-${OFFSET} day")
MONTH=$(date +"%m" -d "-${OFFSET} day")
YEAR=$(date +"%Y" -d "-${OFFSET} day")

# Command to set the color to SUCCESS (Green)
SETCOLOR_SUCCESS="echo -en \\033[1;32m"
# Command to set the color to FAILED (Red)
SETCOLOR_FAILURE="echo -en \\033[1;31m"
# Command to set the color back to normal
SETCOLOR_NORMAL="echo -en \\033[0;39m"

for TABLE in $(echo ${USE_HIVE} | sed "s/,/ /g")
do
    $SETCOLOR_SUCCESS
    COMPACTED_TABLE="${TABLE}_compacted"

    SQL="
    SET hive.exec.compress.output=true;
    SET hive.merge.mapredfiles=true;
    SET hive.hadoop.supports.splittable.combineinputformat=true;
    SET parquet.compression=SNAPPY;
    SET hive.merge.mapfiles=true;
    SET hive.exec.dynamic.partition.mode=nonstrict;

    CREATE DATABASE IF NOT EXISTS compaction;
    USE compaction;
    DROP TABLE IF EXISTS compaction.${COMPACTED_TABLE};
    CREATE TABLE IF NOT EXISTS compaction.${COMPACTED_TABLE} LIKE ${SOURCE_DATABASE}.${TABLE};

    INSERT INTO TABLE compaction.${COMPACTED_TABLE} PARTITION (year,month,day)
    SELECT * FROM ${SOURCE_DATABASE}.${TABLE} WHERE year=${YEAR} AND month=${MONTH} AND day=${DAY}"

    echo "Executing HIVE INSERT OVERWRITE for compaction.${COMPACTED_TABLE} PARTITION(year=${YEAR}, month=${MONTH}, day=${DAY})"
    $SETCOLOR_NORMAL
    beeline -u jdbc:hive2://engdcdr09502.engd.local:10000/default -n "" -p "" --fastConnect=true --color=true -e "${SQL}"

    if [ $? -eq 0 ]
    then
        $SETCOLOR_SUCCESS
        echo "Successfully compacted table ${TABLE}"
        echo "Copying data back to primary table"
        hadoop fs -rm -r /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/*
        hadoop fs -cp /user/hive/warehouse/compaction.db/${COMPACTED_TABLE,,}/year=${YEAR}/month=${MONTH}/day=${DAY}/* /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/
    
        $SETCOLOR_NORMAL
        impala-shell -i ${IMPALA} -q "INVALIDATE METADATA ${TABLE}; COMPUTE INCREMENTAL STATS ${TABLE}" -d ${SOURCE_DATABASE}
    else
        $SETCOLOR_FAILURE
        echo "Failure compacting table ${TABLE}."
    fi
    $SETCOLOR_NORMAL
done

$SETCOLOR_NORMAL

for TABLE in $(echo ${USE_IMPALA} | sed "s/,/ /g")
do
    $SETCOLOR_SUCCESS
    COMPACTED_TABLE="${TABLE}_compacted"

    SQL="
    CREATE DATABASE IF NOT EXISTS compaction;
    USE compaction;
    DROP TABLE IF EXISTS compaction.${COMPACTED_TABLE};
    CREATE TABLE compaction.${COMPACTED_TABLE} LIKE ${SOURCE_DATABASE}.${TABLE};

    INVALIDATE METADATA compaction.${COMPACTED_TABLE};
    INVALIDATE METADATA ${SOURCE_DATABASE}.${TABLE};

    INSERT OVERWRITE TABLE compaction.${COMPACTED_TABLE} PARTITION (year,month,day)
    SELECT * FROM ${SOURCE_DATABASE}.${TABLE} WHERE year='${YEAR}' AND month='${MONTH}' AND day='${DAY}'"

    echo "Executing Impala INSERT OVERWRITE for compaction.${COMPACTED_TABLE} PARTITION(year=${YEAR}, month=${MONTH}, day=${DAY})"
    $SETCOLOR_NORMAL
    impala-shell -i ${IMPALA} -q "${SQL}"
    if [ $? -eq 0 ]
    then
        $SETCOLOR_SUCCESS
        echo "Successfully compacted table ${TABLE}"
        echo "Copying data back to primary table"
        hadoop fs -rm -r /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/*
        hadoop fs -cp /user/hive/warehouse/compaction.db/${COMPACTED_TABLE,,}/year=${YEAR}/month=${MONTH}/day=${DAY}/* /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/
        $SETCOLOR_NORMAL
        impala-shell -i ${IMPALA} -q "INVALIDATE METADATA ${TABLE}; COMPUTE INCREMENTAL STATS ${TABLE}" -d ${SOURCE_DATABASE}
    else
        $SETCOLOR_FAILURE
        echo "Failure compacting table ${TABLE}."
    fi

    $SETCOLOR_NORMAL
done
