package utils

import (
	"encoding/hex"
	"fmt"
	"os"
)

func UploadFile(sourcePath string, vaultPath string, session *Session) error {
	repoURL := fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username)

	// 1. Read source
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// 2. Determine Storage Name and File Key
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
		realName = GenerateRandomName()
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

	// 3. Encrypt file data with the per-file key
	encryptedData, err := EncryptWithKey(data, fileKey)
	if err != nil {
		return err
	}

	// 4. Encrypt updated index
	indexBytes, err := session.Index.ToBytes(session.Password)
	if err != nil {
		return err
	}

	// 5. Push to Git
	filesToPush := map[string][]byte{
		realName:        encryptedData,
		".config/index": indexBytes,
	}

	err = PushFiles(repoURL, session.RawKey, filesToPush, "Nexus: Updated Vault")
	if err != nil {
		return err
	}

	// 6. Save updated index to local session to bypass cache
	return nil
}
