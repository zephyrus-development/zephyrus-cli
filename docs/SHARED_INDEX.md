# Shared File Index Module

The Shared File Index module manages the encrypted index of shared files, tracking which files have been shared and their associated encryption metadata.

## Overview

The Shared Index stores information about shared files in an encrypted file structure (`shared_index.json`) in the `.zephyrus` repository. This allows the system to:
- Track which files have been shared
- Store share reference IDs and share passwords
- Map share IDs back to original files
- Revoke access by managing share entries

## Data Structure

### Shared Index Format

```json
{
  "shares": {
    "72cTWg": {
      "filename": "report.pdf",
      "vault_path": "documents/report.pdf",
      "share_password": "encrypted_password",
      "created_at": "2026-02-04T15:30:00Z",
      "access_count": 0
    },
    "AbXkLm": {
      "filename": "budget.xlsx",
      "vault_path": "financial/budget.xlsx",
      "share_password": "encrypted_password",
      "created_at": "2026-02-03T10:15:00Z",
      "access_count": 0
    }
  }
}
```

### Share Entry Fields

| Field | Type | Purpose |
|-------|------|---------|
| `filename` | string | Original filename for reference |
| `vault_path` | string | Full path in vault |
| `share_password` | string | Encrypted share password |
| `created_at` | ISO 8601 | Timestamp of share creation |
| `access_count` | integer | Number of times shared file was accessed |

## Functions

### `LoadSharedIndex`

Load shared index from encrypted storage.

**Function Signature:**
```go
func LoadSharedIndex(username string, password string) (*SharedIndex, error)
```

**Parameters:**
- `username`: GitHub username
- `password`: Vault password

**Returns:**
- Populated `SharedIndex` struct
- Error if file not found or decryption fails

**Process:**
1. Fetch encrypted `shared_index.json` from GitHub
2. Decrypt using vault password
3. Parse JSON into SharedIndex struct
4. Return index

### `SaveSharedIndex`

Save shared index to encrypted storage.

**Function Signature:**
```go
func (si *SharedIndex) Save(username string, password string, keyPath string) error
```

**Parameters:**
- `username`: GitHub username
- `password`: Vault password
- `keyPath`: Path to SSH private key

**Returns:**
- Error if save fails

**Process:**
1. Serialize SharedIndex to JSON
2. Encrypt using vault password with PBKDF2
3. Push to `.zephyrus/shared_index.json` via git
4. Return error or nil

### `AddShare`

Add a new share entry to the index.

**Function Signature:**
```go
func (si *SharedIndex) AddShare(shareID string, filename string, vaultPath string, sharePassword string) error
```

**Parameters:**
- `shareID`: Unique reference ID for share
- `filename`: Original filename
- `vaultPath`: Full vault path
- `sharePassword`: Password for share decryption

**Returns:**
- Error if entry already exists

### `RemoveShare`

Remove a share entry from the index.

**Function Signature:**
```go
func (si *SharedIndex) RemoveShare(shareID string) error
```

**Parameters:**
- `shareID`: Share reference to revoke

**Returns:**
- Error if share not found

### `GetShare`

Retrieve a specific share entry.

**Function Signature:**
```go
func (si *SharedIndex) GetShare(shareID string) (*ShareEntry, error)
```

**Returns:**
- ShareEntry struct
- Error if not found

## Encryption

### Share Password Storage

Share passwords are encrypted before storage:
1. Generate random share password (user-provided or generated)
2. Encrypt with vault password using AES-256-GCM
3. Store in index as encrypted hex string
4. Decrypt when needed for share operations

### Index Encryption

The entire shared index is encrypted:
- Serialized to JSON
- Encrypted with vault password using PBKDF2 + AES-256-GCM
- Stored as `.zephyrus/shared_index.json`
- Same encryption as vault index

## Data Flow

### Creating a Share

```
User calls: share documents/report.pdf
  ↓
Generate share ID (base62, 6 chars)
  ↓
Encrypt file with share password
  ↓
Store share reference in index
  ↓
Save encrypted index to GitHub
  ↓
Return share string to user
```

### Revoking a Share

```
User calls: shared rm 72cTWg
  ↓
Load shared index
  ↓
Find and remove share entry
  ↓
Save encrypted index to GitHub
  ↓
Share link becomes invalid
```

## File Storage

The shared index is stored at:
```
.zephyrus/shared_index.json (encrypted)
```

Format:
```
[16-byte salt][12-byte nonce][encrypted JSON + auth tag]
```

## Integration

The Shared Index is used by:
- **Share Module**: Creates entries when sharing files
- **Shared Manage**: Removes entries when revoking
- **Shared Search**: Queries index for shares
- **Authentication**: Loaded during session setup

## Backup and Recovery

### Backup Shared Index

```bash
# Download shared index manually if needed
./zep download .zephyrus/shared_index.json ./shared_index_backup.json
```

### Recovery

If shared index is corrupted:
1. Delete `.zephyrus/shared_index.json` from GitHub
2. Existing shares remain functional (still encrypted on GitHub)
3. New shares will start with empty index
4. Can manually recreate shares using `share` command

## Limitations

### No Version History

- Shared index overwrites with each change
- No history of share operations
- But: Git commits show what changed

### Synchronization

- Single source of truth on GitHub
- Local cache exists in session
- Multiple concurrent edits may conflict
- Use `connect`/`disconnect` for explicit sync

## Security Considerations

### Share Password Protection

- Share passwords are encrypted in the index
- Only vault password holder can see them
- Recipients use different share password

### Index Compromise

If shared index is leaked:
- Individual share passwords are encrypted
- Attacker cannot decrypt without vault password
- Shared files themselves are still encrypted

### Access Tracking

- `access_count` field can track usage (not currently incremented)
- Could be enhanced for audit logging
- No automatic expiration of shares

## See Also

- [Share Module](SHARE.md) - Creating shares
- [Shared Manage Module](SHARED_MANAGE.md) - Revoking shares
- [Shared Search Module](SHARED_SEARCH.md) - Finding shares
- [Encryption Module](ENCRYPTION.md) - Cryptographic details
- [Index Module](INDEX.md) - Vault index structure
