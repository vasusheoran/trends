@echo off 


echo "Starting ..."


docker-compose -f trends.docker-compose.yml down

docker-compose -f trends.docker-compose.yml pull 

docker-compose -f trends.docker-compose.yml up -d

@REM start http://localhost/