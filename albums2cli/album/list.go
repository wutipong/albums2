package album

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums2/albums2cli/profile"
	"github.com/wutipong/albums2/albums2cli/server/api"
)

func listAlbum(ctx context.Context, profileName string, collection string, dryRun bool) (err error) {
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

	r, err := api.GetCollectionByName(ctx, server, collection)
	if err != nil {
		return fmt.Errorf("unable to find collection: %w", err)
	}

	if !r.Existed {
		return fmt.Errorf("collection does not exist.")
	}

	collectionID := r.Collection.ID

	albumList, err := api.GetAlbumList(ctx, server, collectionID)
	if err != nil {
		return err
	}

	for _, album := range albumList.Albums {
		slog.Info("album",
			slog.String("id", album.ID.String()),
			slog.String("name", album.Name),
		)
	}

	return nil
}
