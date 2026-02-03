package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"nexus-cli/utils"
)

func main() {
	const (
		username   = "Auchrio"
		repoURL    = "git@github.com:Auchrio/.nexus.git"
		keyPath    = ".config/key"
		sourceFile = "test.txt"
		vaultPath  = "test.txt"
	)

	// 1. Credentials
	fmt.Print("Enter Vault Password: ")
	var password string
	fmt.Scanln(&password)

	rawKey, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatalf("Failed to read deploy key: %v", err)
	}

	// 2. Fetch and Check Index
	fmt.Println("Fetching current index...")
	var index utils.VaultIndex
	var realName string // We define this here to use it later

	rawIndex, err := utils.FetchRaw(username, ".config/index")
	if err != nil {
		fmt.Println("No existing index found, creating a new one.")
		index = utils.NewIndex()
	} else {
		index, err = utils.FromBytes(rawIndex, password)
		if err != nil {
			log.Fatalf("Failed to decrypt index: %v", err)
		}

		fmt.Println("--- Current Decoded Index ---")
		index.PrintDebug()
		fmt.Println("-----------------------------")

		// CONFLICT CHECK & ID REUSE
		if entry, exists := index[vaultPath]; exists {
			fmt.Printf("⚠ CONFLICT: '%s' already exists (points to %s).\n", vaultPath, entry.RealName)
			fmt.Print("Do you want to overwrite it? (y/N): ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response != "y" && response != "yes" {
				fmt.Println("Upload cancelled by user.")
				return
			}

			// SUCCESS: Reuse the existing 16-char filename
			realName = entry.RealName
			fmt.Printf("Proceeding with overwrite using existing ID: %s\n", realName)
		}
	}

	// 3. Prepare File
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		log.Fatalf("Failed to read source file: %v", err)
	}

	encryptedFile, err := utils.Encrypt(content, password)
	if err != nil {
		log.Fatalf("Encryption failed: %v", err)
	}

	// If realName is still empty (new file), generate a new one
	if realName == "" {
		realName = utils.GenerateRandomName()
		index.AddFile(vaultPath, realName)
	}

	encryptedIndex, err := index.ToBytes(password)
	if err != nil {
		log.Fatalf("Failed to encrypt index: %v", err)
	}

	// 4. Atomic Push
	filesToPush := map[string][]byte{
		realName:        encryptedFile,
		".config/index": encryptedIndex,
	}

	fmt.Printf("Pushing %s (as %s) and updated index...\n", vaultPath, realName)
	err = utils.PushFiles(repoURL, rawKey, filesToPush, "Nexus: Updated "+vaultPath)
	if err != nil {
		log.Fatalf("Git push failed: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("✔ SUCCESS: Vault updated.")
	fmt.Println("--------------------------------------------------")
}
