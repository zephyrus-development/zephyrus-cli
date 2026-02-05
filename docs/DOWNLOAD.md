# download.go Documentation

## Package utils

This module handles downloading files from the Nexus vault, including retrieval, decryption, and local storage operations.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system file operations

### Functions

#### DownloadFile

```go
func DownloadFile(vaultPath string, outputPath string, session *Session) error
```

Downloads a file from the vault, decrypts it, and saves it to the local filesystem.

**Parameters:**
- `vaultPath`: The path to the file in the vault (e.g., "documents/report.pdf")
- `outputPath`: The local filesystem path where the file will be saved
- `session`: The active session containing authentication credentials and vault index

**Return:**
- `error`: Returns an error if any operation fails (file lookup, fetching, decryption, or writing)

**Process:**

1. **Find Entry**: Uses the vault index to locate the file entry at the specified vault path
2. **Validate Type**: Ensures the target is a file, not a folder (folders cannot be downloaded directly)
3. **Fetch Encrypted File**: Retrieves the encrypted file from the remote repository using its storage ID
4. **Decrypt**: Decrypts the file using the session password
5. **Save Locally**: Writes the decrypted file to the specified output path on the local filesystem

**Error Handling:**
- Returns an error if the file is not found in the vault
- Returns an error if attempting to download a folder
- Returns an error if the remote file fetch fails
- Returns an error if decryption fails (usually indicates invalid password)
- Returns an error if local file write fails

**Example Usage:**

```go
err := DownloadFile("documents/report.pdf", "./report.pdf", session)
if err != nil {
    fmt.Println("Download failed:", err)
}
```

**Notes:**

- Files are stored in the vault with hex-based storage IDs (RealName) for privacy
- The actual filename stored remotely does not reveal the original file name
- Downloaded files are created with permission mode 0644 (read/write for owner, read-only for others)
- Decryption failures typically indicate an incorrect password or corrupted file data

---

## Directory Download

### DownloadDirectory

```go
func DownloadDirectory(vaultPath string, outputDir string, session *Session) error
```

Downloads an entire directory from the vault recursively, preserving folder structure and decrypting all files.

**Parameters:**
- `vaultPath`: The path to folder in the vault (e.g., "documents" or "backups/archive")
- `outputDir`: The local filesystem directory where files will be saved
- `session`: The active session with authentication and decryption credentials

**Return:**
- `error`: Returns error if vault path doesn't exist, isn't a directory, or any file fetch/decrypt operation fails

### Directory Download Process

1. **Verify Vault Path**: Confirms path exists and is a directory in vault index
2. **Create Local Directory**: Creates the output directory structure with `os.MkdirAll()`
3. **Recursively Walk Vault**: Traverses all entries in the vault directory
4. **Process Each File**:
   - Maintains relative paths from source directory
   - Fetches encrypted file from remote storage
   - Decrypts with session password
   - Creates any intermediate directories as needed
   - Saves file to local path
5. **Preserve Structure**: All nested folders are recreated exactly as they appear in vault

### Usage Examples

```bash
# Download entire vault directory
zep download vault/documents ./documents

# Download backed-up project
zep download backup/my-project ./restored-project

# Download archived files to custom location
zep download archive/2023 ./archives/2023
```

### Directory Structure Preservation

```
Vault Structure:
vault/documents/
  report.pdf
  notes.txt
  subfolder/
    memo.pdf
    archive.zip

Downloaded to local:
./documents/
  report.pdf
  notes.txt
  subfolder/
    memo.pdf
    archive.zip
```

### Progress Output

```
Downloading directory: vault/documents
Creating directory structure...
Processing file (1): report.pdf
  → Fetching and decrypting...
  → Saved to ./documents/report.pdf
Processing file (2): notes.txt
  → Fetching and decrypting...
  → Saved to ./documents/notes.txt
Processing file (3): subfolder/memo.pdf
  → Fetching and decrypting...
  → Saved to ./documents/subfolder/memo.pdf

✔ Successfully downloaded 3 files from vault
```

### Features

- ✅ **Recursive Processing**: Handles nested folders of any depth
- ✅ **Atomic Decryption**: Each file decrypted with correct per-file key
- ✅ **Structure Preservation**: Folder hierarchy maintained exactly
- ✅ **Progress Tracking**: Shows file count and processing status
- ✅ **Auto-Folder Creation**: Creates intermediate directories as needed
- ✅ **Error Handling**: Continues on file errors and reports summary
- ✅ **Permission Preservation**: Files created with secure permissions (0644)

### Use Cases

1. **Restore Complete Backups**
   ```bash
   zep download backup/full-backup ./restored
   ```

2. **Download Project Archive**
   ```bash
   zep download archive/project-2024 ./projects/2024
   ```

3. **Retrieve Bulk Photo Collection**
   ```bash
   zep download photos/vacation ./vacation-photos
   ```

4. **Extract Nested Documents**
   ```bash
   zep download vault/documents/legal ./legal-docs
   ```

### Download Speed Considerations

- Download speed depends on:
  - Number of files in directory
  - Size of individual files
  - Network bandwidth
  - GitHub API rate limits
- Large directories with thousands of files may take several minutes
- Consider using `zep list` first to see directory contents and estimate size

### Error Handling

**Password Decryption Issues:**
- If download fails with "decryption failed": verify the session password is correct
- Files encrypted with a different password cannot be decrypted

**Directory Not Found:**
- Error if specified vault path doesn't exist
- Use `zep list` to verify path exists

**Disk Space Issues:**
- Ensure output directory location has sufficient disk space
- Large downloads may fail if insufficient space available

### Comparing Single File vs Directory Download

```bash
# Single file download (fast)
zep download vault/documents/report.pdf ./report.pdf

# Multiple directory download (uses batch operations for efficiency)
zep download vault/documents ./documents
# This downloads all files in documents/ recursively
```

### Related Commands

- [`zep upload`](UPLOAD.md): Upload directories to vault
- [`zep list`](LIST.md): View vault structure before downloading
- [`zep transfer-vault`](TRANSFER.md): Migrate directories to another vault
- [`zep delete`](DELETE.md): Remove directories from vault
