# GitHub Search Service

This is a gRPC-based microservice in Go that wraps GitHub's REST API to perform code searches across public repositories. It accepts a search term and an optional user filter, and returns a list of file URLs along with their respective repositories.

---

## Features
- gRPC server using Protobuf definitions
- Integration with GitHub's Search API (`/search/code`)
- Optional filtering by GitHub username
- Docker support for containerized execution

---

## Code Flow Overview

1. **Client (`client/main.go`)**
    - Constructs a `SearchRequest` with a `search_term` and an optional `user`.
    - Sends the request to the gRPC server at `localhost:50051` (or Docker Compose service name).

2. **Server (`server/main.go` & `server/service.go`)**
    - Listens on port 50051.
    - On receiving a request, constructs a GitHub API query URL.
    - Adds authentication and headers.
    - Parses the JSON response.
    - Converts results to gRPC-compatible `SearchResponse`.

3. **API Protobuf (`api/githubsearch.proto`)**
    - Defines service `GithubSearchService` and messages `SearchRequest`, `SearchResponse`, `Result`.

---

## Limitations

- By default, the GitHub API only returns **30 results** per request (first page only).
- The current implementation does **not support pagination** and therefore only fetches the first page.

### Current Request Behavior

#### gRPC Request Parameters in Client Code
```go
resp, err := client.Search(ctx, &pb.SearchRequest{
    SearchTerm: "filename:Dockerfile", // exact filename match
    User:       "",                    // optional GitHub username (empty = global search)
})
```

#### GitHub API Query Constructed
```http
GET https://api.github.com/search/code?q=filename:Dockerfile
```
- The `q` parameter combines both the search term and (optionally) the user qualifier.
- No `page` or `per_page` parameters are passed, so GitHub returns **only the first 30 results** by default.
- Constructed GitHub API call: `GET /search/code?q=<search_term> [user:<username>]`
- Returns at most 30 results due to GitHub API's default `per_page` setting.

### How to Fix Pagination

#### Option 1: Extend the Protobuf Definition (Recommended)
Update `api/githubsearch.proto` to:
```proto
message SearchRequest {
  string search_term = 1;
  string user = 2;
  int32 page = 3;         // page number
  int32 per_page = 4;     // results per page (max 100)
}
```
Then update server to build GitHub query like:
```go
fmt.Sprintf("https://api.github.com/search/code?q=%s&page=%d&per_page=%d", query, page, perPage)
```

#### Option 2: Handle Pagination Internally (No Protobuf Change)
Loop over `page=1..N` inside the server and merge all results:
```go
for page := 1; page <= 3; page++ {
  // build URL with ?page=page&per_page=30
  // merge results
}
```
This hides pagination from the client, but removes fine control over page/limit.

---

## Prerequisites

### Local Requirements (Non-Docker)
- Go (>= 1.18)
- `protoc` (Protocol Buffers Compiler)
- Go plugins for protoc:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
- GitHub Personal Access Token (set as environment variable `GITHUB_TOKEN`)

---

## Environment Variables
- `GITHUB_TOKEN`: GitHub personal access token for authenticated requests to increase rate limits.

---

## Running Locally (Without Docker)

### Step 1: Generate gRPC code
```bash
make
```

### Step 2: Export GitHub token
```bash
export GITHUB_TOKEN=ghp_your_token_here
```

### Step 3: Start server
```bash
go run server/main.go server/service.go
```

### Step 4: In another terminal, run client
```bash
go run client/main.go
```

---

## Running With Docker

### Step 1: Build Docker image
```bash
docker build -t github-search-service .
```

### Step 2: Run container with GitHub token
```bash
docker run -e GITHUB_TOKEN=ghp_your_token_here -p 50051:50051 github-search-service
```

---

## Running With Docker Compose

### Step 1: Create `.env` file (optional)
```env
GITHUB_TOKEN=ghp_your_token_here
```

### Step 2: Start all services
```bash
docker-compose up --build
```

- The server will be reachable as `github-search-server:50051` from inside Docker.
- The client will log results returned from the gRPC server.

---

## Repository Structure
```
.
├── api/                    # Proto definitions and generated code
├── client/                 # gRPC client code
├── server/                 # gRPC server implementation
├── Dockerfile              # Dockerfile for building server image
├── Dockerfile.client       # Dockerfile for client (used in Docker Compose)
├── docker-compose.yml      # Compose file to run server + client
├── Makefile                # Builds proto files
├── go.mod / go.sum         # Go dependencies
└── README.md
```
