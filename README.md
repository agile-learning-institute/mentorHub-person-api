# institute-person-api

## Overview

This is a simple GoLang API that provides Get/Post/Patch services for docuements in the People collection. This API uses data from a [backing Mongo Database](https://github.com/agile-learning-institute/mentorHub-mongodb), and supports a [VueJS Single Page Appliaction.](https://github.com/agile-learning-institute/mentorHub-person-ui)

The OpenAPI specifications for the api can be found in the ``docs`` folder, and are served [here](https://agile-learning-institute.github.io/mentorHub-person-api/)

## Prerequisits

- [Mentorhub Developer Edition](https://github.com/agile-learning-institute/mentorHub/blob/main/mentorHub-developer-edition/README.md)
- [Go Language](https://go.dev/doc/install)

### Optional

- [Mongo Compass](https://www.mongodb.com/try/download/compass) - if you want a way to look into the database

## Install Go Dependencies
```bash
make install
```

## Run Unit Testing
```bash
make test
```

## Run the API locally 
```bash
make local
```
Serves up the API locally with a backign mongodb database, ctrl-c to exit

## Build the API Container
```bash
make container
```
This will build the new container, and start the mongodb and API container ready for testing. The test script ./test/test.sh is also run so you should see information about an inserted document. You will get a ``failed. Received HTTP code 000`` message if there are problems

## Generate Test Data
```bash
make generate
```
Generattes loads of test data, ctrl-c to exit

## API Testing with CURL
If you want to do more manual testing, here are the curl commands to use

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
curl http://localhost:8082/api/person/aaaa00000000000000000000

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

The [Dockerfile](./Dockerfile) uses a 2-stage build, and supports both amd64 and arm64 architectures. See [docker-build.sh](./src/docker/docker-build.sh) for details about how to build in the local architecture for testing, and [docker-push.sh] for details about how to build and push multi-architecture images.
