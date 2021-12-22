@echo off 


echo "Starting ..."


docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml pull 

docker-compose -f docker-compose.yml up -d

start http://localhost/