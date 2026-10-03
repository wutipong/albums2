package api

import (
	"context"
	"net/url"

	"github.com/wutipong/albums/albums2cli/server/types"
)

type CollectionListResponse struct {
	Collections []types.Collection `json:"collections"`
}

func GetCollectionList(ctx context.Context, server ServerConfig) (resp CollectionListResponse, err error) {
	c := NewClient(server)
	c.Get("api/album").
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)
	return
}

func GetCollection(ctx context.Context, server ServerConfig, collectionID string) (resp types.Collection, err error) {
	c := NewClient(server)
	c.Get("api/album/" + collectionID).
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)
	return
}

type CollectionByNameResponse struct {
	Existed    bool             `json:"existed"`
	Collection types.Collection `json:"collection"`
}

func GetCollectionByName(ctx context.Context, server ServerConfig, collectionName string) (resp CollectionByNameResponse, err error) {
	c := NewClient(server)

	u := url.URL{
		Path: "api/album/by-name",
	}
	u.Query().Set("name", collectionName)

	c.Get(u.String()).
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)
	return
}
