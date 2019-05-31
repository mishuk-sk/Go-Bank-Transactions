# EPAM golang task
### Implemented small server with CRUD api
Actually, application runs in docker, so that there're several bash scripts provided to comfortably run docker containers during development. `start-dev.sh` starts 4 docker containers (including postgres, pgAdmin, newman for running postman tests and also builds application container). It would just rebuild container and run newman tests on each `.go` file change. On other hand  `wait-for-it.sh` (was forked from [vishnubob/wait-for-it](https://github.com/vishnubob/wait-for-it)) runs inside container to start go application only after db is ready to accept connections.
Also, after running tests, newman outputs results into `testing_env/Test_output.html`, that presents test result in readable for human form.
By the way, there are several questions on package organization, testing environment, etc.
Testing looks REALLY messy. Like this, always restarting:
[code sample](https://github.com/mishuk-sk/Go-Bank-Transactions/blob/66fefbd04bed6ecc39f1d548d160272ba72abce6/start-dev.sh#L1-L5)

So that question appears. May be you can suggest any approach to do it more readable and using less resources?


Another question appears about packaging. Actually, I'm doing it wrong, I think. In [main.go](https://github.com/mishuk-sk/Go-Bank-Transactions/blob/master/main.go) I'm importing [handlers.go](https://github.com/mishuk-sk/Go-Bank-Transactions/blob/master/handlers/handlers.go), whether it just impliments functionality needed ONLY in main.go (so that, it's not reusable). 
Possible solution is to do like in my `handlers` package, where logical code parts are just organized in different files. Actually, I'm asking for feedback on this topic.

Also there's a question about organizing handleFunctions like it is in code. Is it the best way to do it? Could I do better?:)
