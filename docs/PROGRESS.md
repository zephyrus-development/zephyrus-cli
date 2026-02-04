# Progress Module

The Progress module provides visual feedback during long-running vault operations, including multi-step progress indicators with messages and completion status.

## Overview

Progress indicators improve user experience by showing operation status and preventing timeout confusion during network operations. All major vault operations (upload, download, share, delete, purge) use progress feedback.

## Functions

### `PrintProgressStep`

Display a progress step with step number and message.

**Function Signature:**
```go
func PrintProgressStep(step int, totalSteps int, message string)
```

**Parameters:**
- `step`: Current step number (1-indexed)
- `totalSteps`: Total number of steps in operation
- `message`: Status message to display

**Example Output:**
```
[1/5] Encrypting file...
[2/5] Uploading to GitHub...
[3/5] Updating vault index...
[4/5] Finalizing...
[5/5] Complete!
```

### `PrintCompletionLine`

Display a completion message with checkmark.

**Function Signature:**
```go
func PrintCompletionLine(message string)
```

**Parameters:**
- `message`: Success message to display

**Example Output:**
```
✔ Upload successful.
```

### `PrintErrorLine`

Display an error message with error symbol.

**Function Signature:**
```go
func PrintErrorLine(message string)
```

**Parameters:**
- `message`: Error message to display

**Example Output:**
```
❌ Upload failed: insufficient permissions
```

### `ClearProgress`

Clear progress lines from terminal (currently a no-op for compatibility).

**Function Signature:**
```go
func ClearProgress()
```

## Usage Examples

### Upload Operation

```go
func Upload(session *Session, vaultPath string, localPath string) error {
    PrintProgressStep(1, 5, "Reading file...")
    // Read file...
    
    PrintProgressStep(2, 5, "Generating encryption key...")
    // Generate key...
    
    PrintProgressStep(3, 5, "Encrypting file content...")
    // Encrypt...
    
    PrintProgressStep(4, 5, "Pushing to GitHub...")
    // Push...
    
    PrintProgressStep(5, 5, "Updating vault index...")
    // Update index...
    
    PrintCompletionLine("Upload successful.")
    return nil
}
```

### Error Handling

```go
err := Download(session, vaultPath, localPath)
if err != nil {
    PrintErrorLine(fmt.Sprintf("Download failed: %v", err))
    return err
}
```

## Operations Using Progress

| Operation | Steps | Purpose |
|-----------|-------|---------|
| Upload | 5 | Read → Encrypt → Upload → Index → Complete |
| Download | 5 | Fetch → Decrypt → Write → Verify → Complete |
| Share | 4 | Fetch → Encrypt → Share → Index |
| Delete | 4 | Fetch → Remove → Push → Update |
| Purge | 3 | Prepare → Wipe → Clean |

## Implementation Details

### Step Display Format

Steps are displayed as `[current/total]` followed by the message:
```
[1/5] Message here...
[2/5] Next step...
```

### Symbols

- ✔ (checkmark) - Success/completion
- ❌ (cross mark) - Error/failure

### Terminal Output

Progress messages are printed to standard output using `fmt.Println()` and `fmt.Printf()`. The terminal handles line wrapping and character encoding.

## Performance Considerations

Progress updates introduce minimal overhead:
- Single `fmt.Print` call per step
- No network delays
- Straightforward string formatting

The visual feedback improvement outweighs the negligible performance cost.

## Future Enhancements

Potential improvements to progress module:
- Progress bars with percentage completion
- Spinner animations for indeterminate operations
- Custom progress output formats
- Suppression flag for non-interactive mode
- Speed/rate display (e.g., "2.3 MB/s")
- Time estimates for long operations

## See Also

- [Upload Module](UPLOAD.md) - File upload with progress
- [Download Module](DOWNLOAD.md) - File download with progress
- [Share Module](SHARE.md) - File sharing with progress
- [Delete Module](DELETE.md) - File deletion with progress
- [Purge Module](PURGE.md) - Vault wipe with progress
