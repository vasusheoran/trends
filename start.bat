@echo off 


echo "Starting ..."


docker-compose down

docker-compose pull 

docker-compose up -d

start http://localhost/