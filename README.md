# test_task

## About

This service is a API which receives data for HTTP request and gives the result of request.

Work algorithm.
The server is waiting HTTP-request from client (curl, for example). In request's body there

should be message in JSON format. For example:

```json
POST /task
{
"method": "GET",
"url": "http://google.com",
"headers": {
"Authentication": "Basic
bG9naW46cGFzc3dvcmQ=",
....
}
}
```

Response:

```json
200 OK
{
 "id": <generated unique id>
}
```

Server forms valid HTTP-request to 3rd-party service with data from client's message and responses to client with JSON object:
Request:

```json
GET task/<taskId>
```

Response:

```json
200 OK
{
"id": <unique id>,
"status": "done/in_process/error/new"
"httpStatusCode": <HTTP status of 3rd-party service response>, "headers": {
<headers array from 3rd-party service response>
},
"length": <content length of 3rd-party service response> }
```

## Stack

- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit for Go
- [zap](https://github.com/uber-go/zap) - Logger
- [migrate](https://github.com/golang-migrate/migrate) - Database migrations. CLI and Golang library.
- [testify](https://github.com/stretchr/testify) - Testing toolkit
- [gomock](https://github.com/golang/mock) - Mocking framework
- [Docker](https://www.docker.com/) - Docker

## Running

In order to run this you need running database. You can run postgreqsl container as described in Instructions or if you have already running database, you can use it (just write its dsn in conf.yml).

### Instruction

1. Build database container:

   ```bash
   docker-compose up -d db
   ```

   You will have postgres container with database "task"

2. Now, you should migrate table creation schema into database container:

   ```bash
   migrate -path ./schema -database 'postgres://postgres:postgres@localhost:5432/task?sslmode=disable' up
   ```

   If you do not have migrate then install it:

   ```bash
       curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
       apt-get update
       apt-get install -y migrate
   ```

   You have table "tasks" in you database. To check it you can run:

   ```bash
   docker exec -it postgres psql  -U postgres -d task
   ```

   and then, look for a list of tables: `\d`

3. Run service:
   Choose correct 'data_source_name' in conf.yml depending on way of running service.

   - Run in command line:

   ```bash
       go run cmd/main.go
   ```

   - Run in docker container:

   ```bash
       docker-compose up -d test_task
   ```

You can test API by Postman or going to swagger [page](http://localhost:8000/documentation/index.html).

## Envs

List of all required envs. All envs must be set before build.

```docker
For db:
        POSTGRES_DB: postgres
        POSTGRES_USER: postgres
        POSTGRES_PASSWORD: secret
For test_task:
        APP_MODE: dev #also could be prod

```

## Maintainers

- Developed by @nnickie23

## Notes

In docker-compose.yml file on line 12:

```docker
    volumes:
      - ../postgresql:/var/lib/postgresql/data
```

This is written because postgresql must be outside of projects directory. If it is in project directory, there will be trouble building test-proxy

To regenerate swagger documentation:

```bash
    swag init -g cmd/main.go
```
