package utils

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	cryptossh "golang.org/x/crypto/ssh" // Needed for the HostKeyCallback
)

func PushFile(repoURL string, rawPrivateKey []byte, filename string, content []byte, commitMsg string) error {
	// 1. Setup SSH Auth
	publicKeys, err := ssh.NewPublicKeys("git", rawPrivateKey, "")
	if err != nil {
		return fmt.Errorf("auth error: %w", err)
	}

	// FIX: Use the x/crypto/ssh version of the callback
	publicKeys.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()

	// 2. Initialize in-memory repo
	// FIX: Init expects (Storer, Filesystem). memory.NewStorage is a Storer, memfs.New is a Filesystem.
	fs := memfs.New()
	storer := memory.NewStorage()

	r, err := git.Init(storer, fs)
	if err != nil {
		return err
	}

	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// 3. Create the file in the virtual filesystem
	file, err := fs.Create(filename)
	if err != nil {
		return err
	}
	file.Write(content)
	file.Close()

	// 4. Add and Commit
	_, err = w.Add(filename)
	if err != nil {
		return err
	}

	// FIX: Commit takes a pointer to git.CommitOptions, which contains the Signature
	_, err = w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Nexus CLI",
			Email: "nexus@cli.io",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// 5. Push
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKeys,
		Force:      true,
	})
}

func PushFiles(repoURL string, rawPrivateKey []byte, files map[string][]byte, commitMsg string) error {
	publicKeys, _ := ssh.NewPublicKeys("git", rawPrivateKey, "")
	publicKeys.HostKeyCallback = cryptossh.InsecureIgnoreHostKey()

	fs := memfs.New()
	storer := memory.NewStorage()
	r, _ := git.Init(storer, fs)
	w, _ := r.Worktree()

	// Add all files to the virtual filesystem
	for path, content := range files {
		// Ensure subdirectories exist for files like .config/index
		dir := filepath.Dir(path)
		if dir != "." {
			fs.MkdirAll(dir, 0755)
		}

		file, _ := fs.Create(path)
		file.Write(content)
		file.Close()
		w.Add(path)
	}

	w.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{Name: "Nexus CLI", Email: "nexus@cli.io", When: time.Now()},
	})

	_, err := r.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{repoURL}})
	if err != nil {
		return err
	}

	return r.Push(&git.PushOptions{RemoteName: "origin", Auth: publicKeys})
}
