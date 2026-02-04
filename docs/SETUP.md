# setup.go Documentation

## Package utils

This module handles the initial setup and configuration of a Nexus vault, including repository verification and encryption key initialization.

### Imports

- `bufio`: Buffered I/O
- `fmt`: String formatting and printing
- `net/http`: HTTP client for repository verification
- `os`: Operating system file operations
- `strings`: String manipulation utilities
- `time`: Time package for commit timestamps
- `github.com/go-git/go-billy/v5/memfs`: In-memory filesystem
- `github.com/go-git/go-git/v5`: Git operations library
- `github.com/go-git/go-git/v5/config`: Git configuration
- `github.com/go-git/go-git/v5/plumbing/object`: Git object types
- `github.com/go-git/go-git/v5/plumbing/transport/ssh`: SSH transport for git
- `github.com/go-git/go-git/v5/storage/memory`: In-memory git storage
- `golang.org/x/crypto/ssh`: SSH cryptography utilities

### Functions

#### SetupVault

```go
func SetupVault(githubUser string, keyFilePath string, password string) error
```

Initializes a new Nexus vault by verifying the GitHub repository, encrypting the private key, and pushing the encrypted key to the repository.

**Parameters:**
- `githubUser`: GitHub username hosting the vault (e.g., "myusername")
  - If empty, prompts user for input
- `keyFilePath`: Path to the GitHub SSH private key file
  - If empty, prompts user for input
  - Example: "~/.ssh/id_ed25519" or "C:\Users\User\.ssh\id_ed25519"
- `password`: The vault password for encrypting the key
  - If empty, prompts user for input (always prompted if not provided)

**Return:**
- `error`: Returns error if repository verification, key reading, encryption, or push fails

**Setup Process:**

1. **Resolve Username**:
   - Uses provided username or prompts user if empty
   - Trims whitespace from input

2. **Verify Repository**:
   - Constructs GitHub repository URLs
   - Makes HTTP HEAD request to verify repository exists
   - Returns error if repository not found or inaccessible

3. **Resolve Key Path**:
   - Uses provided path or prompts user if empty
   - Trims whitespace from input

4. **Resolve Password**:
   - Prompts user for password (always prompted if not provided)
   - Uses standard input (echo enabled during setup)

5. **Encrypt Private Key**:
   - Reads the SSH private key file
   - Encrypts it using the provided password
   - Stores encrypted key in `.config/key`

6. **Initialize Git Repository**:
   - Creates empty git repository in memory
   - Creates `.config` directory
   - Writes encrypted key to `.config/key`
   - Stages the file for commit

7. **Create Initial Commit**:
   - Commits with message "Nexus: Setup Complete"
   - Author: "Nexus" <setup@cli.io>

8. **Push to GitHub**:
   - Sets up SSH authentication using the private key
   - Creates remote configuration
   - Force-pushes commit to master branch

**Error Handling:**
- Returns error if repository not found at expected URL
- Returns error if private key file cannot be read
- Returns error if encryption fails
- Returns error if git operations fail

**Pre-requisites:**

Before running setup:
1. Create a GitHub repository named `.zephyrus` in your account
2. Have SSH access to GitHub configured (or provide path to private key)
3. Know the password you want to use for the vault

**Example Usage:**

```go
// Interactive setup (prompts for all inputs)
err := SetupVault("", "", "")
if err != nil {
    log.Fatal(err)
}

// Programmatic setup with parameters
err = SetupVault(
    "myusername",
    "/home/user/.ssh/id_ed25519",
    "myVaultPassword123",
)
if err != nil {
    log.Fatal(err)
}
```

**Example Interactive Session:**

```
Enter GitHub Username: myusername
Enter Path to GitHub Private Key (e.g., ~/.ssh/id_ed25519): ~/.ssh/id_ed25519
Create a Vault Password (to encrypt your cloud key): ••••••••••
```

### What Gets Stored

After setup, your GitHub `.zephyrus` repository will contain:
```
.config/
  key  (your encrypted SSH private key)
```

### Security Considerations

- The private key is encrypted with your vault password using AES-256-GCM
- The password is never stored; it must be provided each time you connect
- The encrypted key is stored on GitHub (can be accessed if repository is public)
- SSH key is protected by vault password in transit and at rest

### Typical Workflow

```go
// 1. User creates .zephyrus repository on GitHub
// 2. User has SSH private key for GitHub
// 3. Run setup
err := SetupVault("", "", "")
if err != nil {
    fmt.Println("Setup failed:", err)
    return
}

// 4. Now user can connect
err = Connect("myusername", "myVaultPassword123")
if err != nil {
    fmt.Println("Connection failed:", err)
}
```

### Notes

- The repository must be created manually on GitHub before running setup
- Repository name must be exactly `.zephyrus` (with the dot)
- SSH key path can be relative (e.g., `~/.ssh/id_ed25519`) - paths will be expanded
- Password setup does not use `GetPassword()` (echo enabled); consider updating for security
- Force push is used to ensure setup succeeds even if repository isn't empty
