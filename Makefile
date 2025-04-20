PROTO_FILE=api/githubsearch.proto
GO_PKG_PATH=github-search-service/api

all: proto

proto:
	protoc \
		--go_out=paths=source_relative,M$(PROTO_FILE)=$(GO_PKG_PATH):. \
		--go-grpc_out=paths=source_relative,M$(PROTO_FILE)=$(GO_PKG_PATH):. \
		$(PROTO_FILE)

clean:
	rm -f api/githubsearch.pb.go api/githubsearch_grpc.pb.go

.PHONY: all proto clean
