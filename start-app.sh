#!/usr/bin/env bash
docker-compose build app
if [!"$(docker ps -a | grep newman )"]
then
    docker-compose restart newman
else
    docker-compose up -d newman
fi
docker-compose up app
