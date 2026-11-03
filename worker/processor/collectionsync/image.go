package collectionsync

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/filetypes"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/s3"
	"github.com/wutipong/albums2/worker/vips"
)

func calculateThumbnailSize(w, h int32) (width, height int32) {
	width = (w * THUMBNAIL_HEIGHT) / h
	height = THUMBNAIL_HEIGHT

	return
}

func ProcessImageMedia(ctx context.Context,
	media *types.Media,
	reader io.ReadCloser,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("processing image media")

	err := ctx.Err()
	if err != nil {
		zlog.Info().Str("error", err.Error()).Msg("context.")
		return fmt.Errorf("context cancelled: %w", err)
	}

	source := vips.NewSource(reader)
	defer source.Close()

	zlog.Info().Msg("read original image file.")

	params := vips.DefaultLoadOptions()
	if filetypes.HasAnimationExt(filepath.Ext(media.Name)) {
		params.N = -1
	}

	original, err := vips.NewImageFromSource(source, params)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	originalExt := filepath.Ext(media.Original)
	mt := mime.TypeByExtension(originalExt)
	if mt == "image/jpeg" {
		media.Preview = media.Original
		media.Thumbnail = media.Original
		tw, th := calculateThumbnailSize(int32(original.Width()), int32(original.PageHeight()))
		media.ThumbnailWidth = tw
		media.ThumbnailHeight = th

		media.View = media.Original
		media.ViewWidth = int32(original.Width())
		media.ViewHeight = int32(original.PageHeight())

		media.ImageFrames = 1

		return nil
	}

	view, err := populateView(ctx, media, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	if view == nil {
		view = original
	} else {
		defer view.Close()
	}

	err = populatePreview(ctx, media, view)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateThumbnail(ctx, media, source)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	return nil
}

func populateView(
	ctx context.Context,
	media *types.Media,
	original *vips.Image,
) (view *vips.Image, err error) {
	zlog.Info().
		Str("name", media.Name).
		Msg("populating view media for media")

	err = ctx.Err()
	if err != nil {
		err = fmt.Errorf("context cancelled: %w", err)
		return
	}

	if original == nil {
		err = fmt.Errorf("invalid image")
		return
	}

	media.ViewWidth = int32(original.Width())
	media.ViewHeight = int32(original.Height())
	if original.Pages() > 1 {
		media.ViewHeight = int32(original.PageHeight())
	}

	if mime.TypeByExtension(filepath.Ext(media.Original)) == "image/webp" {
		// image is already webp so just reuse the original.
		media.View = media.Original
		return
	}

	view, err = original.Copy(nil)
	if err != nil {
		err = fmt.Errorf("unable to copy original image: %w", err)
		return
	}

	buf, err := view.WebpsaveBuffer(nil)
	if err != nil {
		err = fmt.Errorf("unable to save to webp image: %w", err)
		return
	}

	if media.View == "" || media.View == media.Original {
		media.View = s3.CreateAssetKey(".webp")
	}

	_, err = s3.PutObject(
		ctx,
		media.View,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		err = fmt.Errorf("unable to put object to S3: %w", err)
		return
	}

	return
}

func populatePreview(
	ctx context.Context,
	media *types.Media,
	view *vips.Image,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("populating preview media for media")

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	media.ImageFrames = int32(view.Pages())

	preview, err := view.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to create copy image: %w", err)
	}
	defer preview.Close()

	err = view.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
		Height: THUMBNAIL_HEIGHT,
		Crop:   vips.InterestingNone,
	})

	if err != nil {
		return fmt.Errorf("fails to resize image: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY
	params.PageHeight = preview.PageHeight()

	buf, err := preview.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if media.Preview == "" || media.Preview == media.View {
		media.Preview = s3.CreateAssetKey(".webp")
	}

	_, err = s3.PutObject(
		ctx,
		media.Preview,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		return fmt.Errorf("unable to put preview object to S3: %w", err)
	}

	return nil
}

func populateThumbnail(
	ctx context.Context,
	media *types.Media,
	source *vips.Source,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("populating thumbnail media for media")

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	thumbnail, err := vips.NewImageFromSource(source, &vips.LoadOptions{})

	if err != nil {
		return fmt.Errorf("unable to read image from source: %w", err)
	}

	defer thumbnail.Close()

	err = thumbnail.ThumbnailImage(1_000_000, &vips.ThumbnailImageOptions{
		Height: THUMBNAIL_HEIGHT,
	})

	if err != nil {
		return fmt.Errorf("unable to create thumbnail: %w", err)
	}

	media.ThumbnailHeight = int32(thumbnail.Height())
	media.ThumbnailWidth = int32(thumbnail.Width())

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if media.Thumbnail == "" || media.Thumbnail == media.Original {
		media.Thumbnail = s3.CreateAssetKey(".webp")
	}

	_, err = s3.PutObject(
		ctx,
		media.Thumbnail,
		bytes.NewReader(buf),
		int64(len(buf)),
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	if err != nil {
		return fmt.Errorf("unable to put object to S3: %w", err)
	}

	return nil
}
