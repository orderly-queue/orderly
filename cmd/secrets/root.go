package secrets

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Secrets utilities",
	}

	cmd.AddCommand(newKey())
	cmd.AddCommand(newJwt())

	return cmd
}
