package main

import (
	"log"
	"net"
	"net/http"

	pb "github-search-service/api"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting gRPC server on port 50051...")

	client := &http.Client{}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGithubSearchServiceServer(grpcServer, &githubSearchServer{
		httpClient: client,
	})

	log.Println("gRPC server running on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
