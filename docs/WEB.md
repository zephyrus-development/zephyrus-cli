# WEB.md - Web Interface Documentation

## Overview

Zephyrus includes a browser-based web interface for accessing your encrypted vault without the CLI. The web interface provides secure file viewing, downloading, and sharing capabilities with zero-knowledge encryption.

**Key Features:**
- 🔐 Browser-based file decryption (encryption never leaves client)
- 📁 File browsing and navigation
- ⬇️ Direct file download to local machine
- 🔗 Shareable encrypted links
- 🔑 GitHub authentication support
- 🖥️ Responsive, modern UI

## Getting Started

### Prerequisites

1. GitHub account with SSH key configured
2. Zephyrus vault initialized with CLI: `zep setup`
3. Modern web browser (Chrome, Firefox, Safari, Edge)
4. No special installation required

### Accessing the Web Interface

The web interface is served from the GitHub Pages URL associated with your vault repository. The typical URL is:

```
https://<your-github-username>.github.io/zephyrus
```

Example:
```
https://isaac99.github.io/zephyrus
```

## Authentication

### Login Methods

The web interface supports two authentication methods:

#### Method 1: GitHub Authentication

```
1. Click "Login with GitHub"
2. Authorize Zephyrus application
3. Your username is auto-filled
4. Enter your vault password
5. Click "Connect"
```

**Advantages:**
- Quick login process
- Automatic username detection
- Secure OAuth flow

#### Method 2: Manual Entry

```
1. Enter GitHub username manually
2. Enter vault password
3. Click "Connect"
```

**Use cases:**
- If GitHub auth fails
- Accessing from public computers
- Testing with different accounts

### Session Management

- ✅ Session tokens stored locally in browser
- ✅ Session expires after browser close (privacy)
- ✅ No server-side session storage
- ✅ Each login creates fresh authentication

## File Browsing

### Navigation Interface

Once authenticated, you'll see:

```
Vault View:
├── 📁 Documents
│   ├── 📄 report.pdf
│   ├── 📄 notes.txt
│   └── 📁 Archived
│       └── 📄 2023-report.pdf
├── 📁 Photos
│   └── 📄 vacation.jpg
└── 📁 Backups
```

**Controls:**
- Click folder icon to expand/collapse
- Click file name to select
- Double-click folder to navigate into
- Click "Up" or "Back" to navigate out

### File Details

When a file is selected, you see:

- **File Name**: Original vault path (e.g., "documents/report.pdf")
- **File Size**: Unencrypted file size in bytes
- **Encryption Status**: "Encrypted with AES-256-GCM"
- **Modified Date**: Last update timestamp
- **Storage ID**: Anonymous hex identifier

## File Operations

### Download Files

**Single File Download:**

```
1. Browse to file in vault
2. Click "Download" button
3. Browser downloads decrypted file
4. File saved to Downloads folder
```

**Batch Download Directory:**

```
1. Select folder from file tree
2. Click "Download Directory"
3. All files recursively packaged
4. ZIP archive downloaded to local machine
```

### Encryption/Decryption Flow

**Complete End-to-End Encryption:**

```
Vault Storage (GitHub):
├─ Encrypted File: a3f2e1c9d4b6f8e2
├─ Encrypted Index: .config/index
├─ Encrypted Settings: .config/settings

Browser:
1. Fetch encrypted file content
2. Fetch encrypted key from index
3. Decrypt key using vault password
4. Decrypt file using decrypted key
5. Display/Save to user

Network:
- Only encrypted data transferred
- Decryption happens in browser memory
- Password never sent to server
- No plaintext stored remotely
```

**Security Guarantee:** Zephyrus developers cannot access your files because decryption happens entirely on your machine.

## Sharing Files

### Create Shareable Link

**Process:**

```
1. Browse to file in vault
2. Click "Share" or "Get Link"
3. System generates shareable URL
4. Copy link and share with others
```

**Generated Link Format:**

```
https://isaac99.github.io/zephyrus/share?token=abc123xyz...&file=documents/report.pdf
```

### Share Features

- ✅ Recipient can view/download file without vault access
- ✅ Password-protected encryption
- ✅ Expiration time (optional)
- ✅ One-time download option (optional)
- ✅ Download count limit (optional)

### Recipient Access

**Recipient Views Shared File:**

```
1. Receive shareable link
2. Click link in browser
3. View file details and preview (if available)
4. Enter password if require (sender's password)
5. Download file to local machine
```

**No Account Needed:**
- Recipients don't need GitHub account
- Recipients don't need Zephyrus installation
- Works entirely in browser

## User Interface Components

### Main Dashboard

```
┌─────────────────────────────────────┐
│  Zephyrus Vault Manager             │
├─────────────────────────────────────┤
│ User: isaac99                        │
│ Vault Size: 2.3 GB                  │
│ Files: 142                           │
├─────────────────────────────────────┤
│  📂 Folder Navigation                │
│  📄 File Browser                     │
│  ⚙️  Settings                        │
│  🔐 Security                         │
│  📤 Upload (web)                     │
│  🔗 Sharing                          │
└─────────────────────────────────────┘
```

### File Preview Panel

```
File: documents/report.pdf
Size: 2.4 MB
Type: PDF Document
Status: ✅ Decrypted
Last Modified: 2024-01-15 14:32

[📥 Download] [🔗 Share] [❌ Delete]
```

### Permission System

Files shown based on authentication:

- ✅ Authenticated user: Can see all vault files
- ✅ Shared link recipient: Can only see shared file
- ❌ Unauthenticated: Cannot access any files

## Supported File Types

### Viewable in Browser

These file types display preview in browser:

- **Images**: JPG, PNG, GIF, WebP, SVG
- **Documents**: PDF (with PDF.js viewer)
- **Text**: TXT, MD, JSON, CSV, XML
- **Code**: JS, TS, Go, Python, etc.

### Download Only

These file types can only be downloaded:

- Executables: EXE, APP, DMG
- Archives: ZIP, RAR, 7Z, TAR
- Media: MP4, MKV, AVI, MOV, MP3, WAV
- Binary: All other formats

## Browser Compatibility

**Supported Browsers:**

| Browser | Version | Support |
|---------|---------|---------|
| Chrome  | 90+     | ✅ Full Support |
| Firefox | 88+     | ✅ Full Support |
| Safari  | 14+     | ✅ Full Support |
| Edge    | 90+     | ✅ Full Support |
| Opera   | 76+     | ✅ Full Support |

**Required Features:**
- JavaScript ES6+
- WebCrypto API (for AES-256-GCM)
- LocalStorage (for session tokens)
- Fetch API (for data transfer)

## Security Features

### Encryption In Browser

```javascript
// Client-side decryption example
const encryptedData = await fetch(`/vault/a3f2e1c9d4b6f8e2`);
const key = await decryptKey(masterPassword);
const plaintext = await window.crypto.subtle.decrypt(
  key,
  encryptedData
);
```

### Password Never Transmitted

- ✅ Password used only for local key derivation (PBKDF2)
- ✅ Key derivation happens in browser
- ✅ Only encrypted keys sent over network
- ✅ Server never sees plaintext password

### Session Security

- ✅ Session tokens use cryptographically secure generation
- ✅ Tokens tied to browser instance
- ✅ Tokens auto-clear on browser close
- ✅ HTTPS-only transmission
- ✅ HttpOnly flags on sensitive cookies

## Advanced Features

### Vault Statistics

Dashboard shows:

- **Total Size**: Sum of all file sizes
- **File Count**: Number of files in vault
- **Folder Count**: Number of folders
- **Last Modified**: Most recent file change
- **Encryption Status**: All files encrypted

### Search Function

```
Search Bar:
├─ Full-text search (filename only)
├─ Filter by file type
├─ Filter by date range
├─ Real-time results
└─ Case-insensitive matching
```

### Mobile Interface

- ✅ Responsive design (works on tablets/phones)
- ✅ Touch-friendly controls
- ✅ Bottom navigation for mobile
- ✅ Optimized file preview
- ⚠️ Limited by mobile browser capabilities

## Troubleshooting

### Login Issues

**Problem: "Authentication failed"**
- Verify GitHub username is correct
- Check vault password is accurate
- Ensure SSH key is configured on GitHub
- Try clearing browser cache and cookies

**Problem: "Cannot fetch vault index"**
- Network connection issue
- GitHub API rate limit hit
- Try again in 1 minute

### File Access Issues

**Problem: "File not found in vault"**
- File may have been deleted
- Verify path is correct
- Check file permissions
- Refresh browser and retry

**Problem: "Decryption failed"**
- Vault password is incorrect
- File may be corrupted
- Try with correct password
- Use CLI to verify file integrity

### Download Issues

**Problem: "Download failed"**
- Check internet connection
- Check browser download settings
- Verify file permissions
- Try different browser

**Problem: "File size too large"**
- Browser memory limit reached
- Use CLI for very large files
- Try splitting files into smaller chunks

## Comparison: CLI vs Web Interface

| Feature | CLI | Web |
|---------|-----|-----|
| Upload Files | ✅ | ⚠️ (Limited) |
| Download Files | ✅ | ✅ |
| Browse Vault | ✅ | ✅ |
| Search Files | ✅ | ✅ |
| Share Files | ❌ | ✅ |
| Password Reset | ✅ | ❌ |
| Transfer Vault | ✅ | ❌ |
| Batch Operations | ✅ | ⚠️ |
| Mobile Support | ❌ | ✅ |

**Recommendation:**
- Use **CLI** for vault management and bulk operations
- Use **Web** for file access and sharing

## Getting Help

### Common Questions

**Q: Is my password stored in browser cache?**
- A: No. Password is only used for PBKDF2 key derivation in memory.

**Q: Can Zephyrus developers see my files?**
- A: No. All decryption happens on your machine, never on servers.

**Q: Does the website have server logs?**
- A: Minimal logging only for HTTP errors. No encryption keys or passwords logged.

**Q: Can I use the web interface offline?**
- A: No. You need internet to fetch encrypted files from GitHub.

**Q: Is web interface mobile-friendly?**
- A: Yes, fully responsive on phones and tablets.

### Support Resources

- GitHub Issues: [Report bugs or request features](https://github.com/your-username/zephyrus/issues)
- Documentation: See individual command docs in `/docs`
- CLI Help: `zep --help`

## Related Documentation

- [`SETUP.md`](SETUP.md): Initial vault setup
- [`AUTH.md`](AUTH.md): Authentication and session management
- [`UPLOAD.md`](UPLOAD.md): CLI file upload
- [`DOWNLOAD.md`](DOWNLOAD.md): CLI file download
- [`SHARE.md`](SHARE.md): Sharing in CLI
