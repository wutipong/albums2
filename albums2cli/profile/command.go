package profile

import (
	"context"

	"github.com/urfave/cli/v3"
)

func Command(profile *string) *cli.Command {
	return &cli.Command{
		Name:  "profile",
		Usage: "Manage profiles for albums importer",
		Commands: []*cli.Command{
			{
				Name:  "me",
				Usage: "Show the user information associated to the profile.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return profileDetail(ctx, *profile, false)
				},
			},
			{
				Name: "setup",
				Usage: "Setup profile file interactively. " +
					"Existing configuration file will be overwritten.",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return setupProfile(ctx, *profile)
				},
			}, {
				Name:  "show",
				Usage: "Show profile configuration",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return showProfile(ctx, *profile)
				},
			},
		},
	}
}
