# institute-people-api

This is a simple GoLang API that was written by a polyglot software engineer with the help of ChatGPT, with only a cursory understaqnding of the Go language. See [here](https://chat.openai.com/share/dcb8b738-7e73-40da-8b08-38024f1c9997) for the chat that was used.

To run locally, first start Mongo and MongoExpress in docker desktop

```bash
docker compose up --detach
```

Then install dependencies and run the API

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

Test update a person

```bash
curl -X PATCH http://localhost:8080/api/person/[ID] \
     -H "Content-Type: application/json" \
     -d '{"name":"Bar", "description":"Some long description"}'

```
