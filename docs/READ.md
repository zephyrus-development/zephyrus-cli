# Read Module

The Read module enables reading and displaying file contents directly from the vault without downloading to disk.

## Overview

The Read module allows you to view file contents in the terminal, useful for:
- Quickly viewing text files without saving to disk
- Verifying file contents before full download
- Sharing content via terminal without creating temporary files
- Reducing disk I/O for frequent lookups

## Command

### `read` - Display File Contents

Read and display the contents of a file directly from the vault.

**Usage:**
```bash
./zep read <vault-path> [--shared <share-string>]
```

**Aliases:** `cat`, `view`, `display`

**Arguments:**
- `vault-path`: Path to file in vault

**Flags:**
- `--shared <share-string>`: Read a shared file using share reference

## Examples

### Read from Your Vault

```bash
./zep read documents/notes.txt
```

**Output:**
```
Here are my notes for today...
Line 1
Line 2
...
```

### Read a Text File

```bash
./zep read config/settings.json
```

### Read a Shared File

Recipients can read shared files without downloading:

```bash
./zep read _ --shared "username:shareref:password:RmlsZW5hbWU="
```

### Pipeline Output

Combine with other commands:

```bash
# Search for files and read matches
./zep search notes | xargs -I {} ./zep read {}

# Process content with grep
./zep read documents/log.txt | grep ERROR
```

## Implementation

### Data Flow

1. **Find file** in vault index
2. **Fetch encrypted file** from GitHub
3. **Decrypt** using file's encryption key
4. **Display** to stdout

### File Format Handling

The Read module:
- Outputs raw file content to stdout
- Handles binary files (may appear as garbage)
- Respects terminal encoding
- Works with pipes and redirection

### Decryption Process

For vault files:
1. Locate file in index → get storage ID and encryption key
2. Fetch encrypted data from GitHub (file stored as hex ID)
3. Decrypt with PBKDF2-derived key using AES-256-GCM
4. Verify authentication tag (GCM)
5. Output to stdout

For shared files:
1. Fetch encrypted pointer file
2. Decrypt pointer with share password
3. Parse JSON to get file's storage ID and encryption key
4. Fetch and decrypt file (same as vault files)
5. Output to stdout

## Use Cases

### Documentation Review

```bash
# Read and review setup guide before proceeding
./zep read docs/README.md
```

### Quick Verification

```bash
# Check file size and type without full download
./zep read metadata/info.txt
```

### Content Search

```bash
# Find specific text in files
./zep read path/to/file.txt | grep "search term"
```

### Configuration Inspection

```bash
# View config files without temporary files
./zep read config/database.conf
```

### Log Analysis

```bash
# Tail recent entries
./zep read logs/app.log | tail -20

# Count log lines
./zep read logs/app.log | wc -l

# Filter errors
./zep read logs/error.log | grep FATAL
```

## Limitations

### Binary Files

Reading binary files (images, PDFs, archives) to terminal:
- May display as gibberish
- Not recommended for non-text files
- Use `download` instead for binary files

### Large Files

Large files read to terminal:
- Load entirely into memory
- May take time with large files
- No progress indicator during read
- Consider download for large files

### Encoding

Terminal encoding may affect display:
- UTF-8 files work best
- Non-UTF-8 content may display incorrectly
- Binary content always appears as gibberish

## Performance

### Network Operations

Each read operation:
1. Fetches encrypted file from GitHub (network I/O)
2. Decrypts locally (CPU-bound)
3. Outputs to stdout

Performance depends on:
- File size
- Network latency
- Terminal speed (if piped)

### Memory Usage

- Entire file loaded into memory
- Not suitable for extremely large files
- Consider download + streaming for large logs

## Security Considerations

### Terminal Output

Reading sensitive files displays content in terminal:
- Appears in shell history (consider `history -c`)
- May be logged by terminal multiplexer
- Could be captured by clipboard managers
- Consider environment before reading sensitive data

### Shared File Reading

When reading shared files:
- Share password is passed on command line (visible in process list)
- Considersecure sharing methods for sensitive content
- Share links should use secure distribution channels

## Integration

Read command works with:
- Vault index lookup
- Encryption/decryption
- Shared file system
- Terminal output

## Alternatives

| Task | Command | Benefit |
|------|---------|---------|
| Save to disk | `download` | Persistent file, better for binary |
| Search content | `search` | Find by filename/path |
| List files | `ls` | Browse vault structure |
| View metadata | `info` | File stats without content |

## See Also

- [Download Module](DOWNLOAD.md) - Save files to disk
- [Search Module](SEARCH.md) - Find files by name
- [Shared Files](SHARE.md) - Share via reference
- [Encryption Module](ENCRYPTION.md) - How decryption works
