#!/usr/bin/env bash

export OFFSET=1
export IMPALA=localhost
export SOURCE_DATABASE=default
export COMPACTION_LOCATION=/data/compaction
export USE_HIVE=
export USE_IMPALA=

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
    DROP TABLE IF EXISTS ${COMPACTED_TABLE};
    CREATE TABLE IF NOT EXISTS ${COMPACTED_TABLE} LIKE ${SOURCE_DATABASE}.${TABLE};

    INSERT INTO TABLE compaction.${COMPACTED_TABLE} PARTITION (year,month,day)
    SELECT * FROM ${SOURCE_DATABASE}.${TABLE} WHERE year<=${YEAR} AND month<=${MONTH} AND day<=${DAY}"

    $SETCOLOR_NORMAL
    hive -e "${SQL}"
    if [ $? -eq 0 ]
    then
        $SETCOLOR_SUCCESS
        echo "Successfully compacted table ${TABLE}"
        $SETCOLOR_NORMAL
        impala-shell -i ${IMPALA} -q "invalidate metadata ${TABLE}" -d ${SOURCE_DATABASE}
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
    DROP TABLE IF EXISTS ${COMPACTED_TABLE};
    CREATE TABLE ${COMPACTED_TABLE} LIKE ${SOURCE_DATABASE}.${TABLE};

    INSERT OVERWRITE TABLE compaction.${COMPACTED_TABLE} PARTITION (year,month,day)
    SELECT * FROM ${SOURCE_DATABASE}.${TABLE} WHERE year<='${YEAR}' AND month<='${MONTH}' AND day<='${DAY}'"

    $SETCOLOR_NORMAL
    impala-shell -i ${IMPALA} -q "${SQL}"
    if [ $? -eq 0 ]
    then
        $SETCOLOR_SUCCESS
        echo "Successfully compacted table ${TABLE}"
        $SETCOLOR_NORMAL
        impala-shell -i ${IMPALA} -q "invalidate metadata ${TABLE}" -d ${SOURCE_DATABASE}
    else
        $SETCOLOR_FAILURE
        echo "Failure compacting table ${TABLE}."
    fi

    $SETCOLOR_NORMAL
done
