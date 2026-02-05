# Transfer Vault Feature

## Overview

The `transfer-vault` command allows you to copy all files and folders from one Zephyrus vault to another. This is useful for:
- Migrating to a new GitHub account
- Creating vault backups in another account
- Consolidating multiple vaults
- Sharing vault data with another user

## Command Syntax

```bash
zep transfer-vault <source-username> <dest-username>
zep transfer <source-username> <dest-username>
zep xfer <source-username> <dest-username>
zep copy-vault <source-username> <dest-username>
```

## Usage

### Basic Transfer

```bash
zep transfer-vault alice bob
```

### Step-by-Step

1. **Enter source vault password**:
   ```
   Source vault authentication (alice):
   Source Vault Password: ••••••••••
   ```

2. **Enter destination vault password**:
   ```
   Destination vault authentication (bob):
   Destination Vault Password: ••••••••••
   ```

3. **Confirm the transfer**:
   ```
   ⚠️  You are about to transfer all files from alice to bob.
   This will copy all vault contents. Continue? (y/n): y
   ```

4. **Monitor progress**:
   ```
   🔄 Starting vault transfer from alice to bob
   [1/5] Authenticating with source vault...
   [2/5] Fetching destination vault key...
   [3/5] Scanning source vault...
   Found 42 files to transfer
   
   Transferring file (1): documents/report.pdf
     → Transferred: documents/report.pdf (a3f2e1c9d4b6f8e2)
   Transferring file (2): images/photo.jpg
     → Transferred: images/photo.jpg (b4g3f2d0e5c7h9f3)
   
   [4/5] Preparing transfer package...
   [5/5] Uploading files to destination vault...
   
   ✔ Successfully transferred 42 files from alice to bob
   ```

## How It Works

### Process

1. **Authentication**: Authenticates with both source and destination vaults
2. **Scanning**: Walks through the entire source vault index
3. **Decryption**: Decrypts each file with the source password
4. **Re-encryption**: Re-encrypts each file with a new key using the destination password
5. **Upload**: Uploads all files to the destination vault in a single batch

### Key Features

- ✅ **Full Directory Preservation**: Maintains all folder structures
- ✅ **New Encryption Keys**: Each file gets fresh encryption keys for the destination
- ✅ **Atomic Operation**: All-or-nothing transfer (fails completely if any file has issues)
- ✅ **Progress Tracking**: Real-time feedback for each file transferred
- ✅ **Error Handling**: Detailed error messages if transfer fails

## Requirements

### Prerequisites

- Both source and destination vaults must exist on GitHub
- Must know the passwords for both vaults
- Source vault must contain at least one file
- SSH access to destination vault (deploy key with write permissions)

### Vault Setup

Ensure both vaults are properly initialized:

```bash
# Setup source vault
zep setup alice ~/.ssh/id_ed25519

# Setup destination vault (if not already done)
zep setup bob ~/.ssh/id_ed25519
```

## Examples

### Example 1: Backup to Secondary Account

```bash
# Transfer production vault to backup account
zep transfer-vault production-user backup-user

# Verify transfer by listing destination
zep list -u backup-user
```

### Example 2: Consolidate Multiple Vaults

```bash
# Merge alice's vault into bob's vault
zep transfer-vault alice bob

# Later, transfer charlie's vault too
zep transfer-vault charlie bob
```

### Example 3: Using Aliases

```bash
# These all do the same thing:
zep transfer-vault alice bob
zep transfer alice bob
zep xfer alice bob
zep copy-vault alice bob
```

## Security Considerations

### What Gets Encrypted

- ✅ All file contents (re-encrypted with destination password)
- ✅ All file keys (new keys generated per file)
- ✅ Directory structure (encrypted in vault index)

### What's Safe

- 🔐 Passwords are never transmitted (only stored in memory during transfer)
- 🔐 Files are decrypted locally, then re-encrypted
- 🔐 Destination files use completely new encryption keys
- 🔐 GitHub only sees encrypted data

### Best Practices

1. **Use Strong Passwords**: Ensure both vault passwords are strong (12+ characters)
2. **Verify After Transfer**: List destination vault to confirm all files transferred
3. **Keep SSH Keys Secure**: Both accounts need secure SSH deploy keys
4. **Plan Timing**: Transfer during off-hours if vaults are large
5. **Have Backup**: Maintain backup of source vault before transferring

## Troubleshooting

### "Source and destination vaults must be different"

```
❌ Source and destination vaults must be different.
```

**Solution**: Use different usernames for source and destination.

### "Failed to authenticate with source vault"

```
❌ Transfer failed: failed to authenticate with source vault: auth failed: invalid password
```

**Solution**: Verify the source vault password is correct.

### "Destination vault not found"

```
❌ Transfer failed: destination vault not found
```

**Solution**: Ensure destination vault exists on GitHub and is properly set up.

### "Failed to decrypt destination vault key"

```
❌ Transfer failed: failed to decrypt destination vault key (invalid password)
```

**Solution**: Verify the destination vault password is correct.

### "No files found in source vault"

```
❌ Transfer failed: no files found in source vault
```

**Solution**: Source vault is empty. Upload some files first before transferring.

## Verifying Transfer

After transfer completes, verify success:

```bash
# List files in source
zep list -u alice

# List files in destination
zep list -u bob

# Download a specific file from destination to verify
zep download -u bob documents/report.pdf ./report.pdf
```

## Performance Notes

- **Large Vaults**: Transfer time depends on vault size and number of files
- **Each File**: Decrypted and re-encrypted individually
- **Memory Usage**: Minimal (files processed one at a time, not buffered)
- **Network**: Single batch upload after all files processed

### Example Times

- 10 files under 100MB: ~5-10 seconds
- 50 files under 500MB: ~20-30 seconds
- 100+ files or large files: 1-5 minutes

## Limitations

- ⚠️ **Linear Processing**: Files are transferred one at a time (no parallelization)
- ⚠️ **Direction**: Transfer is one-way (source → destination)
- ⚠️ **Existing Files**: Destination files are overwritten if paths conflict
- ⚠️ **Network Required**: Requires stable internet connection for entire duration

## Related Commands

- [`zep upload`](UPLOAD.md): Upload individual files
- [`zep download`](DOWNLOAD.md): Download files
- [`zep reset-password`](AUTH.md): Change vault password
- [`zep connect`](AUTH.md): Create persistent session
- [`zep list`](LIST.md): List vault contents
