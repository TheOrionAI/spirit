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

func checkpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkpoint [message]",
		Short: "Create a manual checkpoint",
		Long:  `Create a git checkpoint of the current SPIRIT state without pushing to remote. Useful for local saves before experimenting.`,
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := "Manual checkpoint"
			if len(args) > 0 {
				message = args[0]
			}
			return createCheckpoint(message)
		},
	}
}

func createCheckpoint(message string) error {
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

	// Load tracked files
	tracked, err := loadTrackedFiles()
	if err != nil {
		// Fallback to defaults
		tracked = []string{
			"IDENTITY.md", "SOUL.md", "AGENTS.md", "TOOLS.md",
			"PROJECTS.md", "HEARTBEAT.md", "README.md",
			"spirit.json", ".spirit-tracked", "memory/*.md",
			"projects/*.md", "context/*.md",
		}
	}

	// Check which tracked files exist
	existingFiles := []string{}
	missingFiles := []string{}
	for _, pattern := range tracked {
		if strings.Contains(pattern, "*") {
			matches, _ := filepath.Glob(filepath.Join(ConfigDir, pattern))
			for _, match := range matches {
				if _, err := os.Stat(match); err == nil {
					relPath, _ := filepath.Rel(ConfigDir, match)
					existingFiles = append(existingFiles, relPath)
				}
			}
		} else {
			fullPath := filepath.Join(ConfigDir, pattern)
			if _, err := os.Stat(fullPath); err == nil {
				existingFiles = append(existingFiles, pattern)
			} else {
				missingFiles = append(missingFiles, pattern)
			}
		}
	}

	if len(existingFiles) == 0 {
		return fmt.Errorf("no files to checkpoint")
	}

	// Stage files
	fmt.Printf("➕ Staging %d files...\n", len(existingFiles))
	if err := gitAddFiles(existingFiles); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	// Create commit with provided message
	commitMsg := fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), message)
	fmt.Println("💾 Creating checkpoint...")
	if err := gitCommit(commitMsg); err != nil {
		if strings.Contains(err.Error(), "nothing") {
			fmt.Println("✅ Already up to date (no changes)")
			return nil
		}
		return fmt.Errorf("git commit failed: %w", err)
	}

	// Get commit hash for reference
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = ConfigDir
	output, _ := cmd.Output()
	commitHash := strings.TrimSpace(string(output))

	fmt.Printf("🌌 Checkpoint created: %s\n", commitHash)
	fmt.Printf("   Message: %s\n", message)
	fmt.Printf("   Files: %d tracked\n", len(existingFiles))
	
	// Show hint about sync if remote exists
	if remoteURL, _ := getRemoteURL(); remoteURL != "" {
		fmt.Println("\n   Tip: Run 'spirit sync' to push to remote")
	} else {
		fmt.Println("\n   Tip: Set up remote with: git remote add origin <url>")
	}

	return nil
}
