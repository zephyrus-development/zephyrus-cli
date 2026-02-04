package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"path/filepath"
	"time"
)

// GenerateShareReference generates a base62 reference with default 6 characters
func GenerateShareReference() (string, error) {
	return GenerateShareReferenceWithLength(6)
}

// GenerateShareReferenceWithLength generates a base62 reference with configurable length
func GenerateShareReferenceWithLength(length int) (string, error) {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ref := make([]byte, length)
	for i := 0; i < length; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		ref[i] = charset[idx.Int64()]
	}
	return string(ref), nil
}

// ShareFile generates a share string with a new 6-char reference and share password
// File is encrypted with share password and uploaded to /shared/{ref}
func ShareFile(vaultPath string, sharePassword string, session *Session) (string, error) {
	// 1. Find the file entry in the index
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return "", fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Ensure it's a file, not a folder
	if entry.Type == "folder" {
		return "", fmt.Errorf("'%s' is a directory, you can only share individual files", vaultPath)
	}

	// 3. Fetch the original encrypted file from GitHub
	encryptedData, err := FetchRaw(session.Username, entry.RealName)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from remote: %w", err)
	}

	// 4. Decrypt the file with the vault's file key
	encryptedKey, err := hex.DecodeString(entry.FileKey)
	if err != nil {
		return "", fmt.Errorf("invalid file key in index: %w", err)
	}
	fileKey, err := Decrypt(encryptedKey, session.Password)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt file key: check your password")
	}
	decryptedData, err := DecryptWithKey(encryptedData, fileKey)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	// 5. Generate a new reference with configurable length from settings
	ref, err := GenerateShareReferenceWithLength(session.Settings.ShareHashLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate share reference: %w", err)
	}

	// 6. Encrypt the file with the share password
	shareEncrypted, err := Encrypt(decryptedData, sharePassword)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt with share password: %w", err)
	}

	// 7. Upload to /shared/{ref}
	sharedPath := fmt.Sprintf("shared/%s", ref)
	filesToPush := map[string][]byte{
		sharedPath: shareEncrypted,
	}

	err = PushFilesWithAuthor(
		fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username),
		session.RawKey,
		filesToPush,
		session.Settings.CommitMessage,
		session.Settings.CommitAuthorName,
		session.Settings.CommitAuthorEmail,
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload shared file: %w", err)
	}

	// 8. Add entry to shared index
	if session.SharedIndex == nil {
		session.SharedIndex = NewSharedIndex()
	}
	indexEntry := SharedFileEntry{
		Name:         vaultPath,
		Reference:    ref,
		Password:     sharePassword,
		SharedAt:     time.Now(),
		OriginalPath: vaultPath,
	}
	session.SharedIndex.AddEntry(indexEntry)

	// 9. Upload the updated shared index
	indexJSON, err := session.SharedIndex.EncryptForRemote(session.Password)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt shared index: %w", err)
	}

	indexFilesToPush := map[string][]byte{
		"shared/.config/index": indexJSON,
	}

	err = PushFilesWithAuthor(
		fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username),
		session.RawKey,
		indexFilesToPush,
		session.Settings.CommitMessage,
		session.Settings.CommitAuthorName,
		session.Settings.CommitAuthorEmail,
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload shared index: %w", err)
	}

	// 10. Generate the share string: username:reference:sharepassword:base64filename
	filename := filepath.Base(vaultPath)
	encodedFilename := base64.StdEncoding.EncodeToString([]byte(filename))
	shareString := fmt.Sprintf("%s:%s:%s:%s", session.Username, ref, sharePassword, encodedFilename)

	return shareString, nil
}
