@echo off 

echo Stopping running apps ...

set DOCKER_FILE="docker-compose.yml"

docker-compose  -f %DOCKER_FILE% down

pause