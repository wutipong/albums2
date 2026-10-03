package process

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums/albums2cli/profile"
	"github.com/wutipong/albums/albums2cli/server/api"
)

func Command(profileStr *string) *cli.Command {
	dryRun := false
	id := ""

	return &cli.Command{
		Name:  "process",
		Usage: "Manually queueing process asset command.",
		Commands: []*cli.Command{
			{
				Name:  "single",
				Usage: "Process a single asset.. ",
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
					c, err := profile.LoadProfile(ctx, *profileStr)
					if err != nil {
						return fmt.Errorf(
							"unable to load configuration. please run 'albums2cli setup' first: %w",
							err,
						)
					}

					slog.Info("Albums instance",
						slog.String("url", c.URL),
					)

					url, err := url.Parse(c.URL)
					if err != nil {
						return fmt.Errorf("invalid Albums url: %w", err)
					}

					server := api.ServerConfig{
						URL:    url,
						DryRun: dryRun,
						APIKey: c.APIKey,
					}

					return processSingle(ctx, server, dryRun, id)
				},
			}, {
				Name:  "pending",
				Usage: "Process pending assets.. ",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "dry-run",
						Value:       false,
						Usage:       "Don't make actual API calls to the server. Useful for testing.",
						Destination: &dryRun,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					c, err := profile.LoadProfile(ctx, *profileStr)
					if err != nil {
						return fmt.Errorf(
							"unable to load configuration. please run 'albums2cli setup' first: %w",
							err,
						)
					}

					slog.Info("Albums instance",
						slog.String("url", c.URL),
					)

					url, err := url.Parse(c.URL)
					if err != nil {
						return fmt.Errorf("invalid Albums url: %w", err)
					}

					server := api.ServerConfig{
						URL:    url,
						DryRun: dryRun,
						APIKey: c.APIKey,
					}

					return processPending(ctx, server, dryRun)
				},
			},
		},
	}
}
