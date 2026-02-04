package utils

import (
	"os"
	"os/exec"
	"runtime"
)

// LocalLS runs the 'ls' command from the main terminal
func LocalLS(args []string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// On Windows, try to use ls if available (Git Bash, WSL, etc.)
		// Fall back to dir if ls doesn't exist
		cmd = exec.Command("ls", args...)
		// Set stdio to current
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			// If ls fails, fall back to dir
			return LocalDir(args)
		}
		return nil
	}

	// On Unix/Linux, use ls
	cmd = exec.Command("ls", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// LocalDir runs the 'dir' command on Windows or equivalent on other platforms
func LocalDir(args []string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		// On Windows, use dir command
		// Need to prepend dir to args
		dirArgs := append([]string{}, args...)
		cmd = exec.Command("cmd", append([]string{"/c", "dir"}, dirArgs...)...)
	} else {
		// On Unix/Linux, use ls with more detail (equivalent to dir)
		lsArgs := append([]string{"-la"}, args...)
		cmd = exec.Command("ls", lsArgs...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
