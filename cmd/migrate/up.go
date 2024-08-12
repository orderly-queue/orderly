package migrate

import (
	"github.com/spf13/cobra"
)

func up() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Run the up migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return m.Up()
		},
	}
}
