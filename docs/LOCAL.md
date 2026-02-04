# Local Filesystem Module

The Local module provides commands to access the local filesystem from within the Zephyrus interactive shell, enabling convenient filesystem browsing without exiting the REPL.

## Overview

The Local module bridges the gap between the REPL shell and the local filesystem, allowing you to:
- List local files and directories
- Navigate the local filesystem
- Verify file availability before upload
- Monitor local directory structures
- All without exiting the interactive shell

## Commands

### `localls` - List Local Files (ls equivalent)

List files in the local filesystem using the `ls` command.

**Usage:**
```bash
localls [arguments]
```

**Aliases:** `lls`

**Arguments:**
Pass any arguments supported by your system's `ls` (Unix/Linux) or `dir` (Windows)

**Behavior:**
- **Linux/macOS**: Runs `ls` directly with provided arguments
- **Windows**: Attempts to run `ls` if available (Git Bash, WSL), falls back to `dir`

**Examples:**

```bash
# List current directory
zep> lls
file1.txt
file2.pdf
Documents/

# List with detailed output (Linux/macOS)
zep> lls -la
drwxr-xr-x  5 user group   160 Feb  4 19:13 .
drwxr-xr-x 14 user group   448 Feb  3 19:52 ..
-rw-r--r--  1 user group  1024 Feb  4 10:30 file1.txt
-rw-r--r--  1 user group 52384 Feb  3 15:22 file2.pdf

# List specific directory
zep> lls Documents
report.pdf
budget.xlsx

# List with wildcard (Unix/Linux)
zep> lls *.txt
notes.txt
report.txt
```

### `localdir` - List Local Files (dir equivalent)

List files in the local filesystem using detailed directory listing.

**Usage:**
```bash
localdir [arguments]
```

**Aliases:** `ldir`

**Arguments:**
Pass any arguments supported by your system's `dir` (Windows) or `ls -la` (Unix/Linux)

**Behavior:**
- **Windows**: Runs `dir` with provided arguments
- **Linux/macOS**: Runs `ls -la` for equivalent detailed output

**Examples:**

```bash
# List current directory with details
zep> ldir

# Windows output:
 Directory of C:\Users\YourName\Documents
04/02/2026  19:13    <DIR>          .
04/02/2026  19:13    <DIR>          ..
04/02/2026  15:22            52,384 file2.pdf
04/02/2026  10:30             1,024 file1.txt

# List specific directory
zep> ldir Downloads

# Windows file attributes
zep> ldir /A
```

## Use Cases

### Pre-Upload Verification

Verify files exist before uploading:

```bash
zep> lls Documents
report.pdf
budget.xlsx

zep> upload Documents/report.pdf vault/reports/q1.pdf
✔ Upload successful.
```

### Finding Files to Share

Locate files in Documents before sharing:

```bash
zep> lls
Documents/
Desktop/
Downloads/

zep> lls Documents | grep "confidential"
confidential-report.pdf

zep> share Documents/confidential-report.pdf
```

### Directory Navigation

Monitor local directory structure:

```bash
zep> ldir
zep> lls Documents
zep> ldir Downloads
```

### Script Preparation

Prepare backup batches before uploading:

```bash
zep> lls backups/
backup-2026-01-01.zip
backup-2026-02-01.zip

zep> upload backups/backup-2026-02-01.zip vault/monthly/feb.zip
```

## Cross-Platform Behavior

### Windows

```bash
# localls with fallback to dir
zep> lls Documents
(outputs using dir command)

# localdir uses dir
zep> ldir
 Directory of C:\Users\YourName\Documents
...

# With Git Bash installed, lls may use actual ls
zep> lls -l
```

### Linux/macOS

```bash
# localls uses ls
zep> lls Documents
file1.txt
file2.pdf

# localls with flags
zep> lls -lah
drwxr-xr-x  3 user  group  96 Feb  4 19:13 .
-rw-r--r--  1 user  group 1.1K Feb  4 10:30 file1.txt
-rw-r--r--  1 user  group 51K Feb  3 15:22 file2.pdf

# localdir uses ls -la
zep> ldir
drwxr-xr-x  3 user  group  96 Feb  4 19:13 .
-rw-r--r--  1 user  group 1.1K Feb  4 10:30 file1.txt
```

## Implementation

### Command Execution

Commands use `os/exec` to spawn local processes:

**localls:**
```go
cmd := exec.Command("ls", args...)  // Unix/Linux
cmd := exec.Command("dir", args...) // Windows fallback
```

**localdir:**
```go
cmd := exec.Command("dir", args...)        // Windows
cmd := exec.Command("ls", append([]string{"-la"}, args...)...) // Unix/Linux
```

### Error Handling

- If command not found: Error message displayed
- If directory doesn't exist: System error shown
- If insufficient permissions: Access denied error shown

### Stdin/Stdout/Stderr

- `Stdin`: Connected to terminal (interactive)
- `Stdout`: Connected to terminal (visible output)
- `Stderr`: Connected to terminal (error messages)

## Limitations

### Command Availability

- **localls**: Requires `ls` command (Windows: requires Git Bash/WSL or falls back to `dir`)
- **localdir**: Requires `dir` (Windows) or `ls -la` (Unix/Linux)

### Arguments

Arguments passed through to underlying OS command:
- Platform-specific flags may not work on other platforms
- `-la` works on Unix/Linux, not Windows `dir`
- `/A` works on Windows `dir`, not Unix/Linux `ls`

### Special Characters

Paths with spaces or special characters:
- Quote paths on command line: `lls "My Documents"`
- Shell escaping applies

## REPL-Only

These commands only work within the interactive shell:
- Not available in stateless mode (`./zep upload ...`)
- Not available outside REPL environment
- Only useful for local filesystem operations

## Examples Workflow

```bash
# Start REPL
./zep
# Username: myuser
# Password: ••••••••••
# ✔ Welcome, myuser. Session Active.

# Check what files we have locally
zep> lls
Documents/
Documents-2/
backup.zip

# List detailed info
zep> ldir Documents

# Upload a file after verifying it exists
zep> upload Documents/important.pdf vault/documents/important.pdf
✔ Upload successful.

# Check another directory
zep> ldir backup.zip
(shows file details)

# Upload backup
zep> upload backup.zip vault/backups/latest.zip
✔ Upload successful.

# Verify vault contents
zep> ls vault/
NAME              TYPE
----              ----
documents/        [FOLDER]
backups/          [FOLDER]

zep> exit
```

## See Also

- [REPL Shell](../main.go) - Interactive mode documentation
- [Upload Module](UPLOAD.md) - Uploading files to vault
- [List Module](LIST.md) - Listing vault contents
