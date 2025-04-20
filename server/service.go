package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	pb "github-search-service/api"
)

type githubSearchServer struct {
	pb.UnimplementedGithubSearchServiceServer
}

type GitHubItem struct {
	HTMLURL string `json:"html_url"`
	Repo    struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type GitHubResponse struct {
	Items []GitHubItem `json:"items"`
}

func (s *githubSearchServer) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {

	// Construct GitHub search query with optional user filter
	query := req.SearchTerm
	if req.User != "" {
		query += " user:" + req.User
	}

	// Build GitHub API URL using the constructed query
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s", strings.ReplaceAll(query, " ", "+"))

	// Prepare authenticated HTTP request to GitHub API
	reqHTTP, _ := http.NewRequest("GET", url, nil)
	reqHTTP.Header.Set("Accept", "application/vnd.github+json")
	reqHTTP.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Set GitHub token if available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		reqHTTP.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		log.Println("GitHub API error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var ghResp GitHubResponse

	// Parse GitHub API response into Go structs
	json.Unmarshal(body, &ghResp)

	var results []*pb.Result

	// Convert GitHub search results into gRPC-compatible format
	for _, item := range ghResp.Items {
		results = append(results, &pb.Result{
			FileUrl: item.HTMLURL,
			Repo:    item.Repo.FullName,
		})
	}

	return &pb.SearchResponse{Results: results}, nil
}
