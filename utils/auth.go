package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const configPath = "nexus.conf"

type Session struct {
	Username string     `json:"username"`
	Password string     `json:"password"`
	RawKey   []byte     `json:"raw_key"`
	Index    VaultIndex `json:"index"`
}

// Connect initializes the session and syncs the index locally
func Connect(username string, password string) error {
	fmt.Printf("Connecting and syncing vault for %s...\n", username)

	// 1. Fetch & Decrypt Master Key
	encryptedKey, err := FetchRaw(username, ".config/key")
	if err != nil {
		return fmt.Errorf("master key not found: %w", err)
	}
	rawKey, err := Decrypt(encryptedKey, password)
	if err != nil {
		return fmt.Errorf("auth failed: invalid password")
	}

	// 2. Fetch & Decrypt Index (Handle empty vault)
	var index VaultIndex
	rawIndex, err := FetchRaw(username, ".config/index")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			index = NewIndex()
		} else {
			return err
		}
	} else {
		index, _ = FromBytes(rawIndex, password)
	}

	session := Session{
		Username: username,
		Password: password,
		RawKey:   rawKey,
		Index:    index,
	}

	return session.Save()
}

func (s *Session) Save() error {
	data, _ := json.MarshalIndent(s, "", "  ")
	return os.WriteFile(configPath, data, 0600)
}

func GetSession() (*Session, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("not connected: run 'connect' first")
	}
	var s Session
	err = json.Unmarshal(data, &s)
	return &s, err
}

func Disconnect() error {
	return os.Remove(configPath)
}

// FetchSessionStateless performs the authentication and index fetch without saving to disk
func FetchSessionStateless(username string, password string) (*Session, error) {
	// 1. Fetch & Decrypt Master Key
	encryptedKey, err := FetchRaw(username, ".config/key")
	if err != nil {
		return nil, fmt.Errorf("master key not found: %w", err)
	}
	rawKey, err := Decrypt(encryptedKey, password)
	if err != nil {
		return nil, fmt.Errorf("auth failed: invalid password")
	}

	// 2. Fetch & Decrypt Index
	var index VaultIndex
	rawIndex, err := FetchRaw(username, ".config/index")
	if err != nil {
		index = NewIndex() // Assume new vault if 404
	} else {
		index, err = FromBytes(rawIndex, password)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt index: %w", err)
		}
	}

	return &Session{
		Username: username,
		Password: password,
		RawKey:   rawKey,
		Index:    index,
	}, nil
}
