# upload.go Documentation

## Package utils

This module handles uploading files to the vault, including encryption, index updates, and remote repository synchronization.

### Imports

- `fmt`: String formatting and printing
- `os`: Operating system file operations

### Functions

#### UploadFile

```go
func UploadFile(sourcePath string, vaultPath string, session *Session) error
```

Uploads a file from the local filesystem to the vault, encrypting it and updating the vault index.

**Parameters:**
- `sourcePath`: The path to the file on the local filesystem (e.g., "./document.pdf")
- `vaultPath`: The destination path in the vault (e.g., "documents/report.pdf")
- `session`: The active session containing authentication credentials, encryption password, and vault index

**Return:**
- `error`: Returns error if file reading, encryption, index update, or push operations fail

**Upload Process:**

1. **Construct Repository URL**: Builds the git repository URL from session username

2. **Read Source File**: Reads the entire file from the local filesystem
   - Returns error if file doesn't exist or cannot be read

3. **Determine Storage Name**:
   - **Existing Files**: Checks if file already exists in vault index
     - Uses existing storage ID (RealName)
     - Displays "Updating existing file" message
   - **New Files**: Generates new random hex-based storage ID
     - Adds entry to vault index with `AddFile()`
     - Displays "Uploading new file" message

4. **Encrypt File Data**:
   - Encrypts file using the session password
   - Produces encrypted bytes ready for storage

5. **Encrypt Updated Index**:
   - Serializes and encrypts the updated vault index
   - Creates `.config/index` entry

6. **Prepare Push Package**:
   - Creates map containing:
     - The encrypted file (storage ID as key)
     - The encrypted index (`.config/index` as key)

7. **Push to Repository**:
   - Calls `PushFiles()` to push both files to GitHub
   - Uses SSH authentication from session
   - Commit message: "Nexus: Updated Vault"

**File Handling:**

**New Files:**
```
Local: ./document.pdf
Vault: documents/report.pdf
Storage: a3f2e1c9d4b6f8e2 (hex name, content = encrypted)
```

**Updating Existing Files:**
```
Vault: documents/report.pdf
Storage: (existing hex name - reused)
Action: File content replaced, index unchanged
```

**Error Handling:**
- Returns error if source file cannot be read
- Returns error if encryption fails
- Returns error if git push fails
- Returns error if index serialization fails

**Example Usage:**

```go
// Upload a new file
err := UploadFile("./quarterly_report.pdf", "documents/Q1.pdf", session)
if err != nil {
    fmt.Println("Upload failed:", err)
}

// Update an existing file
err = UploadFile("./updated_report.pdf", "documents/Q1.pdf", session)
if err != nil {
    fmt.Println("Update failed:", err)
}

// Upload to nested folders (creates folders if needed)
err = UploadFile("./archive.zip", "backups/2024/archive.zip", session)
if err != nil {
    fmt.Println("Upload failed:", err)
}
```

**Output Examples:**

Uploading new file:
```
Uploading new file: documents/report.pdf as a3f2e1c9d4b6f8e2
```

Updating existing file:
```
Updating existing file: documents/report.pdf (a3f2e1c9d4b6f8e2)
```

### Implementation Details

The upload process:
- Checks local index cache (not remote) for existing file
- Creates intermediate folders automatically
- Uses stateless push operation (doesn't redownload entire repository)
- Updates both file content and vault index in single push

### Security Considerations

- Files are encrypted with AES-256-GCM before upload
- Storage IDs are randomly generated (original filenames not exposed)
- Index is encrypted and stored in `.config/index`
- Encryption key is derived from session password

### File Size Considerations

- No explicit file size limits in the code
- Practical limits depend on GitHub's API and SSH transfer limits
- Large files should be zipped before upload

### Notes

- If an upload fails mid-process, the local session index may be updated but the push might fail
- The vault index is updated before pushing (optimistic approach)
- Consider saving the session after successful upload
- Intermediate folders are created automatically (no need to pre-create structure)
- Files with the same vault path but different source paths will overwrite

### Typical Upload Workflow

```go
// After connecting to vault
err := UploadFile("./my_document.pdf", "documents/my_document.pdf", session)
if err != nil {
    fmt.Println("Upload failed:", err)
    return
}

// File is now encrypted and stored in vault
// Index is updated and pushed to GitHub
// Both files and index are synced with remote
```

---

## Directory Upload

### UploadDirectory

```go
func UploadDirectory(sourceDirPath string, vaultPath string, session *Session) error
```

Uploads an entire directory recursively to the vault, preserving folder structure and encrypting all files.

**Parameters:**
- `sourceDirPath`: The path to the local directory (e.g., "./documents")
- `vaultPath`: The destination path in the vault (e.g., "backups/documents")
- `session`: The active session with authentication and encryption credentials

**Return:**
- `error`: Returns error if directory reading, any file encryption, or push fails

### Directory Upload Process

1. **Verify Directory**: Confirms path is a directory, not a file
2. **Create Vault Path**: Creates necessary folder structure in vault index
3. **Recursively Walk Files**: Iterates through all files using `filepath.Walk()`
4. **Process Each File**:
   - Maintains relative paths from source directory
   - Reads file content
   - Generates new encryption key for each file
   - Encrypts with vault password
   - Adds entry to vault index
   - Collects for batch upload
5. **Batch Upload**: Uploads all files to GitHub in single push
6. **Index Update**: Encrypts and uploads updated vault index

### Usage Examples

```bash
# Upload entire local directory
zep upload ./my-folder vault/my-folder

# Upload project directory
zep upload ./project backup/project

# Upload with custom name
zep upload ./documents archive/docs-backup
```

### Directory Structure Preservation

```
Local Directory:
./documents/
  report.pdf
  notes.txt
  subfolder/
    memo.pdf
    archive.zip

Vault Structure:
vault/documents/
  report.pdf
  notes.txt
  subfolder/
    memo.pdf
    archive.zip
```

### Progress Output

```
Scanning directory: ./documents
Processing file (1): report.pdf
  → New file: vault/documents/report.pdf as a3f2e1c9d4b6f8e2
Processing file (2): notes.txt
  → New file: vault/documents/notes.txt as b4g3f2d0e5c7h9f3
Processing file (3): subfolder/memo.pdf
  → New file: vault/documents/subfolder/memo.pdf as c5h4g3e1f6d8i0j4

Uploading 3 files to vault...
[1/2] Encrypting vault index...
[2/2] Uploading to GitHub...

✔ Successfully uploaded 3 files from directory
```

### Features

- ✅ **Recursive Processing**: Handles nested folders of any depth
- ✅ **Per-File Encryption**: Each file gets unique encryption key
- ✅ **Atomic Upload**: All files and index uploaded in single batch
- ✅ **Structure Preservation**: Folder hierarchy maintained exactly
- ✅ **Progress Tracking**: Shows file count and processing status
- ✅ **Error Handling**: Reports failures with detailed context

### Use Cases

1. **Backup Entire Projects**
   ```bash
   zep upload ./project backup/my-project
   ```

2. **Archive Bulk Documents**
   ```bash
   zep upload ./documents vault/documents-archive
   ```

3. **Upload Photo Collections**
   ```bash
   zep upload ./photos vault/photos/2024
   ```

### Related Commands

- [`zep download`](DOWNLOAD.md): Download entire directories
- [`zep list`](LIST.md): View vault structure
- [`zep transfer-vault`](TRANSFER.md): Migrate directories to another vault
