package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// TrackedConfig represents the .spirit-tracked file
type TrackedConfig struct {
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

func syncCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "Sync state to configured backend",
		Long:  `Push current state to remote repository.

Only files that exist will be synced. Missing files are skipped silently.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			return runSync(verbose)
		},
	}
}

func runSync(verbose bool) error {
	if verbose {
		fmt.Println("🔍 Verbose mode enabled")
	}

	// Check if spirit is initialized
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		return fmt.Errorf("spirit not initialized. Run: spirit init")
	}

	// Check for git repo
	dotGit := filepath.Join(ConfigDir, ".git")
	if _, err := os.Stat(dotGit); os.IsNotExist(err) {
		fmt.Println("📦 Initializing git repository...")
		if err := gitInit(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}
	}

	// Check for remote
	remoteURL, err := getRemoteURL()
	if err != nil || remoteURL == "" {
		fmt.Println("🔗 No remote configured.")
		fmt.Println("   Set up with: git remote add origin <url>")
		return fmt.Errorf("no remote configured")
	}

	// Load tracked files
	tracked, err := loadTrackedFiles()
	if err != nil {
		if verbose {
			fmt.Printf("   ⚠️  Using default tracked files: %v\n", err)
		}
		// Fallback to defaults
		tracked = []string{
			"IDENTITY.md", "SOUL.md", "AGENTS.md", "TOOLS.md",
			"PROJECTS.md", "HEARTBEAT.md", "README.md",
			"spirit.json", ".spirit-tracked",
		}
	}

	// Check which tracked files exist
	existingFiles := []string{}
	missingFiles := []string{}

	for _, pattern := range tracked {
		// Handle wildcards
		if strings.Contains(pattern, "*") {
			matches, _ := filepath.Glob(filepath.Join(ConfigDir, pattern))
			for _, match := range matches {
				if _, err := os.Stat(match); err == nil {
					relPath, _ := filepath.Rel(ConfigDir, match)
					existingFiles = append(existingFiles, relPath)
				}
			}
		} else {
			// Direct file check
			fullPath := filepath.Join(ConfigDir, pattern)
			if _, err := os.Stat(fullPath); err == nil {
				existingFiles = append(existingFiles, pattern)
			} else {
				missingFiles = append(missingFiles, pattern)
			}
		}
	}

	if verbose {
		fmt.Printf("📁 Found %d files to sync:\n", len(existingFiles))
		for _, f := range existingFiles {
			fmt.Printf("   ✓ %s\n", f)
		}
		if len(missingFiles) > 0 {
			fmt.Printf("   ⏭️  Skipped %d missing files\n", len(missingFiles))
			for _, f := range missingFiles {
				fmt.Printf("     - %s (not found)\n", f)
			}
		}
	} else {
		fmt.Printf("📦 Syncing %d files...\n", len(existingFiles))
	}

	if len(existingFiles) == 0 {
		fmt.Println("⚠️  No files to sync")
		return nil
	}

	// Stage existing files
	fmt.Println("➕ Staging changes...")
	if err := gitAddFiles(existingFiles); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Fetch remote changes first (for multi-agent sync)
	fmt.Println("📥 Fetching remote...")
	gitFetch() // Best effort, ignore errors

	// Pull/merge remote changes to avoid conflicts
	fmt.Println("🔄 Syncing with remote...")
	if err := gitPull(); err != nil {
		return fmt.Errorf("sync failed (conflict?): %w", err)
	}

	// Commit
	commitMsg := fmt.Sprintf("SPIRIT checkpoint: %s (%d files)",
		time.Now().Format("2006-01-02 15:04:05"), len(existingFiles))
	fmt.Println("💾 Creating commit...")
	if err := gitCommit(commitMsg); err != nil {
		// Might be no changes
		if strings.Contains(err.Error(), "nothing") {
			fmt.Println("✅ Already up to date")
			return nil
		}
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Push
	fmt.Println("☁️  Pushing to remote...")
	if err := gitPush(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	fmt.Println("✅ Sync complete!")
	fmt.Printf("   Remote: %s\n", remoteURL)
	fmt.Printf("   Files: %d\n", len(existingFiles))
	return nil
}

func loadTrackedFiles() ([]string, error) {
	trackedPath := filepath.Join(ConfigDir, ".spirit-tracked")
	data, err := os.ReadFile(trackedPath)
	if err != nil {
		return nil, err
	}

	var config TrackedConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return config.Files, nil
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

func gitAddFiles(files []string) error {
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
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
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func gitFetch() error {
	cmd := exec.Command("git", "fetch", "origin")
	cmd.Dir = ConfigDir
	// Fetch may fail if no remote, that's ok
	_, _ = cmd.CombinedOutput()
	return nil
}

func gitPull() error {
	// Try to pull/rebase from main
	cmd := exec.Command("git", "pull", "--rebase", "origin", "main")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Try master if main fails
		if strings.Contains(string(output), "main") || strings.Contains(string(output), "couldn't find remote ref") {
			cmd = exec.Command("git", "pull", "--rebase", "origin", "master")
			cmd.Dir = ConfigDir
			output, err = cmd.CombinedOutput()
			if err != nil {
				// Might be nothing to pull (first push)
				if strings.Contains(string(output), "no such ref") || strings.Contains(string(output), "could not resolve") {
					return nil // First time, no remote commits yet
				}
				return fmt.Errorf("git pull failed: %s", string(output))
			}
		} else if strings.Contains(string(output), "no such ref") || strings.Contains(string(output), "could not resolve") {
			// First time, no remote commits
			return nil
		} else {
			return fmt.Errorf("git pull failed: %s", string(output))
		}
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
				return fmt.Errorf("git push failed: %s", string(output))
			}
			return nil
		}
		return fmt.Errorf("git push failed: %s", string(output))
	}
	return nil
}
