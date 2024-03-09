# institute-person-api

## Table of Contents

- [Overview](#overview)
- [Prerequisits](#prerequisits)
- [Using the Database Container](#using-the-database-container)
- [Install Dependencies and Run](#install-dependencies-and-run-the-api-locally)
- [Build and Test the container](#building-and-testing-the-container-locally)
- [Local API Testing with CURL](#local-api-testing-with-curl)
- [Observability and Configuration](#observability-and-configuration)

## Overview

This is a simple GoLang API that provides Get/Post/Patch services for docuements in the People collection, as well as Get services for a number of related collections. This API uses data from a [backing Mongo Database](https://github.com/agile-learning-institute/mentorHub-mongodb), and supports a [VueJS Single Page Appliaction.](https://github.com/agile-learning-institute/mentorHub-person-ui)

The OpenAPI specifications for the api can be found in the ``docs`` folder, and are served [here](https://agile-learning-institute.github.io/mentorHub-person-api/)

## Prerequisits

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go Language](https://go.dev/doc/install)

### Optional

- [Mongo Compass](https://www.mongodb.com/try/download/compass) - if you want a way to look into the database

### Using the Database Container

If you want a local database, with test data preloaded, you can find instructions in the [mentorHub repository - docker configurations](https://github.com/agile-learning-institute/mentorHub/blob/main/docker-configurations/README.md#run-the-mongodb-backing-database)

### Install dependencies and run the API locally

If you have started the database separatly, you can run the API locally

```bash
cd ./src
go get 
go run main.go
```

### Building and Testing the container locally

You should build the container and test changes locally before making a pull request. You can use the build script below, and then [run curl tests](#local-api-testing-with-curl) to confirm the build.

```bash
./src/docker/docker-build.sh
```

## Local API Testing with CURL

### Test Health Endpoint

This endpoint supports the promethius monitoring standards for a healthcheck endpoint

```bash
curl http://localhost:8082/api/health/

```

### Test Config Endpoint

```bash
curl http://localhost:8082/api/config/

```

### Get Enumerators

```bash
curl http://localhost:8082/api/enums/

```

### Get Partner Names

```bash
curl http://localhost:8082/api/partners/

```

### Get Mentor Names

```bash
curl http://localhost:8082/api/mentors/

```

### Test find all people with IDs

```bash
curl http://localhost:8082/api/person/
```

### Test get a person

```bash
curl http://localhost:8082/api/person/[ID]

```

### Test add a person

```bash
curl -X POST http://localhost:8082/api/person/ \
     -d '{"userName":"Foo", "description":"Some short description"}'

```

### Test update a person

```bash
curl -X PATCH http://localhost:8082/api/person/aaaa00000000000000000021 \
     -d '{"description":"Some long description"}'

```

## Observability and Configuration

The ```api/config/``` endpoint will return a list of configuration values. These values are either "defaults" or loaded from an Environment Variable, or found in a singleton configuration file of the same name. Environment Variables take precidence. The variable "CONFIG_FOLDER" will change the location of configuration files from the default of ```./```

The ```api/health/``` endpoint is a Promethius Healthcheck endpoint.

The [Dockerfile](./src/docker/Dockerfile) uses a 2-stage build, and supports both amd64 and arm64 architectures. See [docker-build.sh](./src/docker/docker-build.sh) for details about how to build in the local architecture for testing, and [docker-push.sh] for details about how to build and push multi-architecture images.
