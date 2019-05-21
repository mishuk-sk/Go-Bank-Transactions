# EPAM golang task
### Implemented small server with CRUD api
Actually, application runs in docker, so that there're several bash scripts provided to comfortably run docker containers during development. `start-dev.sh` starts 3 docker containers (including postgres, pgAdmin and also builds application container). It would just rebuild container on each `.go` file change. On other hand  `wait-for-it.sh` (was forked from [vishnubob/wait-for-it](https://github.com/vishnubob/wait-for-it)) runs inside container to start go application only after db is ready to accept connections.
