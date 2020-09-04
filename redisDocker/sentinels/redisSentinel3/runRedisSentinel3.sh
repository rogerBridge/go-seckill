#!/usr/bin/env bash
# docker stop redis && docker rm redis
# docker network create redisStore
#  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf  \
docker run -d \
  -v $PWD/sentinel.conf:/usr/local/etc/redis/sentinel.conf \
  --name sentinel3  \
  -p 127.0.0.1:26383:26379 \
  --network=redisStore  \
  --network-alias=sentinel3 \
  redis:latest redis-sentinel /usr/local/etc/redis/sentinel.conf