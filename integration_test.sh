#!/bin/bash

ENDPOINT=$1
PORT=$2

OUTPUT=$(curl -s http://localhost:$PORT/$ENDPOINT)

if [[ "$OUTPUT" =~ "$ENDPOINT" ]]; then
    echo "... OK"
else
    echo "... Not OK"
fi