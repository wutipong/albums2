package collection

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/wutipong/albums/albums2cli/profile"
	"github.com/wutipong/albums/albums2cli/server/api"
)

func getCollectionById(ctx context.Context, profileName string, dryRun bool, collectionId string) (err error) {
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

	collection, err := api.GetCollection(ctx, server, collectionId)
	if err != nil {
		return fmt.Errorf("unable to get collection: %w", err)
	}

	slog.Info("collection",
		slog.String("id", collection.ID),
		slog.String("name", collection.Name),
	)

	return nil
}

func getCollectionByName(ctx context.Context, profileName string, dryRun bool, collectionName string) (err error) {
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

	resp, err := api.GetCollectionByName(ctx, server, collectionName)
	if err != nil {
		slog.Error("Failed to get collection by name", "error", err)
		return err
	}

	if resp.Existed {
		slog.Info("collection",
			slog.String("id", resp.Collection.ID),
			slog.String("name", resp.Collection.Name),
		)
	} else {
		slog.Info("collection not found",
			slog.String("name", collectionName),
		)
	}
	return nil
}
