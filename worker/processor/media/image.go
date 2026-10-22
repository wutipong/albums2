package media

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/gopkg/filetypes"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/gopkg/vips"
	"github.com/wutipong/albums2/worker/util/s3"
)

func calculateThumbnailSize(w, h int32) (width, height int32) {
	width = (w * THUMBNAIL_HEIGHT) / h
	height = THUMBNAIL_HEIGHT

	return
}

func ProcessImageMedia(ctx context.Context, media *types.Media) error {
	slog.Info("processing image media", slog.String("id", media.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("getting object from S3.", slog.String("id", media.Original))
	object, err := s3.GetObject(
		ctx,
		media.Original,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}

	defer func() {
		err = object.Close()
		if err != nil {
			slog.Error("error closing S3 object", "error", err)
		}
	}()

	source := vips.NewSource(object)
	defer source.Close()

	slog.Info("read original image file.")

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
	slog.Info("populating view media for media", slog.String("id", media.ID.String()))

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
		media.View = s3.CreateAssetKey("webp")
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
	slog.Info(
		"populating preview media for media",
		slog.String("id", media.ID.String()),
	)

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	media.ImageFrames = int32(view.Pages())

	preview, err := view.Copy(nil)
	if err != nil {
		return fmt.Errorf("unable to create copy image: %w", err)
	}

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
		media.Preview = s3.CreateAssetKey("webp")
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
	slog.Info("populating thumbnail media for media", slog.String("id", media.ID.String()))

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
		media.Thumbnail = s3.CreateAssetKey("webp")
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
