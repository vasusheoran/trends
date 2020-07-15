@echo off 


echo "Starting ..."
set DOCKER_FILE="docker-compose.yml"

echo Docker File :: %DOCKER_FILE%


docker-compose  -f %DOCKER_FILE% down

docker-compose  -f %DOCKER_FILE% pull 

docker-compose   -f %DOCKER_FILE% up -d

start http://localhost/