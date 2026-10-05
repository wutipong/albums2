package collection

import (
	"context"

	"github.com/urfave/cli/v3"
)

func Command(profile *string) *cli.Command {
	dryRun := false
	byName := false
	byId := true
	strValue := ""

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
			}, {
				Name:  "create",
				Usage: "Create a new collection in the server. ",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "name",
						Destination: &strValue,
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},

				Action: func(ctx context.Context, cmd *cli.Command) error {
					return createCollection(ctx, *profile, dryRun, strValue)
				},
			}, {
				Name:  "get",
				Usage: "Get a collection by ID or name. ",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "name-or-id",
						Destination: &strValue,
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
					&cli.BoolFlag{
						Name:        "id",
						Value:       true,
						Usage:       "Find the collection by ID. This is the default behavior.",
						Destination: &byId,
					},
					&cli.BoolFlag{
						Name:        "name",
						Value:       false,
						Usage:       "Find the collection by name.",
						Destination: &byName,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if byName {
						return getCollectionByName(ctx, *profile, dryRun, strValue)
					}

					return getCollectionById(ctx, *profile, dryRun, strValue)
				},
			},
		},
	}
}
