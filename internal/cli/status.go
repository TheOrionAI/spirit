package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current state status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🌌 SPIRIT status: Active (stub)")
		},
	}
}
