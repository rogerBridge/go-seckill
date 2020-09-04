#!/usr/bin/env bash
# docker stop redis && docker rm redis
# docker network create redisStore
#  -v $PWD/redis.conf:/usr/local/etc/redis/redis.conf  \
docker run -d \
  -v $PWD/sentinel.conf:/usr/local/etc/redis/sentinel.conf \
  --name sentinel2  \
  -p 127.0.0.1:26382:26379 \
  --network=redisStore  \
  --network-alias=sentinel2 \
  redis:latest redis-sentinel /usr/local/etc/redis/sentinel.conf