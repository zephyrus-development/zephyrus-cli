package utils

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func UploadFile(sourcePath string, vaultPath string, session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username)

	// 1. Read source
	PrintProgressStep(1, 5, "Reading file...")
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	PrintCompletionLine("File read successfully")

	// 2. Determine Storage Name and File Key
	PrintProgressStep(2, 5, "Validating vault...")
	var realName string
	var fileKey []byte

	entry, err := session.Index.FindEntry(vaultPath)
	if err == nil && entry.Type == "file" {
		// Existing file: reuse storage name, decrypt existing key
		realName = entry.RealName
		fmt.Printf("Updating existing file: %s (%s)\n", vaultPath, realName)

		// Decrypt the existing file key
		encryptedKey, _ := hex.DecodeString(entry.FileKey)
		fileKey, err = Decrypt(encryptedKey, session.Password)
		if err != nil {
			return fmt.Errorf("failed to decrypt file key: %w", err)
		}
	} else {
		// New file: generate new storage name and file key
		// Use configurable hash length from settings
		hashByteLength := session.Settings.FileHashLength / 2 // Convert hex chars to bytes
		realName = GenerateRandomNameWithLength(hashByteLength)
		fileKey = GenerateFileKey()

		// Encrypt the file key with the vault password
		encryptedKey, err := Encrypt(fileKey, session.Password)
		if err != nil {
			return fmt.Errorf("failed to encrypt file key: %w", err)
		}
		encryptedKeyHex := hex.EncodeToString(encryptedKey)

		session.Index.AddFile(vaultPath, realName, encryptedKeyHex)
		fmt.Printf("Uploading new file: %s as %s\n", vaultPath, realName)
	}
	PrintCompletionLine("File validated")

	// 3. Encrypt file data with the per-file key
	PrintProgressStep(3, 5, "Encrypting file...")
	time.Sleep(time.Millisecond * 100) // Simulate work for visibility
	encryptedData, err := EncryptWithKey(data, fileKey)
	if err != nil {
		return err
	}
	PrintCompletionLine("File encrypted")

	// 4. Encrypt updated index
	PrintProgressStep(4, 5, "Updating vault index...")
	indexBytes, err := session.Index.ToBytes(session.Password)
	if err != nil {
		return err
	}
	PrintCompletionLine("Vault index updated")

	// 5. Push to Git
	PrintProgressStep(5, 5, "Uploading to GitHub...")
	filesToPush := map[string][]byte{
		realName:        encryptedData,
		".config/index": indexBytes,
	}

	err = PushFilesWithAuthor(repoURL, session.RawKey, filesToPush, session.Settings.CommitMessage, session.Settings.CommitAuthorName, session.Settings.CommitAuthorEmail)
	if err != nil {
		return err
	}
	PrintCompletionLine("Upload to GitHub completed")

	// 6. Save updated index to local session to bypass cache
	return nil
}

// UploadDirectory uploads an entire directory recursively to the vault
func UploadDirectory(sourceDirPath string, vaultPath string, session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username)

	// 1. Verify directory exists
	fileInfo, err := os.Stat(sourceDirPath)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("path is not a directory: %s", sourceDirPath)
	}

	filesToPush := make(map[string][]byte)
	fileCount := 0

	fmt.Printf("Scanning directory: %s\n", sourceDirPath)

	// 2. Walk through all files in the directory recursively
	err = filepath.Walk(sourceDirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, only process files
		if info.IsDir() {
			return nil
		}

		// Calculate relative path from source directory
		relPath, err := filepath.Rel(sourceDirPath, filePath)
		if err != nil {
			return err
		}

		// Construct vault path preserving directory structure
		currentVaultPath := vaultPath + "/" + relPath

		fmt.Printf("Processing file (%d): %s\n", fileCount+1, relPath)

		// 3. Read source file
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// 4. Determine Storage Name and File Key
		var realName string
		var fileKey []byte

		entry, err := session.Index.FindEntry(currentVaultPath)
		if err == nil && entry.Type == "file" {
			// Existing file: reuse storage name, decrypt existing key
			realName = entry.RealName
			fmt.Printf("  → Updating: %s (%s)\n", currentVaultPath, realName)

			// Decrypt the existing file key
			encryptedKey, _ := hex.DecodeString(entry.FileKey)
			fileKey, err = Decrypt(encryptedKey, session.Password)
			if err != nil {
				return fmt.Errorf("failed to decrypt file key for %s: %w", currentVaultPath, err)
			}
		} else {
			// New file: generate new storage name and file key
			hashByteLength := session.Settings.FileHashLength / 2
			realName = GenerateRandomNameWithLength(hashByteLength)
			fileKey = GenerateFileKey()

			// Encrypt the file key with the vault password
			encryptedKey, err := Encrypt(fileKey, session.Password)
			if err != nil {
				return fmt.Errorf("failed to encrypt file key for %s: %w", currentVaultPath, err)
			}
			encryptedKeyHex := hex.EncodeToString(encryptedKey)

			session.Index.AddFile(currentVaultPath, realName, encryptedKeyHex)
			fmt.Printf("  → New file: %s as %s\n", currentVaultPath, realName)
		}

		// 5. Encrypt file data with the per-file key
		encryptedData, err := EncryptWithKey(data, fileKey)
		if err != nil {
			return err
		}

		// Collect encrypted file for batch push
		filesToPush[realName] = encryptedData
		fileCount++

		return nil
	})

	if err != nil {
		return fmt.Errorf("directory walk failed: %w", err)
	}

	if fileCount == 0 {
		return fmt.Errorf("no files found in directory: %s", sourceDirPath)
	}

	fmt.Printf("\nUploading %d files to vault...\n", fileCount)

	// 6. Encrypt updated index
	PrintProgressStep(1, 2, "Encrypting vault index...")
	indexBytes, err := session.Index.ToBytes(session.Password)
	if err != nil {
		return err
	}
	PrintCompletionLine("Vault index updated")

	// 7. Add index to push
	filesToPush[".config/index"] = indexBytes

	// 8. Push all files to Git in a single operation
	PrintProgressStep(2, 2, "Uploading to GitHub...")
	err = PushFilesWithAuthor(repoURL, session.RawKey, filesToPush, session.Settings.CommitMessage, session.Settings.CommitAuthorName, session.Settings.CommitAuthorEmail)
	if err != nil {
		return err
	}
	PrintCompletionLine("Upload to GitHub completed")

	fmt.Printf("✔ Successfully uploaded %d files from directory\n", fileCount)
	return nil
}
