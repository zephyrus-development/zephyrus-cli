package utils

import (
	"encoding/json"
	"fmt"
	"time"
)

// SharedFileEntry represents a shared file in the index
type SharedFileEntry struct {
	Name         string    `json:"name"`
	Reference    string    `json:"reference"`
	Password     string    `json:"password"`
	SharedAt     time.Time `json:"shared_at"`
	OriginalPath string    `json:"original_path"`
}

// SharedIndex stores all shared files with encryption
type SharedIndex struct {
	Files map[string]SharedFileEntry `json:"files"` // Key is reference ID
}

// NewSharedIndex creates an empty shared index
func NewSharedIndex() *SharedIndex {
	return &SharedIndex{
		Files: make(map[string]SharedFileEntry),
	}
}

// AddEntry adds a new shared file to the index
func (si *SharedIndex) AddEntry(entry SharedFileEntry) {
	si.Files[entry.Reference] = entry
}

// GetEntry retrieves a shared file by reference
func (si *SharedIndex) GetEntry(reference string) (SharedFileEntry, error) {
	entry, exists := si.Files[reference]
	if !exists {
		return SharedFileEntry{}, fmt.Errorf("shared file with reference '%s' not found", reference)
	}
	return entry, nil
}

// RemoveEntry deletes a shared file from the index
func (si *SharedIndex) RemoveEntry(reference string) error {
	if _, exists := si.Files[reference]; !exists {
		return fmt.Errorf("shared file with reference '%s' not found", reference)
	}
	delete(si.Files, reference)
	return nil
}

// ListEntries returns all shared files
func (si *SharedIndex) ListEntries() []SharedFileEntry {
	entries := make([]SharedFileEntry, 0, len(si.Files))
	for _, entry := range si.Files {
		entries = append(entries, entry)
	}
	return entries
}

// ToJSON encodes the index as JSON
func (si *SharedIndex) ToJSON() ([]byte, error) {
	return json.MarshalIndent(si, "", "  ")
}

// FromJSON decodes the index from JSON
func (si *SharedIndex) FromJSON(data []byte) error {
	return json.Unmarshal(data, si)
}

// EncryptForRemote encrypts the index for storage on GitHub
func (si *SharedIndex) EncryptForRemote(password string) ([]byte, error) {
	jsonData, err := si.ToJSON()
	if err != nil {
		return nil, err
	}
	return Encrypt(jsonData, password)
}

// DecryptFromRemote decrypts the index from GitHub storage
func DecryptSharedIndex(encryptedData []byte, password string) (*SharedIndex, error) {
	jsonData, err := Decrypt(encryptedData, password)
	if err != nil {
		return nil, err
	}

	si := NewSharedIndex()
	err = si.FromJSON(jsonData)
	if err != nil {
		return nil, err
	}

	return si, nil
}
