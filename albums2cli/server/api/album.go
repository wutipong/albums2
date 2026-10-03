package api

import (
	"context"
	"path"

	"github.com/wutipong/albums/albums2cli/server/types"
)

type AlbumListResponse struct {
	Albums []types.Album `json:"albums"`
}

func GetAlbumList(ctx context.Context, server ServerConfig) (resp AlbumListResponse, err error) {
	c := NewClient(server)
	c.Get("api/album").
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)
	return
}

type AlbumDetailResponse struct {
	types.Album
	Assets []string `json:"assets"`
}

func GetAlbumDetail(ctx context.Context, server ServerConfig, albumID string) (resp AlbumDetailResponse, err error) {
	c := NewClient(server)
	c.Get(path.Join("api", "album", albumID)).
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)
	return
}

type CreateAlbumRequest struct {
	Name string `json:"name"`
}

func CreateAlbum(
	ctx context.Context,
	server ServerConfig,
	name string,
) (resp types.Album, err error) {
	req := CreateAlbumRequest{Name: name}

	c := NewClient(server)
	r := c.Post("api/album").
		SetBodyJsonMarshal(req).
		SetSuccessResult(&resp).
		SetErrorResult(err).
		Do(ctx)

	err = r.Err

	return
}

type DeleteAlbumResponse struct {
	Success bool `json:"success"`
}

func DeleteAlbum(
	ctx context.Context,
	server ServerConfig,
	id string,
) (resp DeleteAlbumResponse, err error) {
	c := NewClient(server)
	r := c.Delete(path.Join("api", "album", id)).
		SetSuccessResult(&resp).
		Do(ctx)

	err = r.Err
	return
}

type PopulateAlbumCoverRequest struct {
	AssetID string `json:"asset_id"`
}

type PopulateAlbumCoverResponse struct {
	Status string `json:"status"`
}

func PopulateAlbumCover(
	ctx context.Context,
	server ServerConfig,
	albumID string,
) (resp PopulateAlbumCoverResponse, err error) {
	req := PopulateAlbumCoverRequest{
		AssetID: "",
	}

	c := NewClient(server)
	r := c.Post(path.Join("api", "album", albumID, "cover")).
		SetBodyJsonMarshal(req).
		SetSuccessResult(&resp).
		Do(ctx)

	err = r.Err

	return
}
