package profile

import (
	"context"
	"log/slog"
)

func showProfile(ctx context.Context, profile string) (err error) {
	config, err := LoadProfile(ctx, profile)
	if err != nil {
		return err
	}

	slog.Info("profile",
		slog.String("name", profile),
		slog.String("url", config.URL),
		slog.String("api_key", config.APIKey),
	)

	return nil
}
