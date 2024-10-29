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

			if app.Config.Queue.Snapshot.Enabled {
				snapshot, err := app.Snapshotter.Latest(cmd.Context())
				if err != nil {
					return err
				}
				if snapshot != nil {
					state, err := app.Snapshotter.Open(cmd.Context(), *snapshot)
					if err != nil {
						return err
					}
					app.Queue.Load(state)
				}
				if err := app.Snapshotter.Work(cmd.Context()); err != nil {
					return err
				}
				// Run a snapshot before shutdown so we don't load up stale data
				defer app.Snapshotter.Snapshot(context.Background())
			}

			go app.Queue.Report(cmd.Context())

			go app.Metrics.Start(cmd.Context())

			app.Probes.Ready()
			app.Probes.Healthy()

			return app.Http.Start(cmd.Context())
		},
	}
}
