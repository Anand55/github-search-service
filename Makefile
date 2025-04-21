PROTO_FILE=api/githubsearch.proto
GO_PKG_PATH=github-search-service/api

all: proto

proto:
	protoc \
		--go_out=paths=source_relative,M$(PROTO_FILE)=$(GO_PKG_PATH):. \
		--go-grpc_out=paths=source_relative,M$(PROTO_FILE)=$(GO_PKG_PATH):. \
		$(PROTO_FILE)

server:
	go run server/main.go server/service.go

client:
	go run client/main.go

build:
	go build -o github-search-server ./server

test:
	go test ./server -v

clean:
	rm -f api/githubsearch.pb.go api/githubsearch_grpc.pb.go
	rm -f github-search-server

.PHONY: all proto server client build test clean
