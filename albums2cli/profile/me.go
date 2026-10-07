package profile

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums2/albums2cli/server/api"
)

func profileDetail(ctx context.Context, profileName string, dryRun bool) (err error) {
	config, err := LoadProfile(ctx, profileName)
	if err != nil {
		return err
	}

	slog.Info("profile",
		slog.String("name", profileName),
		slog.String("url", config.URL),
	)

	serverUrl, err := url.Parse(config.URL)
	if err != nil {
		return fmt.Errorf("Unable to parse server URL: %w", err)
	}
	server := api.ServerConfig{
		URL:    serverUrl,
		DryRun: dryRun,
		APIKey: config.APIKey,
	}

	me, err := api.GetMe(ctx, server)
	if err != nil {
		return err
	}

	slog.Info("session",
		slog.String("id", me.Session.ID),
	)

	slog.Info("user",
		slog.String("id", me.User.ID),
		slog.String("name", me.User.Name),
		slog.String("email", me.User.Email),
	)

	return nil
}
