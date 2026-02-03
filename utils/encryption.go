package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	SaltSize   = 16
	Iterations = 100000
	KeySize    = 32 // AES-256
	NonceSize  = 12 // Standard for GCM
)

// Encrypt handles key derivation via PBKDF2 and encryption via AES-GCM
func Encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	// Derive key using PBKDF2
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Bundle as: [Salt][Nonce][Ciphertext]
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt extracts the salt, re-derives the key, and decrypts the data
func Decrypt(data []byte, password string) ([]byte, error) {
	if len(data) < SaltSize+NonceSize {
		return nil, errors.New("ciphertext too short")
	}

	// Extract metadata
	salt := data[:SaltSize]
	nonce := data[SaltSize : SaltSize+NonceSize]
	ciphertext := data[SaltSize+NonceSize:]

	// Re-derive key
	key := pbkdf2.Key([]byte(password), salt, Iterations, KeySize, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nonce, ciphertext, nil)
}
