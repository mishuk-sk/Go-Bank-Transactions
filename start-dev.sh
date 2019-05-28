#!/usr/bin/env bash
docker-compose up -d admin newman
docker-compose build db
docker-compose up -d db
watchexec --restart --exts "go" --watch . "./start-app.sh"
