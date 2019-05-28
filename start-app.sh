#!/usr/bin/env bash
docker-compose build app
docker-compose build newman
docker-compose up app
docker-compose up -d newman