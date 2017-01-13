#!/bin/bash

##./mover.sh | grep hadoop > mover.log

TABLES=sys_twitter_rage_raw,sys_trayport_contracts_proc,sys_trayport_products_proc,sys_trayport_raw_normalised_proc,sys_ebase_power_raw,sys_ebase_power_upscaled_proc,sys_ebase_power_validated_proc,sys_reuters_raw,sys_trayport_trades_proc,sys_twitter_rage_tweets_proc,sys_twitter_rage_users_proc

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
