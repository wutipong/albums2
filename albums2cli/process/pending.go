package process

import (
	"context"
	"fmt"

	"github.com/wutipong/albums2/albums2cli/server/api"
)

func processPending(ctx context.Context, server api.ServerConfig, dryRun bool) error {
	if dryRun {
		return nil
	}

	c := api.NewClient(server)
	resp := c.Get("/api/asset/scan").Do(ctx)

	if resp.Err != nil {
		return fmt.Errorf("unable to queue process pending asset command: %w", resp.Err)
	}

	return nil
}
