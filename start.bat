@echo off 


echo "Starting ..."

set "pwd=%cd%"

cd golang\client-app\http && go build -o trends-client-app.exe main.go

echo %pwd%

cd %pwd%

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml pull 

docker-compose -f docker-compose.yml up -d

start http://localhost/
