package album

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums/albums2cli/profile"
	"github.com/wutipong/albums/albums2cli/server/api"
)

func deleteAlbum(ctx context.Context, profileName string, dryRun bool, id string) (err error) {
	config, err := profile.LoadProfile(ctx, profileName)
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

	resp, err := api.DeleteAlbum(ctx, server, id)
	if err != nil {
		return err
	}

	slog.Info("album delete result",
		slog.Bool("sucess", resp.Success),
	)

	return nil
}
