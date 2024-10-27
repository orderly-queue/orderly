package serve

import (
	"context"

	"github.com/orderly-queue/orderly/internal/app"
	"github.com/spf13/cobra"
)

func New(app *app.App) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the api server",
		RunE: func(cmd *cobra.Command, args []string) error {
			go app.Probes.Start(cmd.Context())

			go func() {
				<-cmd.Context().Done()
				ctx := context.Background()
				app.Probes.Unready()

				app.Metrics.Stop(ctx)
				app.Probes.Stop(ctx)
				app.Http.Stop(ctx)
			}()

			go app.Metrics.Start(cmd.Context())

			app.Probes.Ready()
			app.Probes.Healthy()

			return app.Http.Start(cmd.Context())
		},
	}
}
