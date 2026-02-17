package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync state to configured backend",
		Long:  `Push current state to remote repository (GitHub/GitLab).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSync()
		},
	}
}

func runSync() error {
	// Check if spirit is initialized
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		return fmt.Errorf("spirit not initialized. Run: spirit init")
	}

	// Check for git repo
	dotGit := filepath.Join(ConfigDir, ".git")
	if _, err := os.Stat(dotGit); os.IsNotExist(err) {
		// Initialize git repo
		fmt.Println("📦 Initializing git repository...")
		if err := gitInit(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	// Check for remote
	remoteURL, err := getRemoteURL()
	if err != nil || remoteURL == "" {
		// Prompt for remote setup
		fmt.Println("🔗 No remote configured.")
		fmt.Println("   Set up with: git remote add origin <url>")
		fmt.Println("   Or use: spirit sync --setup")
		return fmt.Errorf("no remote configured")
	}

	// Stage all changes
	fmt.Println("➕ Staging changes...")
	if err := gitAddAll(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Commit
	commitMsg := fmt.Sprintf("SPIRIT checkpoint: %s", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("💾 Creating commit...")
	if err := gitCommit(commitMsg); err != nil {
		// Might be no changes
		fmt.Println("   Nothing to commit (working tree clean)")
	}

	// Push
	fmt.Println("☁️ Pushing to remote...")
	if err := gitPush(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	fmt.Println("✅ Sync complete!")
	fmt.Printf("   Remote: %s\n", remoteURL)
	return nil
}

func gitInit() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

func getRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = ConfigDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func gitAddAll() error {
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

func gitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if nothing to commit
		if strings.Contains(string(output), "nothing to commit") {
			return nil
		}
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}

func gitPush() error {
	cmd := exec.Command("git", "push", "origin", "main")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try master if main fails
		if strings.Contains(string(output), "main") {
			cmd = exec.Command("git", "push", "origin", "master")
			cmd.Dir = ConfigDir
			output, err = cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("%s: %w", string(output), err)
			}
			return nil
		}
		return fmt.Errorf("%s: %w", string(output), err)
	}
	return nil
}
