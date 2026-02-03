package utils

import (
	"fmt"
	"os"
)

func DownloadFile(vaultPath string, outputPath string, session *Session) error {
	// 1. Use your custom FindEntry logic to navigate the nested maps
	entry, err := session.Index.FindEntry(vaultPath)
	if err != nil {
		return fmt.Errorf("could not find file in vault: %w", err)
	}

	// 2. Safety check: Ensure we aren't trying to "download" a folder
	if entry.Type == "folder" {
		return fmt.Errorf("'%s' is a directory, you can only download individual files", vaultPath)
	}

	fmt.Printf("Downloading %s (Storage ID: %s)...\n", vaultPath, entry.RealName)

	// 3. Fetch the encrypted hex-named file from GitHub
	encryptedData, err := FetchRaw(session.Username, entry.RealName)
	if err != nil {
		return fmt.Errorf("failed to fetch storage file from remote: %w", err)
	}

	// 4. Decrypted with the session password
	decryptedData, err := Decrypt(encryptedData, session.Password)
	if err != nil {
		return fmt.Errorf("decryption failed: check your password")
	}

	// 5. Save to the local output path
	return os.WriteFile(outputPath, decryptedData, 0644)
}
