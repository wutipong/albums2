package importing

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/wutipong/albums/albums2cli/server/api"
	"github.com/wutipong/albums/albums2cli/server/types"
)

func ProcessDirectory(
	ctx context.Context,
	server api.ServerConfig,
	album types.Album,
	path string,
) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context error: %w", ctx.Err())
	}
	slog.Debug("processing directory",
		slog.String("path", path),
	)

	filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Error(
				"failed to access path",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
		}
		if d.IsDir() {
			return nil
		}

		if IsMediaFile(path) {
			err = processMediaFile(ctx, server, path, album)
		} else if IsArchiveFile(path) {
			err = ProcessArchive(ctx, server, album, path)
		}

		if err != nil {
			slog.Error(
				"failed to process file",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
			return nil
		}
		return nil
	})

	return nil
}

func processMediaFile(
	ctx context.Context,
	server api.ServerConfig,
	path string,
	album types.Album,
) error {
	slog.Debug("processing media file",
		slog.String("path", path),
	)

	if !IsMediaFile(path) {
		slog.Debug("skipping non-media file",
			slog.String("path", path),
		)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	media, err := api.PostMedia(
		ctx,
		server,
		album.ID,
		path,
		file.Name(),
		file,
		info.Size(),
	)
	if err != nil {
		if errors.Is(err, api.ErrDuplicateMedia) {
			slog.Warn(
				"media already exists. skipping file.",
				slog.String("path", path),
			)
			return nil
		}
		return fmt.Errorf("failed to upload media for file %s: %w", path, err)
	}

	slog.Info("uploaded media", slog.Any("media", media))

	return nil
}
