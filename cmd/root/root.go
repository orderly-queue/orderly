package root

import (
	"github.com/orderly-queue/orderly/cmd/migrate"
	"github.com/orderly-queue/orderly/cmd/routes"
	"github.com/orderly-queue/orderly/cmd/serve"
	"github.com/orderly-queue/orderly/internal/app"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "api",
		Short:   "Golang template API",
		Version: app.Version,
	}

	cmd.AddCommand(serve.New(app))
	cmd.AddCommand(migrate.New(app))
	cmd.AddCommand(routes.New(app))

	cmd.PersistentFlags().StringP("config", "c", "orderly.yaml", "The path to the api config file")

	return cmd
}
