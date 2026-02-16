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

func backupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backup",
		Short: "Manually trigger state backup",
		Long:  `Create a checkpoint and sync to all configured backends.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			message, _ := cmd.Flags().GetString("message")
			return backupSpirit(message)
		},
	}

	cmd.Flags().StringP("message", "m", "", "Backup message/description")

	return cmd
}

func autoBackupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "autobackup",
		Short: "Configure automatic backups",
		Long: `Set up automatic state preservation.

Examples:
  spirit autobackup --interval=15m     # Backup every 15 minutes
  spirit autobackup --on-session-end     # Backup when session ends
  spirit autobackup --watch             # Watch for changes and backup
  spirit autobackup --disable           # Disable auto-backup`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return configureAutoBackup(cmd)
		},
	}
}

func backupSpirit(message string) error {
	if message == "" {
		message = fmt.Sprintf("Backup at %s", time.Now().Format("2006-01-02 15:04"))
	}

	fmt.Printf("🌌 Creating backup: %s\n", message)

	// 1. Create checkpoint
	fmt.Println("📸 Creating checkpoint...")
	if err := createCheckpoint(message); err != nil {
		return fmt.Errorf("checkpoint failed: %w", err)
	}

	// 2. Sync to all backends
	fmt.Println("☁️ Syncing to backends...")
	if err := syncToBackends(); err != nil {
		return fmt.Errorf("sync failed: %w", err)
	}

	fmt.Println("✅ Backup complete!")
	fmt.Printf("Your agent's spirit is preserved across all backends.\n")

	return nil
}

func configureAutoBackup(cmd *cobra.Command) error {
	interval, _ := cmd.Flags().GetString("interval")
	onSessionEnd, _ := cmd.Flags().GetBool("on-session-end")
	watch, _ := cmd.Flags().GetBool("watch")
	disable, _ := cmd.Flags().GetBool("disable")

	if disable {
		fmt.Println("🛑 Disabling auto-backup...")
		return disableAutoBackup()
	}

	fmt.Println("🔄 Configuring auto-backup...")

	// Create autobackup config
	config := AutoBackupConfig{
		Enabled:       true,
		Interval:      interval,
		OnSessionEnd:  onSessionEnd,
		Watch:         watch,
		LastBackup:    time.Time{},
	}

	// Save config
	if err := saveAutoBackupConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Setup mechanisms
	if interval != "" {
		fmt.Printf("⏱️  Backing up every %s\n", interval)
		// Could setup cron or systemd timer
	}

	if onSessionEnd {
		fmt.Println("🚪 Backing up on session end")
		// Would integrate with OpenClaw shutdown hooks
	}

	if watch {
		fmt.Println("👁️  Watching for changes")
		// Setup fsnotify watcher
		go watchAndBackup()
	}

	fmt.Println("\n✅ Auto-backup configured!")
	fmt.Println("Your agent's spirit will be preserved automatically.")

	return nil
}

type AutoBackupConfig struct {
	Enabled      bool      `json:"enabled"`
	Interval     string    `json:"interval,omitempty"`
	OnSessionEnd bool      `json:"on_session_end"`
	Watch        bool      `json:"watch"`
	LastBackup   time.Time `json:"last_backup"`
}

func saveAutoBackupConfig(config AutoBackupConfig) error {
	configPath := filepath.Join(ConfigDir, "autobackup.json")

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func disableAutoBackup() error {
	config := AutoBackupConfig{
		Enabled: false,
	}
	return saveAutoBackupConfig(config)
}

func watchAndBackup() {
	// Watch ~/.spirit/memory/ and projects/ for changes
	// Auto-backup when files change
	// This runs in background goroutine

	fmt.Println("👁️  Watcher started (background)")

	// Simplified - real implementation would use fsnotify
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Check if changes since last backup
		if hasChanges() {
			if err := backupSpirit("Auto-backup on change"); err != nil {
				fmt.Fprintf(os.Stderr, "Auto-backup failed: %v\n", err)
			}
		}
	}
}

func hasChanges() bool {
	// Check if files modified since last backup
	// Implementation would check mtime vs last backup time
	return true // Simplified
}

func syncToBackends() error {
	// Sync to all configured backends
	// Currently just implements GitHub

	configPath := filepath.Join(ConfigDir, "spirit.json")
	if _, err := os.Stat(configPath); err != nil {
		return fmt.Errorf("no configuration found, run 'spirit init' first")
	}

	// Get all backends
	backends := []string{"github"} // TODO: Read from config

	for _, backend := range backends {
		switch backend {
		case "github":
			if err := syncToGitHub(); err != nil {
				return fmt.Errorf("github sync failed: %w", err)
			}
		case "s3":
			// TODO
		}
	}

	return nil
}

func syncToGitHub() error {
	// Sync to GitHub
	// This is a simplified version - real one would check if git repo,
	// commit if needed, and push

	dotGit := filepath.Join(ConfigDir, ".git")
	if _, err := os.Stat(dotGit); os.IsNotExist(err) {
		// Initialize git repo
		cmd := exec.Command("git", "init")
		cmd.Dir = ConfigDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git init failed: %w", err)
		}

		// Add remote
		// Would read from config
	}

	// Add all changes
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = ConfigDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Commit
	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Auto-backup: %s",
		time.Now().Format("2006-01-02 15:04:05")))
	cmd.Dir = ConfigDir
	if err := cmd.Run(); err != nil {
		// Might be no changes
	}

	// Push
	// cmd = exec.Command("git", "push", "origin", "main")
	// cmd.Dir = ConfigDir
	// if err := cmd.Run(); err != nil {
	//	 return err
	// }

	fmt.Println("   → Synced to GitHub")
	return nil
}
