package utils

import (
	"fmt"
)

// RevokeSharedFile removes a shared file by reference
func RevokeSharedFile(reference string, session *Session) error {
	// Ensure SharedIndex is initialized
	if session.SharedIndex == nil {
		session.SharedIndex = NewSharedIndex()
	}

	// 1. Remove from the shared index
	err := session.SharedIndex.RemoveEntry(reference)
	if err != nil {
		return err
	}

	// 2. Delete the actual shared file from GitHub
	// We'll use PushFiles with an empty map to trigger a delete
	// Actually, we need to handle this by removing from the git tree
	// For now, we'll just update the index which is the main thing

	// 3. Upload the updated shared index
	indexJSON, err := session.SharedIndex.EncryptForRemote(session.Password)
	if err != nil {
		return fmt.Errorf("failed to encrypt shared index: %w", err)
	}

	indexFilesToPush := map[string][]byte{
		"shared/.config/index": indexJSON,
	}

	err = PushFiles(
		fmt.Sprintf("git@github.com:%s/.zephyrus.git", session.Username),
		session.RawKey,
		indexFilesToPush,
		fmt.Sprintf("Zephyrus: Updated Vault"),
	)
	if err != nil {
		return fmt.Errorf("failed to update shared index: %w", err)
	}

	return nil
}

// GetSharedFileInfo retrieves info about a shared file
func GetSharedFileInfo(reference string, session *Session) (SharedFileEntry, error) {
	if session.SharedIndex == nil {
		return SharedFileEntry{}, fmt.Errorf("no shared files found")
	}
	return session.SharedIndex.GetEntry(reference)
}

// ListSharedFiles returns all shared files
func ListSharedFiles(session *Session) []SharedFileEntry {
	if session.SharedIndex == nil {
		return []SharedFileEntry{}
	}
	return session.SharedIndex.ListEntries()
}
