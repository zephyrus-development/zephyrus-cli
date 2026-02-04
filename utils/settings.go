package utils

import (
	"encoding/json"
	"fmt"
)

// VaultSettings stores customizable runtime attributes for the vault
type VaultSettings struct {
	CommitAuthorName  string `json:"commit_author_name"`
	CommitAuthorEmail string `json:"commit_author_email"`
	CommitMessage     string `json:"commit_message"`
	FileHashLength    int    `json:"file_hash_length"`
	ShareHashLength   int    `json:"share_hash_length"`
}

// DefaultSettings returns the default vault settings
func DefaultSettings() VaultSettings {
	return VaultSettings{
		CommitAuthorName:  "Zephyrus",
		CommitAuthorEmail: "auchrio@proton.me",
		CommitMessage:     "Zephyrus: Updated Vault",
		FileHashLength:    16,
		ShareHashLength:   6,
	}
}

// FromBytes decrypts and parses the JSON settings
func SettingsFromBytes(data []byte, password string) (VaultSettings, error) {
	decrypted, err := Decrypt(data, password)
	if err != nil {
		return VaultSettings{}, fmt.Errorf("failed to decrypt settings: %w", err)
	}

	var settings VaultSettings
	if err := json.Unmarshal(decrypted, &settings); err != nil {
		return VaultSettings{}, fmt.Errorf("failed to parse settings JSON: %w", err)
	}

	// Validate and apply defaults for missing values
	if settings.CommitAuthorName == "" {
		settings.CommitAuthorName = "Zephyrus"
	}
	if settings.CommitAuthorEmail == "" {
		settings.CommitAuthorEmail = "auchrio@proton.me"
	}
	if settings.CommitMessage == "" {
		settings.CommitMessage = "Zephyrus: Updated Vault"
	}
	if settings.FileHashLength <= 0 {
		settings.FileHashLength = 16
	}
	if settings.ShareHashLength <= 0 {
		settings.ShareHashLength = 6
	}

	return settings, nil
}

// ToBytes serializes and encrypts the settings
func (s VaultSettings) ToBytes(password string) ([]byte, error) {
	plaintext, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return nil, err
	}
	return Encrypt(plaintext, password)
}

// Validate checks that settings have reasonable values
func (s VaultSettings) Validate() error {
	if s.CommitAuthorName == "" {
		return fmt.Errorf("commit author name cannot be empty")
	}
	if s.CommitAuthorEmail == "" {
		return fmt.Errorf("commit author email cannot be empty")
	}
	if s.CommitMessage == "" {
		return fmt.Errorf("commit message cannot be empty")
	}
	if s.FileHashLength < 8 || s.FileHashLength > 64 {
		return fmt.Errorf("file hash length must be between 8 and 64 (got %d)", s.FileHashLength)
	}
	if s.ShareHashLength < 4 || s.ShareHashLength > 32 {
		return fmt.Errorf("share hash length must be between 4 and 32 (got %d)", s.ShareHashLength)
	}
	return nil
}

// SaveSettings encrypts and pushes settings to .config/settings on remote
func SaveSettings(username string, password string, rawKey []byte, settings VaultSettings) error {
	if err := settings.Validate(); err != nil {
		return err
	}

	settingsBytes, err := settings.ToBytes(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt settings: %w", err)
	}

	filesToPush := map[string][]byte{
		".config/settings": settingsBytes,
	}

	repoURL := fmt.Sprintf("git@github.com:%s/.zephyrus.git", username)
	return PushFilesWithAuthor(repoURL, rawKey, filesToPush, settings.CommitMessage, settings.CommitAuthorName, settings.CommitAuthorEmail)
}
