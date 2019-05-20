#!/usr/bin/env bash
docker-compose up -d admin db
watchexec --restart --exts "go" --watch . "./start-app.sh"
