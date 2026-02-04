# Shared File Management Module

The Shared File Management module provides commands for managing shared files, including revoking access and updating share metadata.

## Overview

The Shared Management module enables vault owners to:
- View all active shares
- Get detailed information about shares
- Revoke access to shared files
- Update share information
- Manage share lifecycle

## Command

### `shared rm` - Revoke Share or Remove by Name

Revoke access to a shared file or remove by filename. This command supports both hash-based and name-based revocation.

**Usage:**
```bash
./zep shared rm <reference-or-name>
```

**Aliases:** `shared revoke`, `shared delete`, `shared remove`

**Arguments:**
- `reference-or-name`: Either:
  - Share reference ID (e.g., `72cTWg`) - exact match
  - Filename (e.g., `report.pdf`) - searches all shares
  - Full path (e.g., `documents/report.pdf`) - full match

## Behavior

### Hash-Based Revocation

When using a share reference:

```bash
zep> shared rm 72cTWg
✔ Revoked share: 72cTWg (report.pdf)
```

**Process:**
1. Look up share by ID
2. Confirm revocation
3. Remove from shared index
4. Delete file from GitHub
5. Push changes

### Name-Based Revocation

When using a filename:

```bash
zep> shared rm report.pdf
```

**Process:**
1. Search shared index by filename
2. **If exact match found**: Revoke that share
3. **If multiple matches found**: Show options, ask user to clarify
4. **If no matches found**: Error message

### Path-Based Revocation

Full vault path matches are prioritized:

```bash
zep> shared rm documents/report.pdf
✔ Found 1 match. Revoking: documents/report.pdf
```

## Examples

### Revoke by Reference ID

```bash
# Exact share ID
zep> shared rm 72cTWg
✔ Revoked share: 72cTWg
```

### Revoke by Filename

```bash
# Simple filename (searches all shares)
zep> shared rm report.pdf
✔ Revoked share: report.pdf from documents/report.pdf
```

### Revoke by Full Path

```bash
# Full vault path
zep> shared rm documents/reports/2024/q1.pdf
✔ Revoked share: q1.pdf from documents/reports/2024/q1.pdf
```

### Ambiguous Names (Multiple Matches)

If filename matches multiple shares:

```bash
zep> shared rm report.pdf
⚠️  Multiple matches found:
  1. documents/reports/q1.pdf (share: 72cTWg)
  2. archived/reports/q1.pdf (share: AbXkLm)

Which one? (1-2 or share ID): 1
✔ Revoked share: 72cTWg
```

## Functions

### `RevokeSharedFileByID`

Revoke a share using its reference ID.

**Function Signature:**
```go
func RevokeSharedFileByID(session *Session, shareID string) error
```

**Parameters:**
- `session`: Current vault session
- `shareID`: Share reference to revoke

**Returns:**
- Error if share not found or revocation fails

**Process:**
1. Load shared index
2. Find share by ID
3. Delete file from GitHub
4. Remove from index
5. Save index
6. Return success

### `RevokeSharedFileByName`

Revoke share(s) by filename with fuzzy matching.

**Function Signature:**
```go
func RevokeSharedFileByName(session *Session, namePattern string) error
```

**Parameters:**
- `session`: Current vault session
- `namePattern`: Filename or path to match

**Returns:**
- Error if no matches or revocation fails

**Process:**
1. Load shared index
2. Search for matches (exact, prefix, substring)
3. If one match: revoke it
4. If multiple matches: prompt user
5. Revoke selected share
6. Return success

## Integration with Shared Commands

### View Shares

List all shares before revoking:

```bash
zep> shared ls
ID        FILENAME          CREATED
----      --------          -------
72cTWg    report.pdf        2026-02-04
AbXkLm    budget.xlsx       2026-02-03
```

### Get Share Info

Get details before revoking:

```bash
zep> shared info 72cTWg
Share ID: 72cTWg
Filename: report.pdf
Vault Path: documents/reports/q1.pdf
Created: 2026-02-04 15:30:00
Access Count: 0
```

### Revoke and Verify

```bash
zep> shared ls
ID        FILENAME
----      --------
72cTWg    report.pdf

zep> shared rm 72cTWg
✔ Revoked share: 72cTWg

zep> shared ls
(no entries)
```

## Error Handling

### Share Not Found

```bash
zep> shared rm invalid123
❌ Share not found: invalid123
```

### Filename Not Found

```bash
zep> shared rm nonexistent.pdf
❌ No shares found matching: nonexistent.pdf
```

### Ambiguous Match

```bash
zep> shared rm report.pdf
⚠️  Multiple matches found for "report.pdf"
Use full path or share ID for clarity
```

## Revocation Effects

### Immediate

- Share is removed from shared index
- File deleted from GitHub (share pointer)
- Recipients can no longer access

### Not Immediate

- Shared file content remains encrypted on GitHub (original file)
- Revocation only removes the share pointer, not the file
- Re-sharing same file creates new share ID

## Use Cases

### Temporary Sharing

Share a file, revoke after review period:

```bash
# Initial share
zep> share documents/proposal.pdf
(share link sent to reviewer)

# After review, revoke
zep> shared rm proposal.pdf
✔ Revoked access
```

### Access Control

Revoke shares when permissions change:

```bash
# Employee leaving company
zep> shared rm confidential.pdf
✔ Revoked share
```

### Mistake Correction

Revoke and re-share with different password:

```bash
zep> shared rm budget.xlsx
zep> share financial/budget.xlsx
Enter new share password...
```

### File Rotation

Revoke old shares when file is updated:

```bash
zep> shared rm documents/policy.pdf
zep> upload policy.pdf documents/policy.pdf
zep> share documents/policy.pdf
```

## Security Considerations

### Timing

- Revocation is immediate
- All recipients lose access when revoked
- No way to recover a revoked share link
- Consider before revoking active shares

### Backup

- Revoked shares cannot be unrevoked
- No undo operation
- Keep records of who accessed what

### Re-sharing

- Can re-share same file with new ID
- Old share links become invalid
- New share requires new password

## Performance

### Lookup Speed

- Reference ID lookup: O(1) - instant
- Filename lookup: O(n) - linear search of shares
- Usually negligible (typical vault has <100 shares)

### Revocation Speed

- Depends on file size
- Network dependent (GitHub push)
- Usually completes in seconds

## See Also

- [Shared Files Overview](SHARE.md) - Creating shares
- [Shared Search Module](SHARED_SEARCH.md) - Finding shares
- [Shared Index Module](SHARED_INDEX.md) - Index structure
- [Share Command Reference](../README.md#shared---manage-shared-files)
