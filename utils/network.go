package utils

import (
	"fmt"
	"io"
	"net/http"
	"time" // Add this
)

func FetchRaw(username, path string) ([]byte, error) {
	// Add a timestamp to the end of the URL to bypass GitHub's cache
	cacheBuster := time.Now().UnixNano()
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/.nexus/refs/heads/master/%s?t=%d",
		username, path, cacheBuster)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Nexus-CLI-v1")
	// Standard headers to tell proxies not to cache
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
