package media

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/gopkg/filetypes"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/gopkg/vips"
	"github.com/wutipong/albums2/worker/util"
)

const MAX_VIEW_PIXEL = 50_000_000
const VIEW_HEIGHT = 2000

func ProcessImageMedia(ctx context.Context, minioClient *minio.Client, media *types.Media) error {
	slog.Info("processing image media", slog.String("id", media.ID.String()))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	slog.Info("getting object from S3.", slog.String("id", media.Original))
	object, err := minioClient.GetObject(
		ctx, os.Getenv("S3_BUCKET"),
		media.Original,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer object.Close()

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

	view, err := populateView(ctx, minioClient, media, original)
	if err != nil {
		return fmt.Errorf("unable to populate view image: %e", err)
	}

	if view == nil {
		view = original
	} else {
		defer view.Close()
	}

	err = populatePreview(ctx, minioClient, media, view)
	if err != nil {
		return fmt.Errorf("unable to populate preview image: %e", err)
	}

	err = populateThumbnail(ctx, minioClient, media, view)
	if err != nil {
		return fmt.Errorf("unable to populate thumbnail: %e", err)
	}

	return nil
}

func populateView(
	ctx context.Context,
	minioClient *minio.Client,
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
		err = fmt.Errorf("Invalid image")
		return
	}

	media.ViewWidth = int32(original.Width())
	media.ViewHeight = int32(original.Height())
	if original.Pages() > 1 {
		media.ViewHeight = int32(original.PageHeight())
	}

	if original.Width()*original.Height() < MAX_VIEW_PIXEL || original.Pages() == 1 {
		media.View = media.Original

		return
	}

	view, err = original.Copy(nil)
	if err != nil {
		err = fmt.Errorf("unable to copy original image: %w", err)
		return
	}

	if view.Width()*view.Height() > MAX_VIEW_PIXEL {
		factor := float64(view.Height()) / float64(VIEW_HEIGHT)

		err = view.Resize(factor, &vips.ResizeOptions{
			Kernel: vips.KernelLanczos3,
			Gap:    2,
		})
		if err != nil {
			err = fmt.Errorf("unable to resize view image: %w", err)
			return
		}

		media.ViewWidth = int32(view.Width())
		media.ViewHeight = int32(view.Height())
	}

	buf, err := view.WebpsaveBuffer(nil)
	if err != nil {
		err = fmt.Errorf("unable to save to webp image: %w", err)
		return
	}

	if media.View == "" || media.View == media.Original {
		media.View = createAssetKey("webp")
	}

	_, err = util.PutObject(
		ctx, minioClient,
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
	minioClient *minio.Client,
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

	if media.ImageFrames == 1 {
		media.Preview = media.View

		return nil
	}

	preview, err := createPreviewForAnimationImage(view)
	if err != nil {
		return err
	}

	defer preview.Close()

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY
	params.PageHeight = preview.PageHeight()

	buf, err := preview.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if media.Preview == "" || media.Preview == media.View {
		media.Preview = createAssetKey("webp")
	}

	_, err = util.PutObject(
		ctx, minioClient,
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

func createPreviewForAnimationImage(original *vips.Image) (*vips.Image, error) {
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	preview, err := original.Copy(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create a preview copy from original image: %w", err)
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(original.PageHeight())

	preview.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	preview.SetPageHeight(THUMBNAIL_HEIGHT)

	slog.Debug("preview image",
		slog.Int("width", preview.Width()),
		slog.Int("height", preview.Height()),
		slog.Int("page_height", preview.PageHeight()),
		slog.Int("loop", preview.Loop()),
		slog.Int("pages", preview.Pages()),
	)
	return preview, nil
}

func populateThumbnail(
	ctx context.Context,
	minioClient *minio.Client,
	media *types.Media,
	view *vips.Image,
) error {
	slog.Info("populating thumbnail media for media", slog.String("id", media.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	media.ThumbnailWidth = int32((view.Width() * THUMBNAIL_HEIGHT) / view.Height())
	media.ThumbnailHeight = THUMBNAIL_HEIGHT

	if view.Pages() == 1 {
		media.Thumbnail = media.View

		return nil
	}

	media.ThumbnailWidth = int32((view.Width() * THUMBNAIL_HEIGHT) / view.PageHeight())

	thumbnail, err := createThumbnailForAnimationImage(view, err)
	if err != nil {
		return err
	}
	defer thumbnail.Close()

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = THUMBNAIL_QUALITY

	buf, err := thumbnail.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if media.Thumbnail == "" || media.Thumbnail == media.Original {
		media.Thumbnail = createAssetKey("webp")
	}

	_, err = util.PutObject(
		ctx, minioClient,
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

func createThumbnailForAnimationImage(original *vips.Image, err error) (*vips.Image, error) {
	slog.Debug("original image",
		slog.Int("width", original.Width()),
		slog.Int("height", original.Height()),
		slog.Int("page_height", original.PageHeight()),
		slog.Int("loop", original.Loop()),
		slog.Int("pages", original.Pages()),
	)

	copyOptions := vips.DefaultCopyOptions()

	thumbnail, _ := original.Copy(copyOptions)

	err = thumbnail.Autorot(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to perform auto rotating: %w", err)
	}

	factor := float64(THUMBNAIL_HEIGHT) / float64(original.PageHeight())
	thumbnail.Resize(factor, &vips.ResizeOptions{
		Kernel: vips.KernelLanczos3,
		Gap:    2,
	})

	err = thumbnail.ExtractArea(0, 0, thumbnail.Width(), THUMBNAIL_HEIGHT)
	if err != nil {
		return nil, fmt.Errorf("unable to extract area: %w", err)
	}
	thumbnail.SetPages(1)
	thumbnail.SetPageHeight(THUMBNAIL_HEIGHT)

	slog.Debug("thumbnail image",
		slog.Int("width", thumbnail.Width()),
		slog.Int("height", thumbnail.Height()),
		slog.Int("page_height", thumbnail.PageHeight()),
		slog.Int("loop", thumbnail.Loop()),
		slog.Int("pages", thumbnail.Pages()),
	)
	return thumbnail, nil
}
