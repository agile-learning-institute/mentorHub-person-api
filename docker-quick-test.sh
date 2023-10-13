curl http://localhost:8081/api/config/
curl http://localhost:8081/api/person/
curl http://localhost:8081/api/person/651dfe6c13605cd1946273c2
curl -X POST http://localhost:8081/api/person/ -d '{"name":"Foo"}'
curl -X PATCH http://localhost:8081/api/person/651dfe6c13605cd1946273c2 -d '{"description":"Karmen, Grand PooBa Supreme!"}'