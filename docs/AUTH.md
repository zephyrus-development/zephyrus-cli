# auth.go Documentation

## Package utils

This package provides functionalities for user authentication, session management, and password operations in the Zephyrus CLI application.

### Constants

- **configPath**: The path to the configuration file (`zephyrus.conf`).

### Variables

- **globalSession**: A pointer to a `Session` struct that stores the session in RAM for REPL/Stateless mode.

### Types

#### Session

The `Session` struct holds the following fields:

- **Username**: The username of the authenticated user.
- **Password**: The password of the authenticated user.
- **RawKey**: The raw key used for encryption/decryption.
- **Index**: The vault index associated with the session.

### Functions

#### SetGlobalSession

```go
func SetGlobalSession(s *Session)
```

Injects a session into RAM, used by the REPL.

#### Connect

```go
func Connect(username string, password string) error
```

Initializes the session and syncs the index locally. It takes the username and password as parameters and returns an error if the connection fails.

#### (s *Session) Save

```go
func (s *Session) Save() error
```

Saves the current session to the configuration file. Returns an error if the save operation fails.

#### GetSession

```go
func GetSession() (*Session, error)
```

Checks memory first (REPL cache) for an active session, then falls back to disk. Returns the session and an error if not connected.

#### Disconnect

```go
func Disconnect() error
```

Clears the memory cache and removes the configuration file.

#### FetchSessionStateless

```go
func FetchSessionStateless(username string, password string) (*Session, error)
```

Performs authentication and index fetch without saving to disk. It returns a session and an error if authentication fails.

---

## Password Reset

### ResetPassword

```go
func ResetPassword(session *Session, newPassword string) error
```

Changes the vault password and re-encrypts all vault data with the new password. This function performs complete re-encryption of sensitive data.

**Parameters:**
- `session`: The active session with authenticated user credentials
- `newPassword`: The new password to be used for vault encryption

**Return:**
- `error`: Returns error if any re-encryption or push operation fails

### Password Reset Process

The password reset operation is a 5-step process that ensures all vault data is properly re-encrypted:

#### Step 1: Validate Input
- Verifies new password meets security requirements (minimum length)
- Checks password is not empty

#### Step 2: Re-encrypt Master Key
- Derives new encryption key from new password using PBKDF2 (100,000 iterations)
- Re-encrypts the vault's master encryption key with new key derivation
- Updates session password to reflect new password

#### Step 3: Update All File Keys
- Recursively processes entire vault index (including nested folders)
- For each file in vault:
  - Decrypts original file key with old password
  - Re-encrypts file key with new password
  - Updates index entry with re-encrypted key
- Updates per-file encryption keys throughout entire vault structure

#### Step 4: Re-encrypt Vault Components
- Encrypts updated vault index with new password
- Encrypts settings file with new password
- Encrypts shared index with new password
- Prepares batch push package with all updated components

#### Step 5: Push to Remote
- Uploads all re-encrypted vault components to GitHub
- Uses SSH authentication from session
- Commit message: "Nexus: Password Reset"
- Local session password updated for future operations

### Usage Examples

**Command Line:**

```bash
# Reset password (will prompt for current and new password)
zep reset-password

# Follow interactive prompts:
# 1. Enter current password (for verification)
# 2. Enter new password
# 3. Confirm new password (must match)
```

### Interactive Prompts

```
Enter current password for verification: ••••••••
Enter new password: ••••••••
Confirm new password: ••••••••

Processing vault re-encryption...
[1/5] Validating password...
[2/5] Updating master key...
[3/5] Re-encrypting all file keys...
[4/5] Re-encrypting vault components...
[5/5] Uploading changes to GitHub...

✔ Password successfully reset
```

### Security Considerations

**⚠️ Important Security Notes:**

1. **Data Integrity**: Every piece of encrypted data in vault is re-encrypted
   - Master encryption key updated
   - All individual file keys updated
   - Index structure preserved, only encryption changes

2. **One-Way Operation**: Cannot revert to old password
   - New password replaces old password permanently
   - Keep backup of new password in secure location
   - Password is not recoverable if forgotten

3. **Time Requirement**: Password reset may take time for large vaults
   - Must re-encrypt every file in vault
   - Large vaults with thousands of files may take several minutes
   - Do not interrupt process once started

4. **Network Dependency**: Requires complete push to GitHub
   - Must maintain stable network connection during reset
   - If connection drops mid-reset, vault may be in inconsistent state
   - Retry operation if network error occurs

5. **Session Consistency**: Local session password is updated
   - After successful reset, all future vault operations use new password
   - If you have multiple local sessions, all must be updated
   - Stateless mode (`zep -u` flag) will use new password automatically

### Use Cases

1. **Regular Security Rotation**
   - Change password periodically for security best practices
   ```bash
   zep reset-password
   ```

2. **Compromised Password**
   - If password security is breached, immediately change it
   ```bash
   zep reset-password
   ```

3. **Forgotten Password Recovery**
   - As long as current password is known, it can be changed
   - If current password is forgotten, vault cannot be accessed

### Perfect Password Guidelines

Choose a strong new password with:
- ✅ Minimum 12 characters (longer is better)
- ✅ Mix of uppercase and lowercase letters
- ✅ Mix of numbers and special characters
- ✅ No dictionary words or common phrases
- ✅ Not reused from other accounts
- ✅ Stored securely in password manager

### Troubleshooting

**Error: "Re-encryption failed"**
- Ensure you have stable internet connection
- Ensure sufficient disk space for temporary operations
- Retry the reset operation

**Error: "Index update failed"**
- Vault may have inconsistent state
- Retry the password reset completely
- If problem persists, contact support

**Password Reset Takes Too Long**
- Large vaults take longer to re-encrypt
- This is normal for vaults with thousands of files
- Do not interrupt the process
- Estimated time: ~1 minute per 100 files

### Related Commands

- [`zep login`](SETUP.md): Initial authentication
- [`zep transfer-vault`](TRANSFER.md): Migrate vault to different password
- [`zep settings`](SETTINGS.md): Manage vault configuration
- [`zep logout/disconnect`](AUTH.md): End session

### Password Reset API

For programmatic password resets:

```go
session, err := Connect("username", "oldPassword")
if err != nil {
    fmt.Println("Connection failed:", err)
    return
}

err = ResetPassword(session, "newPassword")
if err != nil {
    fmt.Println("Password reset failed:", err)
    return
}

// Session now uses new password for all operations
err = session.Save()
if err != nil {
    fmt.Println("Session save failed:", err)
    return
}

fmt.Println("✔ Password reset successful")
```

### Additional Notes

- Password reset only affects your vault, not your GitHub account
- GitHub authentication (SSH key) is not affected by password reset
- Multiple users can use same GitHub repository with different passwords
- Per-file encryption keys ensure privacy even during password transition
