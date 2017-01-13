# HDFS Compactor

This Docker installs an Impala shell, Hive and HDFS clients to compact the previous days partitions.

Hive is used for complex type tables, see ``USE_HIVE`` environment variable. Impala for flat structures. This is
because with Impala and complex types we'd need to explicity specify all columns.

## Environment Variables

| Name | Optional | Description |
|------|----------|-------------|
| USE_HIVE | No | List of tables with complex types. Hive will be used for compaction. |
| USE_IMPALA | No | List of table to compact using Impala. |
| SOURCE_DATABASE | No | Name of the source/primary database the table to compact belongs to. |
| IMPALA | No | The Impala deamon to connect to |
| OFFSET| No | Integer day offset to move backwards from current day for partition compaction selection. |
| CLOUDERA_MANAGER_URL | No | The Cloudera manager URL to download the hdfs and hive client configs from. |
| CLOUDERA_CLUSTER_NAME | No| The HDFS nameservice. |