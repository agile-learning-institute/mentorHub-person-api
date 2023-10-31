mkdir api
cd api
curl https://raw.githubusercontent.com/agile-learning-institute/institute-person-api/main/src/docker/docker-compose.yaml > docker-compose.yaml
docker compose up --detach
cd ..
