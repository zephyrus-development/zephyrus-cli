package main

import (
	"bytes"
	"fmt"
	"log"
	"nexus-cli/utils"
)

func main() {
	// 1. Setup test data
	password := "super-secret-vault-password"
	originalData := []byte("This is a secret note for the Nexus Vault.")

	fmt.Println("--- Nexus CLI Encryption Test ---")
	fmt.Printf("Original Plaintext: %s\n", string(originalData))
	fmt.Printf("Password Used:      %s\n", password)

	// 2. Test Encryption
	fmt.Println("\nEncrypting...")
	encryptedBlob, err := utils.Encrypt(originalData, password)
	if err != nil {
		log.Fatalf("Encryption failed: %v", err)
	}

	// Show the packaged data [Salt(16) + Nonce(12) + Ciphertext]
	fmt.Printf("Encrypted Blob (hex): %x...\n", encryptedBlob[:32])
	fmt.Printf("Total Blob Length:    %d bytes\n", len(encryptedBlob))

	// 3. Test Decryption
	fmt.Println("\nDecrypting...")
	decryptedData, err := utils.Decrypt(encryptedBlob, password)
	if err != nil {
		log.Fatalf("Decryption failed: %v", err)
	}

	fmt.Printf("Decrypted Result:   %s\n", string(decryptedData))

	// 4. Verification
	if bytes.Equal(originalData, decryptedData) {
		fmt.Println("\n✅ SUCCESS: Decrypted data matches original.")
	} else {
		fmt.Println("\n❌ FAILURE: Decrypted data does not match.")
	}

	// 5. Test Failure (Wrong Password)
	fmt.Println("\nTesting with incorrect password...")
	_, err = utils.Decrypt(encryptedBlob, "wrong-password")
	if err != nil {
		fmt.Printf("Expected Error caught: %v\n", err)
		fmt.Println("✅ SUCCESS: Decryption blocked with wrong password.")
	}
}