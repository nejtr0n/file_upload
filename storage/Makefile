SHELL:=/bin/bash -O globstar

build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/srv pkg/srv.go
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bin/web pkg/web.go

buf:
	protoc --proto_path=$$GOPATH/src:./proto --go_out=./proto --micro_out=./proto proto/*.proto