package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "github-search-service/api"

	"google.golang.org/grpc"
)

func main() {
	searchTerm := flag.String("term", "", "Search term (e.g., filename:Dockerfile)")
	searchUser := flag.String("user", "", "GitHub username to scope the search")
	flag.Parse()

	if *searchTerm == "" {
		log.Fatalln("Search term is required. Use -term flag to specify it.")
	}

	serverAddr := getServerAddr()
	conn, client := newSearchClient(serverAddr)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Sending search request...")

	resp, err := performSearch(ctx, client, *searchTerm, *searchUser)
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

func getServerAddr() string {
	addr := os.Getenv("SERVER_ADDR")
	if addr == "" {
		return "localhost:50051"
	}
	return addr
}

func newSearchClient(addr string) (*grpc.ClientConn, pb.GithubSearchServiceClient) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	client := pb.NewGithubSearchServiceClient(conn)
	return conn, client
}

func performSearch(ctx context.Context, client pb.GithubSearchServiceClient, term, user string) (*pb.SearchResponse, error) {
	return client.Search(ctx, &pb.SearchRequest{
		SearchTerm: term,
		User:       user,
	})
}
