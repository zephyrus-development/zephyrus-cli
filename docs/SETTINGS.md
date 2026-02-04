# Settings Module

The Settings module provides persistent configuration for Zephyrus vaults, allowing users to customize git commit attribution, commit messages, and encryption hash lengths.

## Overview

Settings are stored encrypted in `.config/settings` in your vault repository and are loaded during authentication. This allows configurations to be shared across machines while maintaining security.

## Data Structure

```go
type VaultSettings struct {
    CommitAuthorName  string // Author name for git commits (default: "Zephyrus")
    CommitAuthorEmail string // Author email for git commits (default: "auchrio@proton.me")
    CommitMessage     string // Default commit message (default: "Zephyrus: Updated Vault")
    FileHashLength    int    // Length of file storage IDs in hex (default: 16, range: 8-64)
    ShareHashLength   int    // Length of share references in base62 (default: 6, range: 4-32)
}
```

## Commands

### `settings show` - Display Current Settings

Display all current vault settings.

**Usage:**
```bash
./zep settings show
```

**Example Output:**
```
=== Vault Settings ===
Commit Author: Zephyrus <auchrio@proton.me>
Commit Message: Zephyrus: Updated Vault
File Hash Length: 16 chars
Share Hash Length: 6 chars
```

### `settings set` - Modify Settings

Update a specific setting value.

**Usage:**
```bash
./zep settings set [key] [value]
```

**Arguments:**
- `key`: Setting name (case-insensitive)
  - `author_name` - Git commit author name
  - `author_email` - Git commit author email
  - `commit_message` - Default git commit message
  - `file_hash_length` - File storage ID length
  - `share_hash_length` - Share reference hash length
- `value`: New value for the setting

**Examples:**
```bash
# Set custom git author
./zep settings set author_name "John Doe"
./zep settings set author_email "john@example.com"

# Customize commit message
./zep settings set commit_message "Backup: $(date)"

# Adjust hash lengths (longer = more unique IDs, shorter = more compact)
./zep settings set file_hash_length 32
./zep settings set share_hash_length 8
```

## Default Values

| Setting | Default | Range | Purpose |
|---------|---------|-------|---------|
| CommitAuthorName | "Zephyrus" | Any string | Name for git commits |
| CommitAuthorEmail | "auchrio@proton.me" | Any email | Email for git commits |
| CommitMessage | "Zephyrus: Updated Vault" | Any string | Message for all operations |
| FileHashLength | 16 | 8-64 | Characters in storage ID |
| ShareHashLength | 6 | 4-32 | Characters in share ref |

## Technical Details

### Encryption and Storage

Settings are encrypted using the same AES-256-GCM encryption as the vault index:
1. Serialize `VaultSettings` to JSON
2. Encrypt with user's vault password using PBKDF2 key derivation
3. Store in `.config/settings` on GitHub

### Loading Settings

Settings are fetched during authentication and cached in the `Session` struct for the duration of the session. Each command uses the cached settings values.

### Persistence

Settings changes are immediately pushed to `.config/settings` on GitHub and also update the local `zephyrus.conf` if in persistent mode.

## Use Cases

### Custom Git Attribution

Track operations by user or team:
```bash
./zep settings set author_name "Engineering Team"
./zep settings set author_email "eng@company.com"
```

All future uploads, downloads, deletes, and shares will use these credentials in git history.

### Backup Identification

Embed timestamps or context in commit messages:
```bash
./zep settings set commit_message "Daily Backup $(date +%Y-%m-%d)"
```

### Security vs. Compactness

Adjust hash lengths based on your needs:
- **Shorter hashes**: Faster lookups, more compact storage
- **Longer hashes**: Lower collision risk, better security

```bash
# Very secure, longer file IDs
./zep settings set file_hash_length 64

# Fast and compact
./zep settings set file_hash_length 8
```

## Validation

Settings are validated when set:
- Hash lengths must be within defined ranges
- Email format is validated
- Empty strings are rejected (except for messages)

## See Also

- [Auth Module](AUTH.md) - Session management and authentication
- [Encryption Module](ENCRYPTION.md) - Cryptographic operations
- [Git Module](GIT.md) - Git operations with custom authors
