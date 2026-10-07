package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
	"github.com/wutipong/albums2/albums2cli/album"
	"github.com/wutipong/albums2/albums2cli/collection"
	"github.com/wutipong/albums2/albums2cli/importing"
	"github.com/wutipong/albums2/albums2cli/log"
	"github.com/wutipong/albums2/albums2cli/process"
	"github.com/wutipong/albums2/albums2cli/profile"
)

func main() {
	slog.SetDefault(slog.New(tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
	})))

	profileStr := "default"
	debug := false

	cmd := &cli.Command{
		Name:  "albums2cli",
		Usage: "import assets to albums",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Value:       "default",
				Usage:       "profile of albums server.",
				Destination: &profileStr,
				Category:    "albums Server",
			},
			&cli.BoolFlag{
				Name:        "debug",
				Value:       false,
				Usage:       "enable loggin debug message",
				Destination: &debug,
			},
		},
		Commands: []*cli.Command{
			collection.Command(&profileStr),
			profile.Command(&profileStr),
			album.Command(&profileStr),
			process.Command(&profileStr),
			importing.Command(&profileStr),
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			level := slog.LevelInfo.String()
			if debug {
				level = slog.LevelDebug.String()
			}
			log.Setup(profileStr, level, true, level)
			return ctx, nil

		},
	}
	defer log.CleanUp()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := cmd.Run(ctx, os.Args); err != nil {
		slog.Error("Error running command", slog.String("error", err.Error()))
	}
}
