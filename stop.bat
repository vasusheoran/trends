@echo off 

echo Stopping running apps ...

docker-compose -f trends.docker-compose.yml down

pause