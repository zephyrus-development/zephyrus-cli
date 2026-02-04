package utils

import (
	"fmt"
	"io"
	"net/http"
	"time" // Add this
)

func FetchRaw(username, path string) ([]byte, error) {
	// Use the most direct raw URL format
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/.zephyrus/master/%s?t=%d",
		username, path, time.Now().UnixNano())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("404") // Explicitly return 404 string for main.go to check
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
