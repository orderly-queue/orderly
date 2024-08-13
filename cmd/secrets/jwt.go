package secrets

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jwtSize int
)

func newJwt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "Generate a new JWT secret",
		RunE: func(cmd *cobra.Command, args []string) error {
			secret := make([]byte, jwtSize/8)
			_, err := rand.Read(secret)
			if err != nil {
				return err
			}

			encoded := base64.RawStdEncoding.EncodeToString(secret)
			fmt.Printf("base64:%s\n", encoded)
			return nil
		},
	}

	cmd.Flags().IntVar(&jwtSize, "size", 256, "The size in bits of the jwt sercet")

	return cmd
}
