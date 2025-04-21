# GitHub Search Service

This is a gRPC-based microservice in Go that wraps GitHub's REST API to perform code searches across public repositories. It accepts a search term and an optional user filter, and returns a list of file URLs along with their respective repositories.

---

## Features
- gRPC server using Protobuf definitions
- Integration with GitHub's Search API (`/search/code`)
- Optional filtering by GitHub username
- Command-line flags for dynamic search term/user
- Docker support for containerized execution
- Unit test coverage for success and error flows
- Buf integration for clean proto generation and linting

---

## Code Flow Overview

1. **Client (`client/main.go`)**
    - Accepts `-term` and `-user` via CLI flags
    - Constructs a `SearchRequest` and sends it to the gRPC server
    - Logs results or errors from the server

2. **Server (`server/main.go` & `server/service.go`)**
    - Initializes a single shared `http.Client`
    - Listens on port 50051 and handles `Search` requests
    - Constructs GitHub API query
    - Handles API responses, including non-200 and parsing errors
    - Returns parsed results in `SearchResponse`

3. **API Protobuf (`api/githubsearch.proto`)**
    - Defines service `GithubSearchService` and messages `SearchRequest`, `SearchResponse`, `Result`
    - Supports generation via Buf

---

## Limitations

- GitHub API returns **only the first 30 results** by default.
- No pagination support implemented yet.

### Current Request Behavior

#### gRPC Request Parameters in Client Code
```go
resp, err := client.Search(ctx, &pb.SearchRequest{
    SearchTerm: "filename:Dockerfile",
    User:       "",
})
```

#### GitHub API Query Constructed
```http
GET https://api.github.com/search/code?q=filename:Dockerfile
```
- Uses `filename:` qualifier to narrow the search
- No `page` or `per_page`, so returns only first 30 results

### How to Fix Pagination

#### Option 1: Extend the Protobuf Definition (Recommended)
```proto
message SearchRequest {
  string search_term = 1;
  string user = 2;
  int32 page = 3;
  int32 per_page = 4;
}
```
And build URLs like:
```go
fmt.Sprintf("https://api.github.com/search/code?q=%s&page=%d&per_page=%d", query, page, perPage)
```

#### Option 2: Paginate Internally (No Proto Change)
Loop inside the server:
```go
for page := 1; page <= 3; page++ {
  // build ?page=X&per_page=30, merge results
}
```
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
- **Buf CLI (for proto linting and generation)**
  - MacOS: `brew install bufbuild/buf/buf`
  - Linux:
    ```bash
    curl -sSL https://github.com/bufbuild/buf/releases/latest/download/buf-$(uname -s)-$(uname -m) \
      -o /usr/local/bin/buf && chmod +x /usr/local/bin/buf
    ```
- GitHub Personal Access Token (export as `GITHUB_TOKEN`)

---

## Environment Variables
- `GITHUB_TOKEN`: GitHub personal access token for authenticated requests
- `SERVER_ADDR`: Optional override for gRPC server address (defaults to `localhost:50051`)

---

## Running Locally (Without Docker)

### Step 1: Generate gRPC code
```bash
make proto-buf
```

### Step 2: Export GitHub token
```bash
export GITHUB_TOKEN=ghp_your_token_here
```

### Step 3: Start server
```bash
make server
```

### Step 4: Run client with flags
```bash
make client term="filename:Dockerfile" user="torvalds"
```

---

## Running With Docker

### Step 1: Build Docker image
```bash
docker build -t github-search-service .
```

### Step 2: Run with GitHub token
```bash
docker run -e GITHUB_TOKEN=ghp_your_token_here -p 50051:50051 github-search-service
```

---

## Running With Docker Compose

### Step 1: Optional `.env` file
```env
GITHUB_TOKEN=ghp_your_token_here
```

### Step 2: Launch stack
```bash
docker-compose up --build
```

---

## Makefile Targets
```makefile
make proto         # (legacy) Generate code using protoc
make proto-buf     # Generate code using Buf
make lint-buf      # Lint proto files using Buf
make build         # Build server binary
make server        # Run server locally
make client        # Run client with flags: make client term=... user=...
make test          # Run unit tests
make clean         # Cleanup generated files and binaries
```

---

## Repository Structure
```
.
├── api/                    # Proto definitions and generated code
├── client/                 # gRPC client code
├── server/                 # gRPC server implementation
├── Dockerfile              # Dockerfile for server
├── Dockerfile.client       # Dockerfile for client
├── docker-compose.yml      # Docker Compose setup
├── Makefile                # Dev and build automation
├── buf.gen.yaml            # Buf code generation config
├── buf.yaml                # Buf lint/breaking config
├── go.mod / go.sum         # Dependencies
└── README.md
```
