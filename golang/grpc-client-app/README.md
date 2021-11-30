#### Build
D:
cd wsl\data\trends\golang\grpc-client-app
go build -o trends-client-app.exe main.go

#### Run
```trends-client-app.exe --server localhost:5001 --date 10-2-1 --high 17782 --close 17772.00 --low 17608.15 --symbol 1```
