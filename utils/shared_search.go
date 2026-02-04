package utils

import (
	"fmt"
	"strings"
)

// SharedFileMatch represents a match result when searching for shared files
type SharedFileMatch struct {
	Reference    string
	OriginalPath string
	FileName     string
	MatchScore   int // Lower score = better match
}

// FindSharedFilesByName searches for shared files matching a name query
// Uses fuzzy matching to allow shortnames and partial matches
func FindSharedFilesByName(nameQuery string, session *Session) ([]SharedFileMatch, error) {
	if session.SharedIndex == nil {
		session.SharedIndex = NewSharedIndex()
	}

	var matches []SharedFileMatch
	nameQueryLower := strings.ToLower(nameQuery)
	entries := session.SharedIndex.ListEntries()

	for _, entry := range entries {
		// Extract filename from the original path
		parts := strings.Split(entry.OriginalPath, "/")
		fileName := parts[len(parts)-1]
		fileNameLower := strings.ToLower(fileName)

		// Check for exact match
		if fileNameLower == nameQueryLower {
			matches = append(matches, SharedFileMatch{
				Reference:    entry.Reference,
				OriginalPath: entry.OriginalPath,
				FileName:     fileName,
				MatchScore:   0, // Exact match = best
			})
			continue
		}

		// Check for prefix match
		if strings.HasPrefix(fileNameLower, nameQueryLower) {
			matchScore := len(fileName) - len(nameQuery) // Shorter matches are better
			matches = append(matches, SharedFileMatch{
				Reference:    entry.Reference,
				OriginalPath: entry.OriginalPath,
				FileName:     fileName,
				MatchScore:   matchScore,
			})
			continue
		}

		// Check for substring match
		if strings.Contains(fileNameLower, nameQueryLower) {
			matchScore := 100 + (len(fileName) - len(nameQuery)) // Substring match = worse than prefix
			matches = append(matches, SharedFileMatch{
				Reference:    entry.Reference,
				OriginalPath: entry.OriginalPath,
				FileName:     fileName,
				MatchScore:   matchScore,
			})
		}
	}

	return matches, nil
}

// RevokeSharedFileByName revokes a shared file using its name (with search if ambiguous)
func RevokeSharedFileByName(nameQuery string, session *Session) (string, error) {
	matches, err := FindSharedFilesByName(nameQuery, session)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no shared files found matching '%s'", nameQuery)
	}

	if len(matches) > 1 {
		// Multiple matches - show all and ask user to be more specific
		fmt.Printf("Multiple files match '%s':\n", nameQuery)
		for i, match := range matches {
			fmt.Printf("  %d. %s (ref: %s)\n", i+1, match.FileName, match.Reference)
		}
		return "", fmt.Errorf("ambiguous file name - please be more specific")
	}

	// Exactly one match - revoke it
	reference := matches[0].Reference
	err = RevokeSharedFile(reference, session)
	if err != nil {
		return "", err
	}

	return reference, nil
}

// PrintSharedFilesFormatted lists all shared files with formatted output
func PrintSharedFilesFormatted(session *Session) error {
	entries := session.SharedIndex.ListEntries()
	if len(entries) == 0 {
		fmt.Println("No files have been shared yet.")
		return nil
	}

	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║            SHARED FILES                                         ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	for i, entry := range entries {
		// Extract filename from path
		parts := strings.Split(entry.OriginalPath, "/")
		fileName := parts[len(parts)-1]

		fmt.Printf("\n[%d] %s\n", i+1, fileName)
		fmt.Printf("    Vault Path:     %s\n", entry.OriginalPath)
		fmt.Printf("    Share Ref:      %s\n", entry.Reference)
		fmt.Printf("    Shared At:      %s\n", entry.SharedAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
	return nil
}

// GetSharedFileByName retrieves a shared file by name query
func GetSharedFileByName(nameQuery string, session *Session) (*SharedFileEntry, error) {
	matches, err := FindSharedFilesByName(nameQuery, session)
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no shared files found matching '%s'", nameQuery)
	}

	if len(matches) > 1 {
		fmt.Printf("Multiple files match '%s':\n", nameQuery)
		for i, match := range matches {
			fmt.Printf("  %d. %s (ref: %s)\n", i+1, match.FileName, match.Reference)
		}
		return nil, fmt.Errorf("ambiguous file name - please be more specific")
	}

	// Find the entry in the SharedIndex
	entries := session.SharedIndex.ListEntries()
	for i := range entries {
		if entries[i].Reference == matches[0].Reference {
			return &entries[i], nil
		}
	}

	return nil, fmt.Errorf("shared file entry not found")
}
