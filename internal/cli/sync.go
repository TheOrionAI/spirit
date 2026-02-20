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

type TrackedConfig struct {
	Version string   `json:"version"`
	Files   []string `json:"files"`
}

func syncCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync state to configured backend",
		Long: `Push current state to remote repository.

Environment variable SPIRIT_SOURCE_DIR:
  When set, reads files from this directory instead of ~/.spirit/
  The .spirit-tracked config is still read from ~/.spirit/ (which may be a symlink)

Examples:
  spirit sync                          # Sync from ~/.spirit/
  SPIRIT_SOURCE_DIR=/workspace spirit sync   # Sync from /workspace
  export SPIRIT_SOURCE_DIR=/workspace
  spirit sync                          # Sync from /workspace
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool("verbose")
			return runSync(verbose)
		},
	}
	cmd.Flags().Bool("verbose", false, "Verbose output")
	return cmd
}

func runSync(verbose bool) error {
	sourceDir := getSourceDir()

	if verbose {
		fmt.Printf("🔍 Source directory: %s\n", sourceDir)
		if sourceDir != ConfigDir {
			fmt.Printf("   (via SPIRIT_SOURCE_DIR)\n")
		}
	}

	// Check if spirit is initialized
	if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
		return fmt.Errorf("spirit not initialized. Run: spirit init")
	}

	// Check for git repo in ConfigDir (not sourceDir)
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
		fmt.Println("🔗 No remote configured. Set up with:")
		fmt.Println("   cd ~/.spirit && git remote add origin <url>")
		return fmt.Errorf("no remote configured")
	}

	// Load tracked files from ConfigDir (or via symlink)
	tracked, err := loadTrackedFiles()
	if err != nil {
		if verbose {
			fmt.Printf("⚠️  Using defaults: %v\n", err)
		}
		tracked = []string{
			"IDENTITY.md", "SOUL.md", "AGENTS.md", "TOOLS.md", "PROJECTS.md",
			"HEARTBEAT.md", "README.md", "spirit.json", ".spirit-tracked",
			"memory/*.md", "projects/*.md", "context/*.md",
		}
	}

	// Clear ConfigDir (except git) and copy fresh files from sourceDir
	// This ensures the git repo always reflects the current workspace state
	existingFiles := []string{}
	missingFiles := []string{}

	for _, pattern := range tracked {
		// Handle wildcards
		if strings.Contains(pattern, "*") {
			// Glob in source directory
			matches, _ := filepath.Glob(filepath.Join(sourceDir, pattern))
			for _, match := range matches {
				if info, err := os.Stat(match); err == nil && !info.IsDir() {
					relPath, _ := filepath.Rel(sourceDir, match)
					// Copy to ConfigDir for sync
					targetPath := filepath.Join(ConfigDir, relPath)
					if err := copyFile(match, targetPath); err == nil {
						existingFiles = append(existingFiles, relPath)
					}
				}
			}
		} else {
			// Direct file
			sourcePath := filepath.Join(sourceDir, pattern)
			if _, err := os.Stat(sourcePath); err == nil {
				targetPath := filepath.Join(ConfigDir, pattern)
				if err := copyFile(sourcePath, targetPath); err == nil {
					existingFiles = append(existingFiles, pattern)
				} else if verbose {
					fmt.Printf("⚠️  Failed to copy %s: %v\n", pattern, err)
				}
			} else {
				// File doesn't exist in source
				// Check if it exists in ConfigDir (maybe user added it manually)
				configPath := filepath.Join(ConfigDir, pattern)
				if _, err := os.Stat(configPath); err == nil {
					existingFiles = append(existingFiles, pattern)
				} else {
					missingFiles = append(missingFiles, pattern)
				}
			}
		}
	}

	// Always include .spirit-tracked itself
	if _, err := os.Stat(filepath.Join(ConfigDir, ".spirit-tracked")); err == nil {
		// Already exists
	}

	if verbose {
		fmt.Printf("📁 Found %d files to sync:\n", len(existingFiles))
		for _, f := range existingFiles {
			fmt.Printf("   ✓ %s\n", f)
		}
		if len(missingFiles) > 0 {
			fmt.Printf("⏭️  Skipped %d missing:\n", len(missingFiles))
			for _, f := range missingFiles {
				fmt.Printf("   - %s\n", f)
			}
		}
	} else {
		fmt.Printf("📦 Syncing %d files...\n", len(existingFiles))
	}

	if len(existingFiles) == 0 {
		return fmt.Errorf("no files to sync (check .spirit-tracked or SPIRIT_SOURCE_DIR)")
	}

	// Stage and sync
	fmt.Println("➕ Staging changes...")
	if err := gitAddAll(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	fmt.Println("📥 Fetching remote...")
	gitFetch()

	fmt.Println("🔄 Syncing with remote...")
	if err := gitPull(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	commitMsg := fmt.Sprintf("SPIRIT sync: %s (%d files)", time.Now().Format("2006-01-02 15:04"), len(existingFiles))
	fmt.Println("💾 Creating commit...")
	if err := gitCommit(commitMsg); err != nil {
		if strings.Contains(err.Error(), "nothing") || strings.Contains(err.Error(), "No changes") {
			fmt.Println("✅ Already up to date")
			return nil
		}
		return fmt.Errorf("git commit failed: %w", err)
	}

	fmt.Println("☁️ Pushing to remote...")
	if err := gitPush(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	fmt.Println("✅ Sync complete!")
	fmt.Printf("   Remote: %s\n", remoteURL)
	fmt.Printf("   Files: %d\n", len(existingFiles))
	if sourceDir != ConfigDir {
		fmt.Printf("   Source: %s\n", sourceDir)
	}
	return nil
}

func copyFile(src, dst string) error {
	// Ensure directory exists
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// Copy file content
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0644)
}

func loadTrackedFiles() ([]string, error) {
	trackedPath := filepath.Join(ConfigDir, ".spirit-tracked")
	data, err := os.ReadFile(trackedPath)
	if err != nil {
		// Check if it's a symlink and read target
		if info, err := os.Lstat(trackedPath); err == nil && info.Mode()&os.ModeSymlink != 0 {
			target, _ := os.Readlink(trackedPath)
			data, err = os.ReadFile(target)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
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

func gitAddAll() error {
	cmd := exec.Command("git", "add", "-A")
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
	output, err :=	cmd.CombinedOutput()
	if err != nil {
		outputStr := strings.TrimSpace(string(output))
		if !strings.Contains(outputStr, "could not resolve") && 
		   !strings.Contains(outputStr, "does not appear to be") &&
		   !strings.Contains(outputStr, "No remote repository") {
			return fmt.Errorf("git fetch failed: %s", outputStr)
		}
	}
	return nil
}

func gitPull() error {
	cmd := exec.Command("git", "pull", "--rebase", "origin", "main")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "main") || strings.Contains(string(output), "couldn't find") {
			cmd = exec.Command("git", "pull", "--rebase", "origin", "master")
			cmd.Dir = ConfigDir
			output, err = cmd.CombinedOutput()
			if err != nil {
				if strings.Contains(string(output), "no such ref") || strings.Contains(string(output), "could not resolve") {
					return nil
				}
				return fmt.Errorf("git pull failed: %s", string(output))
			}
		} else if strings.Contains(string(output), "no such ref") || strings.Contains(string(output), "could not resolve") {
			return nil
		} else {
			return fmt.Errorf("git pull failed: %s", string(output))
		}
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

func gitPush() error {
	cmd := exec.Command("git", "push", "origin", "main")
	cmd.Dir = ConfigDir
	output, err := cmd.CombinedOutput()
	if err != nil {
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
