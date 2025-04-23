package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	pb "github-search-service/api"
)

type githubSearchServer struct {
	pb.UnimplementedGithubSearchServiceServer
	httpClient *http.Client
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
	query := req.SearchTerm
	if req.User != "" {
		query += " user:" + req.User
	}

	url := fmt.Sprintf("https://api.github.com/search/code?q=%s", strings.ReplaceAll(query, " ", "+"))

	reqHTTP, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to build HTTP request: %v", err)
		return nil, fmt.Errorf("failed to build GitHub API request: %w", err)
	}

	reqHTTP.Header.Set("Accept", "application/vnd.github+json")
	reqHTTP.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		reqHTTP.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.httpClient.Do(reqHTTP)
	if err != nil {
		log.Printf("GitHub API call failed: %v", err)
		return nil, fmt.Errorf("GitHub API call failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read GitHub response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("GitHub API returned non-200 status: %d\nResponse: %s", resp.StatusCode, string(body))
		return nil, fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	var ghResp GitHubResponse
	if err := json.Unmarshal(body, &ghResp); err != nil {
		log.Printf("Failed to parse GitHub response: %v", err)
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	var results []*pb.Result
	for _, item := range ghResp.Items {
		results = append(results, &pb.Result{
			FileUrl: item.HTMLURL,
			Repo:    item.Repo.FullName,
		})
	}

	log.Printf("Successfully fetched %d results", len(results))
	return &pb.SearchResponse{Results: results}, nil
}
