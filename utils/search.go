package utils

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func SearchFiles(session *Session, query string) error {
	fmt.Printf("Searching vault for: \"%s\"\n", query)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "VAULT PATH\tTYPE\tSTORAGE ID")
	fmt.Fprintln(w, "----------\t----\t----------")

	lowerQuery := strings.ToLower(query)
	found := false

	// Recursive helper to walk the Entry tree
	var walk func(map[string]Entry, string)
	walk = func(entries map[string]Entry, currentPath string) {
		for name, entry := range entries {
			fullPath := name
			if currentPath != "" {
				fullPath = currentPath + "/" + name
			}

			// Check if the current name or full path matches the query
			if strings.Contains(strings.ToLower(fullPath), lowerQuery) {
				found = true
				displayType := "[FILE]"
				rName := entry.RealName
				if entry.Type == "folder" {
					displayType = "[DIR]"
					rName = "-"
				}
				fmt.Fprintf(w, "%s\t%s\t%s\n", fullPath, displayType, rName)
			}

			// If it's a folder, dive deeper
			if entry.Type == "folder" && entry.Contents != nil {
				walk(entry.Contents, fullPath)
			}
		}
	}

	walk(session.Index, "")
	w.Flush()

	if !found {
		fmt.Println("No matches found.")
	}
	return nil
}
