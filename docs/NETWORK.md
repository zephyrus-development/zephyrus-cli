# network.go Documentation

## Package utils

This module handles HTTP network operations for fetching files from the vault repository on GitHub.

### Imports

- `fmt`: String formatting and printing
- `io`: I/O utilities
- `net/http`: HTTP client and utilities
- `time`: Time package for cache busting

### Functions

#### FetchRaw

```go
func FetchRaw(username string, path string) ([]byte, error)
```

Fetches raw file content from a GitHub repository using the GitHub raw content API.

**Parameters:**
- `username`: The GitHub username hosting the vault (e.g., "myuser")
- `path`: The path to the file in the repository (e.g., "storage_id_hex" or ".config/index")

**Return:**
- `[]byte`: The raw file content
- `error`: Returns error if the request fails or file is not found

**URL Format:**

Files are fetched using GitHub's raw content API with a timestamp cache buster:
```
https://raw.githubusercontent.com/{username}/.zephyrus/master/{path}?t={nanoseconds}
```

**Process:**

1. Constructs the GitHub raw URL with cache buster query parameter
2. Creates an HTTP client with 10-second timeout
3. Makes GET request to fetch the file
4. Returns file content or error

**Error Handling:**
- Returns "404" error if file is not found (StatusCode 404)
- Returns "bad status: {code}" error for other HTTP errors
- Returns network-related errors if the request fails

**Special Features:**
- **Cache Busting**: Uses `time.Now().UnixNano()` as query parameter to bypass browser/CDN caching
- **Timeout**: 10-second timeout prevents hanging on network issues
- **Direct**: Uses raw GitHub API for direct file access without HTML wrapper

**Example Usage:**

```go
// Fetch an encrypted file
data, err := FetchRaw("myusername", "a3f2e1c9d4b6f8e2")
if err != nil {
    log.Fatal(err)
}

// Fetch the vault index
indexData, err := FetchRaw("myusername", ".config/index")
if err != nil {
    if err.Error() == "404" {
        fmt.Println("Index not found - new vault")
    } else {
        log.Fatal(err)
    }
}
```

### Typical Usage in Vault Operations

1. **Downloading Files**: Fetch encrypted file by storage ID
2. **Fetching Master Key**: Fetch `.config/key` at setup
3. **Fetching Index**: Fetch `.config/index` to get vault structure

### Notes

- Files must exist in the `.zephyrus` GitHub repository
- The repository can be private (accessed via SSH for git operations, but FetchRaw requires public access or GitHub token)
- Cache busting ensures fresh data on each request
- 10-second timeout is suitable for typical network conditions
- Error message "404" is specifically checked by caller code
