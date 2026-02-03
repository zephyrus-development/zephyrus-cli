package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Entry represents either a file or a folder in the Nexus vault
type Entry struct {
	Type     string           `json:"type"`
	RealName string           `json:"realName,omitempty"`
	Contents map[string]Entry `json:"contents,omitempty"`
}

// VaultIndex is the top-level structure for .config/index
type VaultIndex map[string]Entry

// NewIndex creates an empty vault index
func NewIndex() VaultIndex {
	return make(VaultIndex)
}

// FromBytes decrypts and parses the JSON index
func FromBytes(data []byte, password string) (VaultIndex, error) {
	decrypted, err := Decrypt(data, password)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt index: %w", err)
	}

	var index VaultIndex
	if err := json.Unmarshal(decrypted, &index); err != nil {
		return nil, fmt.Errorf("failed to parse index JSON: %w", err)
	}
	return index, nil
}

// ToBytes serializes and encrypts the index
func (vi VaultIndex) ToBytes(password string) ([]byte, error) {
	plaintext, err := json.MarshalIndent(vi, "", "  ")
	if err != nil {
		return nil, err
	}
	return Encrypt(plaintext, password)
}

// FindEntry navigates the index based on a path (e.g., "images/vacation.png")
func (vi VaultIndex) FindEntry(path string) (*Entry, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentMap := vi

	for i, part := range parts {
		entry, exists := currentMap[part]
		if !exists {
			return nil, fmt.Errorf("path component '%s' not found", part)
		}

		if i == len(parts)-1 {
			return &entry, nil
		}

		if entry.Type != "folder" {
			return nil, fmt.Errorf("'%s' is a file, not a folder", part)
		}
		currentMap = entry.Contents
	}
	return nil, fmt.Errorf("invalid path")
}

// AddFile inserts a new file into the index
func (vi VaultIndex) AddFile(path string, realName string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	currentMap := vi

	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		entry, exists := currentMap[part]
		if !exists {
			entry = Entry{
				Type:     "folder",
				Contents: make(map[string]Entry),
			}
			currentMap[part] = entry
		}
		currentMap = entry.Contents
	}

	fileName := parts[len(parts)-1]
	currentMap[fileName] = Entry{
		Type:     "file",
		RealName: realName,
	}
}

// PrintDebug prints the index structure to the console
func (vi VaultIndex) PrintDebug() {
	if len(vi) == 0 {
		fmt.Println("Index is currently empty.")
		return
	}
	data, _ := json.MarshalIndent(vi, "", "  ")
	fmt.Println(string(data))
}
