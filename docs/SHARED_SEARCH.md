# Shared File Search Module

The Shared File Search module provides fuzzy matching and search functionality for managing shared files by name, allowing discovery and revocation without memorizing share reference IDs.

## Overview

The Shared Search module enables finding and filtering shared files by:
- Exact filename matches
- Prefix matches
- Substring matches
- Fuzzy scoring for relevance

This makes it easy to locate shares for management without needing to remember or look up reference IDs.

## Search Functionality

### Match Types

The search algorithm uses weighted scoring:

| Match Type | Score | Example |
|-----------|-------|---------|
| Exact match | 0 (best) | `report.pdf` matches `report.pdf` |
| Prefix match | ~filename length | `rep` matches `report.pdf` |
| Substring match | ~100 + position | `port` matches `report.pdf` |
| No match | ∞ (excluded) | `xyz` doesn't match anything |

### Sort Order

Results are sorted by score (lowest = best match):
1. Exact matches first
2. Prefix matches second
3. Substring matches last

## Functions

### `FindSharedFilesByName`

Search shared files with fuzzy matching.

**Function Signature:**
```go
func FindSharedFilesByName(index *SharedIndex, pattern string) ([]*SharedFileMatch, error)
```

**Parameters:**
- `index`: Shared index to search
- `pattern`: Search pattern (filename, path, or partial)

**Returns:**
- Slice of `SharedFileMatch` structs (sorted by score)
- Error if search fails

### `SharedFileMatch` Structure

```go
type SharedFileMatch struct {
    ShareID    string
    Filename   string
    VaultPath  string
    CreatedAt  string
    Score      int
}
```

**Fields:**
- `ShareID`: Reference ID for the share
- `Filename`: Original filename
- `VaultPath`: Full vault path
- `CreatedAt`: Creation timestamp
- `Score`: Match quality (lower is better)

### `PrintSharedFilesFormatted`

Display search results in formatted table.

**Function Signature:**
```go
func PrintSharedFilesFormatted(matches []*SharedFileMatch)
```

**Parameters:**
- `matches`: Search results to display

**Output Format:**
```
ID        FILENAME           VAULT PATH                     CREATED
----      --------           ----------                     -------
72cTWg    report.pdf         documents/reports/q1.pdf       2026-02-04
AbXkLm    budget.xlsx        financial/budget.xlsx          2026-02-03
```

## Command Integration

### `shared ls` - List Shares (with optional search)

List all shares or search by pattern.

**Usage:**
```bash
./zep shared ls [pattern]
```

**Aliases:** `shared list`, `shared find`, `shared search`

**Examples:**

#### List all shares

```bash
zep> shared ls
ID        FILENAME          VAULT PATH
----      --------          ----------
72cTWg    report.pdf        documents/reports/q1.pdf
AbXkLm    budget.xlsx       financial/budget.xlsx
```

#### Search by exact filename

```bash
zep> shared ls report
```

**Output:**
```
ID        FILENAME      VAULT PATH
----      --------      ----------
72cTWg    report.pdf    documents/reports/q1.pdf
```

#### Search by prefix

```bash
zep> shared ls rep
```

**Output:**
```
ID        FILENAME      VAULT PATH
----      --------      ----------
72cTWg    report.pdf    documents/reports/q1.pdf
```

#### Search by substring

```bash
zep> shared ls port
```

**Output:**
```
ID        FILENAME      VAULT PATH
----      --------      ----------
72cTWg    report.pdf    documents/reports/q1.pdf
```

#### Search by vault path

```bash
zep> shared ls documents
```

**Output:**
```
ID        FILENAME      VAULT PATH
----      --------      ----------
72cTWg    report.pdf    documents/reports/q1.pdf
```

## Search Algorithms

### Exact Match

Checks if pattern equals filename:
```
Pattern: "report.pdf"
Filename: "report.pdf"
Match: Yes (score 0)
```

### Prefix Match

Checks if filename starts with pattern:
```
Pattern: "report"
Filename: "report.pdf"
Match: Yes (score 6 = length of pattern)
```

### Substring Match

Checks if pattern appears in filename:
```
Pattern: "port"
Filename: "report.pdf"
Match: Yes (score 104 = 100 + position)
```

### Case Insensitive

All comparisons are case-insensitive:
```
Pattern: "REPORT"
Filename: "report.pdf"
Match: Yes (exact match, score 0)
```

## Use Cases

### Find Shares by Extension

Search for all PDFs:

```bash
zep> shared ls .pdf
ID        FILENAME         VAULT PATH
----      --------         ----------
72cTWg    report.pdf       documents/reports/q1.pdf
AbXkLm    proposal.pdf     projects/2024/proposal.pdf
```

### Find by Project Name

Find all shares from a specific project:

```bash
zep> shared ls 2024
ID        FILENAME          VAULT PATH
----      --------          ----------
AbXkLm    q1_budget.xlsx    2024/financial/q1_budget.xlsx
Cd2nXp    annual_plan.pdf   2024/planning/annual_plan.pdf
```

### Find by Department

Search shares from a department:

```bash
zep> shared ls financial
ID        FILENAME           VAULT PATH
----      --------           ----------
AbXkLm    budget.xlsx        financial/budget.xlsx
```

### Narrow Down Results

Search progressively:

```bash
zep> shared ls report
# Shows 5 results

zep> shared ls report2024
# Shows 2 results

zep> shared ls report2024q
# Shows 1 result: report2024_q1.pdf
```

## Integration with Management

### Find and Revoke

```bash
# Find the share
zep> shared ls report
ID        FILENAME
----      --------
72cTWg    report.pdf

# Revoke it
zep> shared rm 72cTWg
✔ Revoked share: 72cTWg
```

### Find and Get Info

```bash
# Find by name
zep> shared ls budget
ID        FILENAME
----      --------
AbXkLm    budget.xlsx

# Get details
zep> shared info AbXkLm
Share ID: AbXkLm
Filename: budget.xlsx
Vault Path: financial/budget.xlsx
...
```

### Find and Update

```bash
# Find old version
zep> shared ls document
ID        FILENAME
----      --------
72cTWg    document.pdf (old)

# Revoke old
zep> shared rm 72cTWg

# Upload new version
zep> upload document-v2.pdf documents/document.pdf

# Share new version
zep> share documents/document.pdf
```

## Performance

### Search Complexity

- Time: O(n) where n = number of shares
- Space: O(m) where m = number of matches
- Typical: <1ms for <100 shares

### Optimization Tips

If vault has many shares:
- Use specific patterns to narrow results
- Remember share ID for direct access
- Use full paths for unambiguous matches

## Limitations

### Search Scope

Search only covers:
- Filenames
- Vault paths
- Does NOT search file contents

### Exact Reference Required

To revoke by reference (fastest):
```bash
zep> shared rm 72cTWg  # Fast - direct lookup
zep> shared rm report  # Slower - needs search
```

### No Regex

Pattern matching is simple substring/prefix:
- No regex support
- No pattern wildcards
- Case-insensitive only

## Potential Enhancements

Future improvements could include:
- Regex pattern support
- Date range filtering
- Access count filtering
- Sort by creation date or size
- Export search results
- Bulk revocation of matching shares

## See Also

- [Shared Manage Module](SHARED_MANAGE.md) - Revoking shares
- [Shared Index Module](SHARED_INDEX.md) - Index structure
- [Share Module](SHARE.md) - Creating shares
- [Shared Files Overview](../README.md#secure-file-sharing)
