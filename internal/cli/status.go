package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type SpiritStatus struct {
	Initialized   bool       `json:"initialized"`
	ConfigDir     string     `json:"config_dir"`
	Version       string     `json:"version"`
	LastBackup    *time.Time `json:"last_backup,omitempty"`
	TrackedFiles  int        `json:"tracked_files"`
	ExistingFiles int        `json:"existing_files"`
	GitConfigured bool       `json:"git_configured"`
	RemoteURL     string     `json:"remote_url,omitempty"`
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current state status",
		Long:  `Display detailed status of SPIRIT configuration, tracked files, and sync state.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return showStatus()
		},
	}
}

func showStatus() error {
	// Check if initialized
	configExists := false
	if _, err := os.Stat(ConfigDir); err == nil {
		configExists = true
	}

	if !configExists {
		fmt.Println("🌌 SPIRIT Status")
		fmt.Println()
		fmt.Println("   Status:     Not initialized")
		fmt.Printf("   Config dir: %s\n", ConfigDir)
		fmt.Println()
		fmt.Println("   Run 'spirit init' to get started")
		return nil
	}

	status := SpiritStatus{
		Initialized: true,
		ConfigDir:   ConfigDir,
		Version:     Version,
	}

	// Load config for version
	configPath := filepath.Join(ConfigDir, "spirit.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var config Config
		if err := json.Unmarshal(data, &config); err == nil {
			_ = config // could extract version from here if needed
		}
	}

	// Check tracked files
	trackedPath := filepath.Join(ConfigDir, ".spirit-tracked")
	if data, err := os.ReadFile(trackedPath); err == nil {
		var tracked TrackedConfig
		if err := json.Unmarshal(data, &tracked); err == nil {
			status.TrackedFiles = len(tracked.Files)
			// Count existing files
			for _, pattern := range tracked.Files {
				matches, _ := filepath.Glob(filepath.Join(ConfigDir, pattern))
				status.ExistingFiles += len(matches)
			}
		}
	}

	// Check git status
	dotGit := filepath.Join(ConfigDir, ".git")
	if _, err := os.Stat(dotGit); err == nil {
		status.GitConfigured = true

		// Get remote URL
		cmd := exec.Command("git", "remote", "get-url", "origin")
		cmd.Dir = ConfigDir
		if output, err := cmd.Output(); err == nil {
			status.RemoteURL = string(output)
		}

		// Get last commit time
		cmd = exec.Command("git", "log", "-1", "--format=%ct")
		cmd.Dir = ConfigDir
		if output, err := cmd.Output(); err == nil {
			var ts int64
			fmt.Sscanf(string(output), "%d", &ts)
			t := time.Unix(ts, 0)
			status.LastBackup = &t
		}
	}

	// Print status
	fmt.Println("🌌 SPIRIT Status")
	fmt.Println()
	fmt.Printf("   Version:    %s\n", status.Version)
	fmt.Printf("   Config:     %s\n", status.ConfigDir)
	fmt.Printf("   Initialized: Yes\n")

	if status.LastBackup != nil {
		ago := time.Since(*status.LastBackup)
		fmt.Printf("   Last sync:  %s ago\n", formatDuration(ago))
	} else {
		fmt.Println("   Last sync:  Never")
	}

	fmt.Println()
	fmt.Printf("   Tracked files:  %d patterns\n", status.TrackedFiles)
	fmt.Printf("   Existing files: %d matched\n", status.ExistingFiles)

	fmt.Println()
	if status.GitConfigured {
		fmt.Println("   Git: ✓ Configured")
		if status.RemoteURL != "" {
			fmt.Printf("   Remote: %s", status.RemoteURL)
		} else {
			fmt.Println("   Remote: ✗ Not configured")
			fmt.Println("           Run: git remote add origin <url>")
		}
	} else {
		fmt.Println("   Git: ✗ Not initialized")
	}

	fmt.Println()
	fmt.Println("   Commands:")
	fmt.Println("     spirit sync    - Push state to remote")
	fmt.Println("     spirit backup  - Create checkpoint + sync")

	return nil
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}
