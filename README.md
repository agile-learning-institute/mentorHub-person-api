# product-api

This is a simple GoLang API that was written by a polyglot software engineer with the help of ChatGPT, with only a cursory understaqnding of the Go language. See [here](https://chat.openai.com/share/dcb8b738-7e73-40da-8b08-38024f1c9997) for the chat that was used. At this point in time the API does not use any backing services, I will be adding a Mongo database next. (PS I have only a general understanding of Mongo)

To run locally, first start Mongo and MongoExpress in docker desktop

```bash
docker compose up --detach
```

The install dependencies and run the API

```bash
go get -u
go run main.go
```

## Test with CURL

Test Config

```bash
curl http://localhost:8080/api/config/

```

Test add a product

```bash
curl -X POST http://localhost:8080/api/product/ \
     -H "Content-Type: application/json" \
     -d '{"name":"Foo"}'

```

Test update a product

```bash
curl -X PATCH http://localhost:8080/api/product/123 \
     -H "Content-Type: application/json" \
     -d '{"name":"Bar", "description":"Some long description"}'

```
