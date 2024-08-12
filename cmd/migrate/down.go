package migrate

import (
	"github.com/spf13/cobra"
)

func down() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Run the down migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Down()
		},
	}
}
