package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Session stores the decrypted credentials for the current local environment
type Session struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RawKey   []byte `json:"raw_key"`
}

const configPath = "nexus.conf"

// Connect downloads the master key, decrypts it, and saves a local session
func Connect(username string, password string) error {
	fmt.Printf("Connecting to vault for %s...\n", username)

	// 1. Fetch the encrypted key from GitHub using your existing helper
	encryptedKey, err := FetchRaw(username, ".config/key")
	if err != nil {
		return fmt.Errorf("could not find master key in vault: %w", err)
	}

	// 2. Decrypt it to verify the password and get the usable SSH key
	rawKey, err := Decrypt(encryptedKey, password)
	if err != nil {
		return fmt.Errorf("authentication failed: invalid password")
	}

	// 3. Prepare the session object
	session := Session{
		Username: username,
		Password: password,
		RawKey:   rawKey,
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	// 4. Write to local file with restricted permissions (Read/Write for owner only)
	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to save session: %w", err)
	}

	fmt.Println("✔ Connected. Session saved to nexus.conf")
	return nil
}

// Disconnect wipes the local session file
func Disconnect() error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("No active session found.")
		return nil
	}

	err := os.Remove(configPath)
	if err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}

	fmt.Println("✔ Disconnected. Local session cleared.")
	return nil
}

// GetSession is a helper for your other functions to load the active config
func GetSession() (*Session, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("not connected: please run 'connect' first")
	}

	var session Session
	err = json.Unmarshal(data, &session)
	return &session, err
}
