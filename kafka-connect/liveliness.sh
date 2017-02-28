#!/usr/bin/env bash
#
# Usage:
# ./liveness.sh [--test <status> <worker_id>]
#
# --test  Run the probe in test mode with the provided <status> and <worker_id>.
#
# Exit code: 0 if no tasks have failed, 1 if one or more have failed.
#
if [ "$1" == "--test" ]
then
  STATUS=$2
  WORKER_ID=$3
else
  STATUS="$(curl http://$CONNECT_REST_ADVERTISED_HOST_NAME:$CONNECT_REST_PORT/connectors/$CONNECTOR_NAME/status)"
  WORKER_ID="\"$CONNECT_REST_ADVERTISED_HOST_NAME:$CONNECT_REST_PORT\""
fi
if [ "$STATUS" == null ]
then
  echo "Failed to query task status"
  exit 1
fi
# If the connector isn't posted, there will be an error_code
if [ $(echo $STATUS | jq 'has("error_code")') == true ]
then
  echo "Failed to query task status"
  echo $STATUS
  exit 1
fi
JQ_QUERY="[.tasks[] | select(.state == \"FAILED\" and .worker_id == $WORKER_ID) | .id] | length"
FAILED_TASKS=$(echo $STATUS | jq "$JQ_QUERY")
echo "$FAILED_TASKS task(s) failed on $WORKER_ID"
if [ $FAILED_TASKS != 0 ]
then
  echo "$STATUS"
  exit 1
else
  exit 0
fi