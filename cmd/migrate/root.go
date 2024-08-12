package migrate

import (
	"github.com/henrywhitaker3/go-template/internal/app"
	"github.com/henrywhitaker3/go-template/internal/postgres"
	"github.com/spf13/cobra"
)

var (
	m *postgres.Migrator
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			mig, err := postgres.NewMigrator(app.Database)
			if err != nil {
				return err
			}
			m = mig
			return nil
		},
	}

	cmd.AddCommand(up())
	cmd.AddCommand(down())

	return cmd
}
