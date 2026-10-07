package collection

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums2/albums2cli/profile"
	"github.com/wutipong/albums2/albums2cli/server/api"
)

func listCollections(ctx context.Context, profileName string, dryRun bool) (err error) {
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

	collectionList, err := api.GetCollectionList(ctx, server)
	if err != nil {
		slog.Error("Failed to get collection list", "error", err)
		return err
	}

	for _, collection := range collectionList.Collections {
		slog.Info("collections",
			slog.String("id", collection.ID),
			slog.String("name", collection.Name),
		)
	}

	return nil
}
