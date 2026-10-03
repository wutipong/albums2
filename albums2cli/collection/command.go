package collection

import (
	"context"

	"github.com/urfave/cli/v3"
)

func Command(profile *string) *cli.Command {
	dryRun := false

	return &cli.Command{
		Name:  "collection",
		Usage: "Manage collections in the server",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all collections in the server. ",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return listCollections(ctx, *profile, dryRun)
				},
			},
		},
	}
}
