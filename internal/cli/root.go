package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	Version   string
	ConfigDir string
)

func Execute(version string) error {
	Version = version
	ConfigDir = getConfigDir()

	rootCmd := &cobra.Command{
		Use:   "spirit",
		Short: "SPIRIT - State Preservation & Identity Resurrection Infrastructure Tool",
		Long: `🌌 SPIRIT

Resurrection-grade state management for AI agents.

Don't backup. Resurrect.

Complete documentation: https://spirit.theorionai.io`,
		Version: Version,
	}

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(migrateCmd())
	rootCmd.AddCommand(backupCmd())
	rootCmd.AddCommand(autoBackupCmd())
	rootCmd.AddCommand(checkpointCmd())
	rootCmd.AddCommand(restoreCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(statusCmd())

	return rootCmd.Execute()
}

func getConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".spirit"
	}
	return filepath.Join(home, ".spirit")
}
