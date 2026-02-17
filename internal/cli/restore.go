package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func restoreCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restore",
		Short: "Restore state from backup",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🌌 Restore complete (stub)")
		},
	}
}
