package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func ListFiles(session *Session, folderPath string) error {
	// Start with the root map
	currentMap := session.Index

	// If a path is provided, navigate to that entry
	if folderPath != "" && folderPath != "/" {
		entry, err := session.Index.FindEntry(folderPath)
		if err != nil {
			return err
		}
		if entry.Type != "folder" {
			return fmt.Errorf("'%s' is a file, not a folder", folderPath)
		}
		currentMap = entry.Contents
	}

	if len(currentMap) == 0 {
		fmt.Println("Directory is empty.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tTYPE\tSTORAGE ID")
	fmt.Fprintln(w, "----\t----\t----------")

	for name, entry := range currentMap {
		displayType := "[FILE]"
		rName := entry.RealName
		displayName := name

		if entry.Type == "folder" {
			displayType = "[DIR]"
			rName = "-"
			displayName = name + "/"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", displayName, displayType, rName)
	}

	return w.Flush()
}