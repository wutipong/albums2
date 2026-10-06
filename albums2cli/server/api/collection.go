package api

import (
	"context"

	"github.com/wutipong/albums/albums2cli/server/types"
)

type CollectionListResponse struct {
	Collections []types.Collection `json:"collections"`
}

func GetCollectionList(ctx context.Context, server ServerConfig,
) (resp CollectionListResponse, err error) {
	c := NewClient(server)
	_, err = c.R().
		SetSuccessResult(&resp).
		SetContext(ctx).
		Get("api/collections")
	return
}

func GetCollection(ctx context.Context, server ServerConfig, collectionID string,
) (resp types.Collection, err error) {

	c := NewClient(server)
	_, err = c.R().
		SetSuccessResult(&resp).
		SetContext(ctx).
		Get("api/collections/" + collectionID)

	return
}

func CreateCollection(
	ctx context.Context,
	server ServerConfig,
	name string,
) (resp types.Collection, err error) {
	req := types.Collection{Name: name}

	c := NewClient(server)
	_, err = c.R().SetBodyJsonMarshal(req).
		SetSuccessResult(&resp).
		SetContext(ctx).
		Post("api/collections")

	return
}

type CollectionByNameResponse struct {
	Existed    bool             `json:"existed"`
	Collection types.Collection `json:"collection"`
}

func GetCollectionByName(
	ctx context.Context,
	server ServerConfig,
	collectionName string,
) (resp CollectionByNameResponse, err error) {
	c := NewClient(server)

	_, err = c.R().
		SetSuccessResult(&resp).
		SetContext(ctx).
		SetQueryParam("name", collectionName).
		Get("api/collections/by-name")

	return
}
