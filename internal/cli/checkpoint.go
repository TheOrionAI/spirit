package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func checkpointCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "checkpoint",
		Short: "Create a manual checkpoint",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🌌 Checkpoint created (stub)")
		},
	}
}

func createCheckpoint(message string) error {
	return nil
}
