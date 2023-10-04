# institute-person-api

This is a simple GoLang API that was written by a polyglot software engineer with the help of ChatGPT, with only a cursory understaqnding of the Go language and MongoDB. See [here](https://chat.openai.com/share/dcb8b738-7e73-40da-8b08-38024f1c9997) for the chat that was used to start.

[Here](./product-api-openapi.yaml) is the Swagger for the API

[Here](https://github.com/orgs/agile-learning-institute/repositories?q=institute-person&type=all&sort=name) are the repositories in the person triplet

[here](https://github.com/orgs/agile-learning-institute/repositories?q=institute&type=all&sort=name) are all of the repositories in the [Institute](https://github.com/agile-learning-institute/institute/tree/main) system

## Prerequisits

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go Language](https://go.dev/doc/install)
- [Mongo Compass](https://www.mongodb.com/try/download/compass) - if you want a way to look into the database

To run locally, you need to build the database container. Clone [this repo](https://github.com/agile-learning-institute/institute-person-db) and follow the instructions to build the container. Once that container built it will be run by the docker compose command below.

## Install dependencies and run the API locally

If you have started the database separatly, you can run the API locally

```bash
go get -u
go run main.go
```

## To Build the Container

```bash
GOOS=linux GOARCH=amd64 go build -o "institute-person-api" main.go
docker build . --tag institute-person-api
```

## To Run the API and Database Container 

```bash
docker compose up --detach
```

Note: If you see an error that looks like this

```bash
Error response from daemon: driver failed programming external connectivity on endpoint institute-person-api-institute-person-db-1 (f1517663e417de527d1ebf9d30a9ac21e4ca045d15bebb6297a79724f54536e9): Bind for 127.0.0.1:27017 failed: port is already allocated
```

You will need to stop the database container

- issue a ```docker compose down``` command
- cd to the database project and issue a ```docker compose down``` command
- cd back to the api project and try ```docker compose up``` again

## Stop and Start the containers without loosing data

```bash
docker compose stop
docker compose start
```

## Restart the containers (Reset the database)

```bash
docker compose down
docker compose up --deatch
```

## Test with CURL

NOTE: If you are running the API from the command line with ```go run main.go``` the API will be served at port 8080, 
if you run the API in containers with ```docker compose up``` then it will be served at port 8081. 
Adjust the following URI's accordingly.

Test Config Endpoint

```bash
curl http://localhost:8081/api/config/

```

Test add a person

```bash
curl -X POST http://localhost:8081/api/person/ \
     -H "Content-Type: application/json" \
     -d '{"name":"Foo"}'

```

Test get a person

```bash
curl http://localhost:8081/api/person/[ID]

```

Test find all people with ID

```bash
curl http://localhost:8081/api/person/
```

Test update a person

```bash
curl -X PATCH http://localhost:8081/api/person/[ID] \
     -H "Content-Type: application/json" \
     -d '{"name":"Bar", "description":"Some long description"}'

```
