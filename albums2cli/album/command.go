package album

import (
	"context"

	"github.com/urfave/cli/v3"
)

func Command(profile *string) *cli.Command {
	dryRun := false
	id := ""
	collectionId := ""

	return &cli.Command{
		Name:  "album",
		Usage: "Manage albums in the server",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all albums in the server. ",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
					&cli.StringFlag{
						Name:        "collection",
						Value:       "",
						Usage:       "Collection ID to list albums from.",
						Destination: &collectionId,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return listAlbum(ctx, *profile, collectionId, dryRun)
				},
			}, {
				Name:  "show",
				Usage: "Show album details",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "id",
						UsageText:   "Album ID",
						Destination: &id,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return showAlbum(ctx, *profile, dryRun, id)
				},
			},
		},
	}
}
