package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path"
	"path/filepath"
	"time"

	"github.com/wutipong/albums2/gopkg/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostMediaResposnse struct {
	Media   types.Media `json:"Media"`
	Success bool        `json:"success"`
}

var (
	ErrDuplicateMedia = errors.New("duplicate Media")
)

func PostMedia(
	ctx context.Context,
	server ServerConfig,
	albumID string,
	containerPath string,
	filePath string,
	reader io.Reader,
	size int64,
) (result PostMediaResposnse, err error) {
	start := time.Now()
	defer func() {
		slog.Info("PostMedia completed",
			slog.Duration("duration", time.Since(start)),
			slog.String("path", filePath),
		)
	}()

	if ctx.Err() != nil {
		err = fmt.Errorf("context error: %w", ctx.Err())
		return
	}
	if server.DryRun {
		slog.Debug(
			"Dry run: skipping Media upload",
			slog.String("path", filePath),
		)
		result = PostMediaResposnse{
			Media: types.Media{
				ID: primitive.NewObjectID(),
			},
			Success: true,
		}

		return
	}

	c := NewClient(server)

	data, err := io.ReadAll(reader)
	if err != nil {
		err = fmt.Errorf("unable to read from reader: %w", err)
		return
	}

	mediaFileName := filepath.Join(containerPath, filePath)

	slog.Debug("upload request",
		slog.String("album_id", albumID),
		slog.String("filename", mediaFileName),
	)

	var postMediaRequest PostMediaRequestResponse
	_, err = c.R().SetBodyJsonMarshal(PostMediaRequestRequest{
		Filename: mediaFileName,
	}).
		SetSuccessResult(&postMediaRequest).
		SetContext(ctx).
		Post(path.Join("api", "albums", albumID, "media"))

	if err != nil {
		err = fmt.Errorf("request to upload failed: %w", err)
		return
	}

	slog.Debug("upload request results",
		slog.String("ID", postMediaRequest.ID),
		slog.String("url", postMediaRequest.URL),
		slog.Bool("success", postMediaRequest.Success),
	)

	slog.Debug("put object", slog.Int("size", len(data)), slog.String("url", postMediaRequest.URL))

	_, err = c.R().
		SetBodyBytes(data).
		SetRetryCount(10).
		SetContext(ctx).
		Put(postMediaRequest.URL)

	success := true
	if err != nil {
		slog.Error("put object fails", slog.String("error", err.Error()))
		success = false
	}

	slog.Debug("upload commit",
		slog.String("id", postMediaRequest.ID),
		slog.String("url", postMediaRequest.URL),
	)

	_, err = c.R().
		SetSuccessResult(&result).
		SetBodyJsonMarshal(PostMediaCommitRequest{
			Success: success,
		}).
		SetContext(ctx).
		Patch(path.Join("api", "media", postMediaRequest.ID))

	if err != nil {
		err = fmt.Errorf("unable to commit Media upload %s: %w", postMediaRequest.ID, err)
	}

	return
}

type PostMediaRequestRequest struct {
	Filename string `json:"filename"`
}

type PostMediaRequestResponse struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	Success bool   `json:"success"`
}

type PostMediaCommitRequest struct {
	ID      string `json:"id"`
	Success bool   `json:"success"`
}

type PostMediaCommitResponse struct {
	Media   types.Media `json:"Media"`
	Success bool        `json:"success"`
}
