include commons.mk

SERVICE_NAME = dashboard
SERVICE_LANG = ts

USERNAME := $(shell whoami)
WORKSPACE := workspace

.PHONY: all build run

server:
# go build -ldflags '-w -s -extldflags -static' -o main golang/cmd/main.go
	docker run -it --rm --net=host -e GOOS=linux -e GOARCH=amd64 -e GOGC="" \
	-v /home/$(USERNAME)/build/go:/go \
	-v /home/$(USERNAME)/build/go_cache:/root/.cache/go-build \
	-v /home/$(USERNAME)/$(WORKSPACE)/git/trends:/app \
	-w /app/golang golang:1.17 /bin/bash \
	-c CGO_ENABLED=0 go build -ldflags '-w -s -extldflags -static' -o main cmd/main.go
# docker run --rm --net=host -e GOOS=linux -e GOARCH=amd64 -e GOGC="" \
# -v /Users/vasusheoran/go/darwin_amd64/go:/go \
# -v /Users/vasusheoran/go/darwin_amd64/go_cache:/root/.cache/go-build \
# -v /Users/vasusheoran/git/trends:/app \
# -w /app/golang golang:1.17 /bin/bash -c CGO_ENABLED=0 go build -ldflags '-w -s -extldflags -static' -o main cmd/main.go
	
dashboard:
	${NPM_DR_LINUX} "${NPM_BUILD_CMD}"

ui-run:
	${NPM_RUN_LINUX} "${NPM_PKG_STAGE_CMD}"

envoy-run:
	docker run -p 4200:4200 -p 49153:49153 --rm --name web -v /Users/vasusheoran/git/trends:/Users/vasusheoran/git/trends -w /Users/vasusheoran/git/trends/dashboard  -e ENVOY_UID=777 -e ENVOY_GID=777  envoyproxy/envoy:v1.18-latest 

ui-update:
	${NPM_UPDATE_ANGULAR} "${UPDATE_ANGULAR}"

image:
	docker-compose build
	doker-compose push
	
# protobufs-go: 
# 	cd golang/grpc/proto && protoc --go_out=../client --go_opt=paths=source_relative \
# 		--go-grpc_out=../client --go-grpc_opt=paths=source_relative \
# 		ticker.proto

# protobufs-ts: 
# 	protoc --plugin=protoc-gen-ts="/Users/vasusheoran/git/trends/dashboard/node_modules/.bin/protoc-gen-ts" \
#        --js_out="import_style=commonjs,binary:dashboard/src/app/generated" \
#        --ts_out="service=grpc-web:dashboard/src/app/generated" \
#        --proto_path="golang/grpc/proto"  golang/grpc/proto/ticker.proto

# protoc -I="golang/grpc/proto" golang/grpc/proto/ticker.proto \
# --js_out=import_style=commonjs,binary:dashboard/src/app/generated \
# --grpc-web_out=import_style=typescript,mode=grpcwebtext:dashboard/src/app/generated


gen-protobufs: protobufs-go protobufs-ts
# Target pattern to match any from parent
%: ;

# > sudo mv ~/Downloads/protoc-gen-grpc-web-1.2.1-darwin-x86_64 \
#     /usr/local/bin/protoc-gen-grpc-web
# Password:
# >  chmod +x /usr/local/bin/protoc-gen-grpc-web