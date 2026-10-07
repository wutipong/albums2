package importing

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mholt/archives"
	"github.com/wutipong/albums2/albums2cli/server/api"
	"github.com/wutipong/albums2/gopkg/types"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func ProcessArchive(
	ctx context.Context,
	server api.ServerConfig,
	album types.Album,
	albumPath string,
) error {
	if !IsArchiveFile(filepath.Ext(albumPath)) {
		return fmt.Errorf("file is not an archive: %s", albumPath)
	}
	archiveFile, err := os.Open(albumPath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %s: %w.",
			albumPath,
			err,
		)
	}
	defer archiveFile.Close()

	err = WalkArchive(ctx, server, album.ID, albumPath, archiveFile)
	if err != nil {
		return fmt.Errorf("failed to process archive %s: %w", albumPath, err)
	}

	return nil
}

func IsArchiveFile(path string) bool {
	return slices.Contains(archiveExtensions, strings.ToLower(filepath.Ext(path)))
}

var archiveExtensions = []string{
	".zip",
	".7z",
	".rar",
}

func WalkArchive(
	ctx context.Context,
	server api.ServerConfig,
	albumID string,
	archivePath string,
	archive io.Reader,
) error {
	if ctx.Err() != nil {
		return fmt.Errorf("context error: %w", ctx.Err())
	}

	format, stream, err := archives.Identify(ctx, archivePath, archive)
	if err != nil {
		return fmt.Errorf("failed to identify archive format: %w", err)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("format does not support extraction")
	}

	var decoder *encoding.Decoder = nil

	err = extractor.Extract(
		ctx, stream,
		func(ctx context.Context, f archives.FileInfo) error {
			if ctx.Err() != nil {
				err = fmt.Errorf("context error: %w", ctx.Err())
				return err
			}

			if f.IsDir() {
				return nil
			}

			filename := f.NameInArchive
			if decoder == nil {
				en, charset, err := DetectCharSet(filename)
				if err != nil {
					slog.Warn("failed to detect filename character set",
						slog.String("filename", filename),
						slog.String("archive file", archivePath))

					decoder = &encoding.Decoder{
						Transformer: transform.Nop,
					}
				} else {
					slog.Debug(
						"using character set",
						slog.String("charset", charset),
						slog.String("filename", filename),
						slog.String("archive file", archivePath),
					)
					decoder = en.NewDecoder()
				}
			}

			if decoder != nil {
				filename, err = decoder.String(filename)
				if err != nil {
					slog.Warn(
						"failed to decode filename. using filename as is.",
						slog.String("filename", filename),
						slog.String("error", err.Error()),
					)
				}
			}

			if IsArchiveFile(f.NameInArchive) {
				slog.Info(
					"nested archive found.",
					slog.String("filename", filename),
					slog.String("archive", archivePath),
				)

				file, err := f.Open()

				if err != nil {
					return fmt.Errorf("failed to open nested archive %s: %w", filename, err)
				}
				defer file.Close()

				tempFile, err := os.CreateTemp("", filepath.Base(filename))
				if err != nil {
					return fmt.Errorf("unable to create temporary file: %w", err)
				}

				defer os.Remove(tempFile.Name())
				defer tempFile.Close()

				_, err = io.Copy(tempFile, file)
				if err != nil {
					return fmt.Errorf("unable to copy data to the temporary: %w", err)
				}
				tempFile.Seek(0, io.SeekStart)

				err = WalkArchive(
					ctx, server, albumID, filepath.Join(archivePath, filename), tempFile,
				)
				if err != nil {
					return fmt.Errorf("failed to process nested archive %s: %w", filename, err)
				}

				return nil
			}

			if IsMediaFile(f.NameInArchive) {
				media, err := uploadArchiveMedia(ctx, server, albumID, archivePath, filename, f)
				if errors.Is(err, api.ErrDuplicateMedia) {
					slog.Warn(
						"media already exists. skipping file.",
						slog.String("path", archivePath),
					)
					return nil
				}
				if err != nil {
					return fmt.Errorf(
						"failed to upload media for file %s in archive %s: %w",
						filename, archivePath, err,
					)
				}

				slog.Info("uploaded media", slog.Any("media", media))
			}

			return nil
		})

	if err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	return nil
}

func uploadArchiveMedia(
	ctx context.Context,
	server api.ServerConfig,
	albumID string,
	archivePath string,
	filename string,
	f archives.FileInfo,
) (media types.Media, err error) {
	if ctx.Err() != nil {
		err = fmt.Errorf("context error: %w", ctx.Err())
		return
	}

	slog.Info(
		"creating media",
		slog.String("archive", archivePath),
		slog.String("entry", filename),
	)

	file, err := f.Open()
	if err != nil {
		err = fmt.Errorf("failed to open archive entry %s/%s: %w", archivePath, filename, err)
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		err = fmt.Errorf("unable to retrieve file stat:%w", err)
		return
	}

	resp, err := api.PostMedia(
		ctx,
		server,
		albumID,
		archivePath,
		filename,
		file,
		stat.Size(),
	)
	if err != nil {
		err = fmt.Errorf("failed to upload media %s/%s: %w", archivePath, filename, err)
		return
	}
	media = resp.Media
	return
}
