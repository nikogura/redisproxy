#!/bin/bash

redis-cli -h redis set foo foo
redis-cli -h redis set bar bar
redis-cli -h redis set wip wip
redis-cli -h redis set zoz zoz
redis-cli -h redis set ten 10

redisproxy run -c $SIZE -e $EXPIRATION -p $PORT -r $REDIS
