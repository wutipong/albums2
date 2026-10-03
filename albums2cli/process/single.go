package process

import (
	"context"
	"fmt"
	"path"

	"github.com/wutipong/albums/albums2cli/server/api"
)

func processSingle(ctx context.Context, server api.ServerConfig, dryRun bool, id string) error {
	if dryRun {
		return nil
	}

	c := api.NewClient(server)
	resp := c.Get(path.Join("api", "asset", id, "process")).
		Do(ctx)

	if resp.Err != nil {
		return fmt.Errorf("unable to queue process asset command: %w", resp.Err)
	}

	return nil
}
