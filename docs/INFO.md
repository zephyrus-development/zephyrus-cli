# Info Module

The Info module provides detailed information about vault contents, file metadata, and vault statistics.

## Overview

The Info module displays comprehensive information about your vault and individual files, including size, creation/modification times, and metadata from the vault index.

## Commands

### `info` - Display Vault or File Information

Display vault statistics or detailed information about a specific file.

**Usage:**
```bash
./zep info [file-path]
```

**Arguments:**
- `file-path` (optional): Path to file or folder. If omitted, shows vault statistics.

## Vault Information

When called without arguments, displays overall vault statistics.

**Example:**
```bash
./zep info
```

**Output:**
```
=== Vault Statistics ===
Total Files: 42
Total Folders: 8
Vault Size: 2.3 GB
Storage ID Length: 16 chars
Share Hash Length: 6 chars
Committed Author: Zephyrus <auchrio@proton.me>
Commit Message: Zephyrus: Updated Vault
```

**Information Provided:**
- **Total Files**: Count of individual files in vault
- **Total Folders**: Count of directories
- **Vault Size**: Cumulative size of all files
- **Storage ID Length**: File storage identifier length (from settings)
- **Share Hash Length**: Share reference hash length (from settings)
- **Commit Author**: Git author for operations (from settings)
- **Commit Message**: Git commit message template (from settings)

## File Information

When called with a file path, displays metadata for that file.

**Example:**
```bash
./zep info documents/report.pdf
```

**Output:**
```
=== File Information ===
File: report.pdf
Path: documents/report.pdf
Type: File
Storage ID: a3f2e1c9d4b6f8e2
Encrypted Key: 8a4f2e1c...d4b6f8e2 (hex)
Size: 1.2 MB
```

**Information Provided:**
- **File**: Filename
- **Path**: Full vault path
- **Type**: "File" or "Folder"
- **Storage ID**: Random hex identifier for encrypted file
- **Encrypted Key**: AES-256-GCM encryption key (first and last 8 chars)
- **Size**: File size in human-readable format

## Folder Information

When called on a directory, displays folder metadata and contents summary.

**Example:**
```bash
./zep info documents
```

**Output:**
```
=== Folder Information ===
Path: documents
Type: Folder
Files: 12
Subdirectories: 3
Total Size: 450 MB
```

## Functions

### `GetVaultStats`

Calculate overall vault statistics.

**Function Signature:**
```go
func GetVaultStats(index *Index) (map[string]interface{}, error)
```

**Returns:**
- Total file count
- Total folder count
- Total size
- Settings information

### `PrintVaultInfo`

Display formatted vault statistics.

**Function Signature:**
```go
func PrintVaultInfo(stats map[string]interface{})
```

### `GetFileInfo`

Retrieve detailed information about a specific file or folder.

**Function Signature:**
```go
func GetFileInfo(index *Index, path string) (map[string]interface{}, error)
```

**Returns:**
- File/folder metadata
- Storage identifiers
- Size information
- Path information

### `PrintFileInfo`

Display formatted file information.

**Function Signature:**
```go
func PrintFileInfo(info map[string]interface{})
```

## Use Cases

### Vault Audit

Check vault size and file count before operations:
```bash
./zep info
```

### File Verification

Verify file metadata before sharing or deleting:
```bash
./zep info documents/sensitive.pdf
```

### Storage Management

Find which files are consuming the most space:
```bash
./zep ls documents
./zep info documents/large-file.zip
```

### Security Verification

Check file's encryption details and storage ID:
```bash
./zep info path/to/file.txt
# Review the Storage ID and Encrypted Key
```

## Implementation Details

### Size Calculation

Folder sizes are calculated by summing all contained files:
1. Traverse folder tree
2. Sum all file sizes
3. Return total in bytes
4. Format as human-readable (B, KB, MB, GB)

### Storage ID

The Storage ID is the random hex filename used in GitHub storage. It's:
- Generated randomly during upload
- Reused for file updates
- Length determined by settings `FileHashLength`

### Encrypted Key

The file's encryption key is:
- Derived from your vault password using PBKDF2
- Stored encrypted in the vault index
- Displayed truncated (first and last 8 chars) for verification

## Integration

Info is integrated into other operations:
- Upload completion shows storage ID
- Download verification uses file info
- Delete confirmation shows file size
- Share operations reference file metadata

## See Also

- [Settings Module](SETTINGS.md) - Vault configuration
- [Index Module](INDEX.md) - Vault structure and navigation
- [Encryption Module](ENCRYPTION.md) - File encryption details
