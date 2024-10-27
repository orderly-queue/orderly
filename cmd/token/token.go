package token

import (
	"fmt"
	"time"

	"github.com/orderly-queue/orderly/internal/app"
	"github.com/spf13/cobra"
)

var (
	expires time.Duration
)

func New(app *app.App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Generate an authentication token",
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := app.Jwt.New(expires)
			if err != nil {
				return err
			}

			fmt.Println(token)

			return nil
		},
	}

	cmd.Flags().DurationVar(&expires, "expiry", time.Hour*24*90, "The length of time the token is valid for")

	return cmd
}
