#!/bin/bash

##./mover.sh | grep hadoop > mover.log

TABLES=

COUNTER=50

while [ $COUNTER -gt 1 ]
do
        COUNTER=$(($COUNTER -1))
        DAY=$(date +"%d" -d "-${COUNTER} day")
        MONTH=$(date +"%m" -d "-${COUNTER} day")
        YEAR=$(date +"%Y" -d "-${COUNTER} day")

        for TABLE in $(echo ${TABLES} | sed "s/,/ /g")
        do
                 hadoop fs -ls /user/hive/warehouse/compaction.db/${TABLE}_compacted/year=${YEAR}/month=${MONTH}/day=${DAY}/* 2>/dev/null
                 if [ $? -eq 0 ]
                 then
                        rm="hadoop fs -rmr /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/*"
                        cp="hadoop fs -cp /user/hive/warehouse/compaction.db/${TABLE}_compacted/year=${YEAR}/month=${MONTH}/day=${DAY}/* /data/in/${TABLE}/year=${YEAR}/month=${MONTH}/day=${DAY}/"
                        echo $rm
                        echo $cp
                fi
        done
done
