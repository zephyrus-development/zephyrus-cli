package main

import (
	"fmt"
	"log"
	"nexus-cli/utils"
	"os"

	"github.com/spf13/cobra"
)

var (
	username string
	keyPath  string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "nexus",
		Short: "Nexus CLI: A stateless, encrypted git-based vault",
		Long:  `Nexus is a secure, git-backed vault that uses client-side encryption to store your files privately on GitHub.`,
	}

	// --- SETUP ---
	var setupCmd = &cobra.Command{
		Use:   "setup [username] [key-path]",
		Short: "Initialize the vault and encrypt your master key",
		Example: `  nexus setup myuser ~/.ssh/id_ed25519
  nexus setup -u myuser -k /path/to/key`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				username = args[0]
			}
			if len(args) > 1 {
				keyPath = args[1]
			}

			if username == "" || keyPath == "" {
				fmt.Println("Error: Username and Key Path are required via arguments or flags.")
				cmd.Help()
				return
			}

			pass, _ := utils.GetPassword("Create a Vault Password: ")
			err := utils.SetupVault(username, keyPath, pass)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("✔ Setup complete.")
		},
	}
	setupCmd.Flags().StringVarP(&username, "user", "u", "", "GitHub username")
	setupCmd.Flags().StringVarP(&keyPath, "key", "k", "", "Path to local private key")

	// --- CONNECT ---
	var connectCmd = &cobra.Command{
		Use:     "connect [username]",
		Short:   "Login and sync the remote index to your local session",
		Example: "  nexus connect myuser",
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			targetUser := ""
			if len(args) > 0 {
				targetUser = args[0]
			} else {
				fmt.Print("Enter GitHub Username: ")
				fmt.Scanln(&targetUser)
			}
			pass, _ := utils.GetPassword("Enter Vault Password: ")
			err := utils.Connect(targetUser, pass)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	// --- UPLOAD ---
	var uploadCmd = &cobra.Command{
		Use:     "upload [local-path] [vault-path]",
		Short:   "Upload a file to the vault",
		Example: "  nexus upload ./secrets.txt work/finance/secrets.txt",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			err = utils.UploadFile(args[0], args[1], session)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("✔ Upload successful.")
		},
	}

	// --- DOWNLOAD ---
	var downloadCmd = &cobra.Command{
		Use:     "download [vault-path] [local-output-path]",
		Short:   "Download a file from the vault",
		Example: "  nexus download work/finance/secrets.txt ./recovered.txt",
		Args:    cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			err = utils.DownloadFile(args[0], args[1], session)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("✔ Download successful.")
		},
	}

	// --- DELETE ---
	var deleteCmd = &cobra.Command{
		Use:   "delete [vault-path]",
		Short: "Remove a file or folder (recursively) from the vault",
		Example: `  nexus delete work/old_file.txt
  nexus delete work/deprecated_folder`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			err = utils.DeletePath(args[0], session)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("✔ Successfully removed from vault.")
		},
	}

	// --- LIST ---
	var listCmd = &cobra.Command{
		Use:     "ls [folder]",
		Aliases: []string{"list"},
		Short:   "List files and folders at a specific path",
		Example: `  nexus ls
  nexus ls work/finance`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			path := ""
			if len(args) > 0 {
				path = args[0]
			}
			err = utils.ListFiles(session, path)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	// --- SEARCH ---
	var searchCmd = &cobra.Command{
		Use:     "search [query]",
		Short:   "Search for files/folders in the vault index",
		Example: "  nexus search .pdf",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			err = utils.SearchFiles(session, args[0])
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	// --- PURGE / DISCONNECT ---
	var purgeCmd = &cobra.Command{
		Use:   "purge",
		Short: "⚠️  DANGER: Hard-reset the remote repository (Wipe all data)",
		Run: func(cmd *cobra.Command, args []string) {
			session, err := utils.GetSession()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print("⚠️  Confirm PURGE? (y/N): ")
			var confirm string
			fmt.Scanln(&confirm)
			if confirm != "y" {
				return
			}
			err = utils.PurgeVault(session)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("✔ Vault purged.")
		},
	}

	var disconnectCmd = &cobra.Command{
		Use:   "disconnect",
		Short: "Clear the local session and cached index",
		Run: func(cmd *cobra.Command, args []string) {
			utils.Disconnect()
			fmt.Println("✔ Disconnected.")
		},
	}

	rootCmd.AddCommand(setupCmd, connectCmd, disconnectCmd, uploadCmd, downloadCmd, deleteCmd, listCmd, searchCmd, purgeCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
