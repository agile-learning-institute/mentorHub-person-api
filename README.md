# institute-person-api

## Table of Contents

- [Overview](#overview)
- [Prerequisits](#prerequisits)
- [Run Containers for UI Development](#run-the-database-and-api-containers-locally-for-ui-development)
- [for API Engineers](#for-api-engineers)
  - [Using the Database Container](#using-the-database-container)
  - [Install dependencies and run the API locally](#install-dependencies-and-run-the-api-locally)
  - [Building and testing containers locally](#building-and-testing-the-container-locally)
- [Local API Testing with CURL](#local-api-testing-with-curl)
- [Observability and Configuration](#observability-and-configuration)
- [Backlog and Feature Branch info](#backlog-and-feature-branch-info)

## Overview

This is a simple GoLang API that was written by a polyglot software engineer with the help of ChatGPT, with only a cursory understaqnding of the Go language and MongoDB. See [here](https://chat.openai.com/share/dcb8b738-7e73-40da-8b08-38024f1c9997) for the chat that was used to start.

[Here](https://github.com/orgs/agile-learning-institute/repositories?q=institute&type=all&sort=name) are all of the repositories in the [Institute](https://github.com/agile-learning-institute/institute/tree/main) system

## Prerequisits

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
  - For Mac ```brew install docker```
- [Go Language](https://go.dev/doc/install)

### Optional

- [Mongo Compass](https://www.mongodb.com/try/download/compass) - if you want a way to look into the database

## Run the Database and API Containers locally for UI Development

```bash
curl https://raw.githubusercontent.com/agile-learning-institute/institute-person-api/main/src/docker/run-local-api.sh | /bin/bash
```

You can review the script at ./src/docker/run-local-api.sh
After a few seconds that command you should see something like this

```bash
 ✔ Network institute-person-ui_default                   Created
 ✔ Container institute-person-ui-institute-mongodb-1     Healthy
 ✔ Container institute-person-ui-institute-mongosh-1     Exited
 ✔ Container institute-person-ui-institute-person-api-1  Started
 ```

## For API Engineers

### Using the Database Container

If you want a local database, with test data preloaded, you can run the database containers independently. See [this repo](https://github.com/agile-learning-institute/institute-mongodb) for instructions on how to run the database containers.

### Install dependencies and run the API locally

If you have started the database separatly, you can run the API locally

```bash
go get 
go run main.go
```

## Building and Testing the container locally

If you have started the database container seperatly, you will need to stop it before testing this container. Use the following commands to do this.

```bash
cd db
docker compose down
cd ..
```

### Build the API container locally

```bash
./src/docker/docker-build.sh
```

### Start the Containers

```bash
cd ./src/docker
docker compose up --detach
```

### Stoping and Starting the containers without loosing data

```bash
cd ./src/docker
docker compose stop
docker compose start
```

### Restart the containers and Reseting the database

```bash
cd ./src/docker
docker compose down
docker compose up --deatch
```

## Local API Testing with CURL

### A word on ports

NOTE: If you are running the API from the command line with ```go run main.go``` the API will be served at port 8080, if you run the API in containers with ```docker compose up``` then it will be served at port 8081.
Adjust the following URI's accordingly.

### Test Health Endpoint

This endpoint supports the promethius monitoring standards for a healthcheck endpoint

```bash
curl http://localhost:8081/api/health/

```

### Test Config Endpoint

```bash
curl http://localhost:8081/api/config/

```

### Get Enumerators

```bash
curl http://localhost:8081/api/enums/

```

### Get Partner Names

```bash
curl http://localhost:8081/api/partners/

```

### Get Mentor Names

```bash
curl http://localhost:8081/api/mentors/

```

### Test find all people with IDs

```bash
curl http://localhost:8081/api/person/
```

### Test get a person

```bash
curl http://localhost:8081/api/person/[ID]

```

### Test add a person

```bash
curl -X POST http://localhost:8081/api/person/ \
     -d '{"name":"Foo", "description":"Some short description"}'

```

### Test update a person

```bash
curl -X PATCH http://localhost:8081/api/person/[ID] \
     -d '{"description":"Some long description"}'

```

## Observability and Configuration

The ```api/config/``` endpoint will return a list of configuration values. These values are either "defaults" or loaded from an Environment Variable, or found in a singleton configuration file of the same name. Environment Variables take precidence. The variable "CONFIG_FOLDER" will change the location of configuration files from the default of ```./```

The Dockerfile at the root of the project is a single-stage build that expects a linux native binary, and a text file called PATCH_LEVEL to exist, see [docker-build.sh](./docker-build.sh). The Dockerfile in /src/docker is a two stage build used for CI.  

The PATCH_LEVEL file that is located in the same folder as the executable should be populated by CI with the hash of the commit-to-main that triggers CI. This will be used on the Version number reported by the /api/config/ endpoint.

## Backlog and Feature Branch info

- [X] Shift dependency injection to Config object Bad Idea!
  - injection-refactor branch
  - [X] Person-store delivers DB Version to Config
  - [X] ConfigHandler injects config
  - [X] Add unit testing of config-handler
- [x] Improved Error Handling & testing (data dependency)
  - [x] All PersonStore timeout and errors logged and thrown
  - [x] Handlers to catch and return errors
- [x] Add attributes from database v1.1.Test
- [x] Add CID to HTTP Error Responses
- [x] Return full document after patch
- [x] Gorilla Promethius Health endpoint
- [x] Add breadcrumbs
- [x] Refactor Person as Simple Class, PersonStore to abstract mongo specific dependencies

- [x] Implement MongoStore
  - [x] Refactor config to use list of Store references objects {name, version, *Store}
  - [x] Refactor enum_store to mongo_store - move mongo-calls to mongo_store
  - [x] Refactor person_store to contain mongo_store
  - [x] Refactor enum_handlers into mongo_handler
  - [x] Add get/mentors endpoint with readOnlyHandler
  - [x] Add get/partners endpoint with readOnlyHandler
  - [x] Refactor Get /person and Get /people to use MongoStore and MongoHandler

- branch: ```store-housekeeping```
  - [x] Refactor database connect into a method and remove from constuctor
  - [x] Incorporate default query in MongoStore

- branch: ```store-default-query```
  - [x] refactor mongo_store FindMany
  - [x] Defult {$and: {$ne: {name: "VERSION"}, {$ne: {status: "Archived"}}}}
  - [x] $and to the parameters passed in constructor (i.e. mentor:ture)
  
- branch: ```Add-Unit-Testing```
  - [ ] Test config without connect/disconnect
  - [ ] Mock MongoStore for Unit Testing
  - [ ] Separate test of connect/disconnect with database
  - [ ] Build Tests (with ProfSynapse)

- branch: ```Post-Patch-Validation```
  - [ ] Person.IsValidPatch(id, bsonM) reflect hasSetter(signature, parmtype)
  - [ ] State Change Rules enforcement
    - Pending to Active or Drip
    - Active to Drip
    - Drip to Active
    - Anything to Archived
    - Prepend Name on Archive to avoid mystery dups
  - [ ] Referential Integrety enforcement??? (Research first)
  - [ ] NewPerson with BreadCrumb construction parameter
  - [ ] NewPerson with Default Values
  - [ ] Use NewPerson in Post processing

- branch: ```JWT-Authentication```
  - [ ] Add JWT authentication
  - [ ] update Breadcrumb constructor calls
