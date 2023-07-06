#!/bin/bash

cmd=$1

if [ $cmd == "up" ]; then
    go-migrate -database $DATABASE_URI -path migrations up
elif [ $cmd == "down" ]; then
    go-migrate -database $DATABASE_URI -path migrations down
else
    echo "No such command '$cmd'"
    exit 1
fi
