
    FROM golang:1.23 as builder

    WORKDIR /app
    
    # Install protoc (protobuf compiler)
    RUN apt-get update && apt-get install -y protobuf-compiler
    
    # Copy Go module files and download dependencies
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the entire project (including Makefile and source files)
    COPY . .
    
    # Install Go protoc plugins
    RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    
    # Set PATH for protoc-gen binaries
    ENV PATH="/go/bin:$PATH"
    
    # Generate gRPC code
    RUN make
    
    # Build the gRPC server binary
    RUN CGO_ENABLED=0 GOOS=linux go build -o github-search-server ./server
    
    FROM debian:bookworm-slim

    RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
    
    WORKDIR /app
    
    # Copy the compiled server binary from the builder
    COPY --from=builder /app/github-search-server ./
    
    EXPOSE 50051
    
    ENTRYPOINT ["./github-search-server"]
    