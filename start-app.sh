#!/usr/bin/env bash
docker-compose build app
docker-compose up app
docker-compose restart newman