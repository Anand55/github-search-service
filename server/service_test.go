package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	pb "github-search-service/api"
)

type mockRoundTripper struct {
	ResponseBody string
	StatusCode   int
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: m.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.ResponseBody)),
		Header:     make(http.Header),
	}, nil
}

// --- Tests ---

func TestSearch_Success(t *testing.T) {
	mockData := GitHubResponse{
		Items: []GitHubItem{
			{
				HTMLURL: "https://github.com/mock/repo/blob/main/main.go",
				Repo: struct {
					FullName string `json:"full_name"`
				}{FullName: "mock/repo"},
			},
		},
	}
	mockJSON, _ := json.Marshal(mockData)

	client := &http.Client{
		Transport: &mockRoundTripper{
			ResponseBody: string(mockJSON),
			StatusCode:   200,
		},
	}

	server := &githubSearchServer{
		httpClient: client,
	}

	resp, err := server.Search(context.Background(), &pb.SearchRequest{
		SearchTerm: "filename:Dockerfile",
		User:       "",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}

	if resp.Results[0].Repo != "mock/repo" {
		t.Errorf("unexpected repo name: got %s", resp.Results[0].Repo)
	}
}

func TestSearch_Non200Response(t *testing.T) {
	client := &http.Client{
		Transport: &mockRoundTripper{
			ResponseBody: `{"message": "Bad credentials"}`,
			StatusCode:   403,
		},
	}

	server := &githubSearchServer{
		httpClient: client,
	}

	_, err := server.Search(context.Background(), &pb.SearchRequest{
		SearchTerm: "filename:Dockerfile",
		User:       "",
	})

	if err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}
}
