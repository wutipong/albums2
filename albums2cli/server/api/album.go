package api

import (
	"context"
	"path"

	"github.com/wutipong/albums2/gopkg/types"
)

type AlbumListResponse struct {
	Albums []types.Album `json:"albums"`
}

func GetAlbumList(
	ctx context.Context,
	server ServerConfig,
	collectionID string,
) (resp AlbumListResponse, err error) {
	c := NewClient(server)
	_, err = c.R().
		SetSuccessResult(&resp).
		SetContext(ctx).
		Get(path.Join("api", "collections", collectionID, "albums"))

	return
}

type AlbumDetailResponse struct {
	types.Album
}

func GetAlbum(ctx context.Context, server ServerConfig, albumID string) (resp AlbumDetailResponse, err error) {
	c := NewClient(server)
	_, err = c.R().SetSuccessResult(&resp).
		SetContext(ctx).
		Get(path.Join("api", "album", albumID))
	return
}

type CreateAlbumRequest struct {
	Name string `json:"name"`
}

func CreateAlbum(
	ctx context.Context,
	server ServerConfig,
	collectionID string,
	name string,
) (resp types.Album, err error) {
	req := CreateAlbumRequest{
		Name: name,
	}

	c := NewClient(server)
	_, err = c.R().
		SetBodyJsonMarshal(req).
		SetSuccessResult(&resp).
		SetContext(ctx).
		Post(path.Join("api", "collections", collectionID, "albums"))

	return
}
