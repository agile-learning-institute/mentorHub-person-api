# institute-person-api

This is a simple GoLang API that was written by a polyglot software engineer with the help of ChatGPT, with only a cursory understaqnding of the Go language and MongoDB. See [here](https://chat.openai.com/share/dcb8b738-7e73-40da-8b08-38024f1c9997) for the chat that was used to start.

[Here](./product-api-openapi.yaml) is the Swagger for the API

[Here](https://github.com/orgs/agile-learning-institute/repositories?q=institute-person&type=all&sort=name) are the repositories in the person triplet

[here](https://github.com/orgs/agile-learning-institute/repositories?q=institute&type=all&sort=name) are all of the repositories in the [Institute](https://github.com/agile-learning-institute/institute/tree/main) system

## Prerequisits

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go Language](https://go.dev/doc/install)

To run locally, first start Mongo and MongoExpress in docker desktop

```bash
docker compose up --detach
```

Once those containers have started, you can access the Express instance [here](http://localhost:8081).
If this is the first time your've started the containers you will want to create a "agile-learning-institute" database with a "people" collection

If you need to stop or start the containers you can use:

```bash
docker compose stop
```

or

```bash
docker compose start
```

When you are done you can remove the containers with:

```bash
docker compose down
```

NOTE: `docker compose down` will remove the database and collection

## Install dependencies and run the API

```bash
go get -u
go run main.go
```

## Test with CURL

Test Config

```bash
curl http://localhost:8080/api/config/

```

Test add a person

```bash
curl -X POST http://localhost:8080/api/person/ \
     -H "Content-Type: application/json" \
     -d '{"name":"Foo"}'

```

Test get a person

```bash
curl http://localhost:8080/api/person/[ID]

```

Test find all people with ID

```bash
curl http://localhost:8080/api/person/
```

Test update a person

```bash
curl -X PATCH http://localhost:8080/api/person/[ID] \
     -H "Content-Type: application/json" \
     -d '{"name":"Bar", "description":"Some long description"}'

```
