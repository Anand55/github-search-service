package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github-search-service/api"

	"google.golang.org/grpc"
)

func main() {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = "localhost:50051"
	}
	// Set up a gRPC connection to the server (Docker Compose: use service name instead of localhost)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	// Create gRPC client from generated code
	client := pb.NewGithubSearchServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	log.Println("Sending search request...")

	// Build the request with search term and optional GitHub user
	resp, err := client.Search(ctx, &pb.SearchRequest{
		SearchTerm: "filename:Dockerfile",
		User:       "",
	})

	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	if len(resp.Results) == 0 {
		log.Println("No results found.")
	} else {
		for _, result := range resp.Results {
			log.Printf("Repo: %s\nFile: %s\n", result.Repo, result.FileUrl)
		}
	}
}
