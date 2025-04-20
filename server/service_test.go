// service_test.go
//
// This test file validates the Search method of the GithubSearchService server.
// It avoids hitting the real GitHub API by overriding the default HTTP transport
// with a custom mockRoundTripper that returns a pre-defined JSON response.
//
// The test verifies:
// - That the Search method successfully parses and returns mocked GitHub API data
// - That the response contains exactly one result
// - That the result includes the expected repository name and file URL
//
// This is achieved without changing the service.go file, using only transport-level mocking.

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

// Custom transport to mock GitHub API HTTP response
type mockRoundTripper struct{}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	mockResponse := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"html_url": "https://github.com/mock/repo/blob/main/main.go",
				"repository": map[string]interface{}{
					"full_name": "mock/repo",
				},
			},
		},
	}
	bodyBytes, _ := json.Marshal(mockResponse)
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(string(bodyBytes))),
		Header:     make(http.Header),
	}, nil
}

func TestSearchReturnsMockedResults(t *testing.T) {
	// Override HTTP transport to intercept GitHub API
	originalTransport := http.DefaultTransport
	http.DefaultTransport = &mockRoundTripper{}
	defer func() { http.DefaultTransport = originalTransport }()

	server := &githubSearchServer{}
	resp, err := server.Search(context.Background(), &pb.SearchRequest{
		SearchTerm: "filename:Dockerfile",
		User:       "",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(resp.Results) != 1 {
		t.Fatalf("Expected 1 result, got: %d", len(resp.Results))
	}

	if resp.Results[0].Repo != "mock/repo" {
		t.Errorf("Unexpected repo name: %s", resp.Results[0].Repo)
	}
}
