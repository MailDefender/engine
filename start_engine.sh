#!/bin/sh

set -e

echo "########"
echo "Starting migrations..."

/utils/atlas migrate apply --env gorm --url $DATABASE_DNS

echo "Migration result : $?"
echo "#######"

echo "LOG_FILE=$LOG_FILE"

if [ -z $LOG_FILE ]
then
    echo "Starting engine without exporting logs"
     ./engine
else
    echo "Creating logs directory..."
    mkdir -p $(dirname $LOG_FILE)
    echo "Starting engine"
    ./engine > $LOG_FILE 2>&1
fi