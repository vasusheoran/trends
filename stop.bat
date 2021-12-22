@echo off 

echo Stopping running apps ...

docker-compose -f docker-compose.yml down

pause