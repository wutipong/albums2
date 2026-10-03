package api

import (
	"context"

	"github.com/wutipong/albums/albums2cli/server/types"
)

type MeResponse struct {
	User    types.User    `json:"user"`
	Session types.Session `json:"session"`
}

func GetMe(ctx context.Context, server ServerConfig) (resp MeResponse, err error) {
	c := NewClient(server)
	r := c.Get("api/me").
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)

	err = r.Err
	return
}
