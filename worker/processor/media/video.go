package media

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/gopkg/vips"
	"github.com/wutipong/albums2/worker/util/s3"
)

const VIDEO_WIDTH = 1280
const VIDEO_HEIGHT = 720

func ProcessVideoMedia(ctx context.Context, minioClient *minio.Client, media *types.Media) error {
	slog.Info("process video asset", slog.Any("id", media.ID))

	err := ctx.Err()
	if err != nil {
		slog.Info("context.", slog.String("error", err.Error()))
		return fmt.Errorf("context cancelled: %w", err)
	}

	s3Obj, err := s3.GetObject(
		ctx,
		media.Original,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return fmt.Errorf("unable to get object from s3: %w", err)
	}
	defer s3Obj.Close()

	originalFile, err := os.CreateTemp("",
		fmt.Sprintf("*.%s", filepath.Base(media.Original)),
	)

	if err != nil {
		return fmt.Errorf("unable to create temp file for original file: %w", err)
	}
	defer os.Remove(originalFile.Name())

	io.Copy(originalFile, s3Obj)

	probe, err := ffmpeg.Probe(originalFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var info Probe
	json.Unmarshal([]byte(probe), &info)

	err = processVideoThumbnail(ctx, media, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}

	err = processVideoPreview(ctx, media, originalFile, info)
	if err != nil {
		return fmt.Errorf("unable to process video asset preview: %w", err)
	}

	err = processVideoView(ctx, media, originalFile)
	if err != nil {
		return fmt.Errorf("unable to process video asset view: %w", err)
	}

	return nil
}

func processVideoView(
	ctx context.Context, media *types.Media,
	originalFile *os.File,
) error {
	slog.Info("process video asset view media", slog.Any("id", media.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if media.View == "" || media.View == media.Original {
		media.View = s3.CreateAssetKey("mp4")
	}
	outputFile, err := os.CreateTemp("", "*view.mp4")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	err = ffmpeg.Input(originalFile.Name()).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"vf": fmt.Sprintf(
				"scale=%d:%d:force_original_aspect_ratio=decrease,scale=trunc(iw/2)*2:trunc(ih/2)*2",
				VIDEO_WIDTH, VIDEO_HEIGHT,
			),
			"c:v":      "libx264",
			"preset":   "superfast",
			"crf":      "30",
			"pix_fmt":  "yuv420p",    // Fixes browser incompatibility [6]
			"c:a":      "aac",        // Standard audio
			"movflags": "+faststart", // Enables progressive loading [5]
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create view asset for video asset: %w", err)
	}

	probe, err := ffmpeg.Probe(outputFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	var viewInfo Probe
	json.Unmarshal([]byte(probe), &viewInfo)

	viewVideoStream, err := viewInfo.Video()
	if err != nil {
		return fmt.Errorf("unable to get video stream from video asset: %w", err)
	}
	media.ViewWidth = int32(viewVideoStream.Width)
	media.ViewHeight = int32(viewVideoStream.Height)

	outputFile.Seek(0, io.SeekStart)

	_, err = s3.PutObject(
		ctx,
		media.View,
		outputFile,
		-1,
		minio.PutObjectOptions{
			ContentType: "video/mp4",
		},
	)

	return nil
}

func processVideoThumbnail(
	ctx context.Context, media *types.Media, originalFile *os.File, info Probe,
) error {
	slog.Info("process video asset thumbnail", slog.Any("id", media.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	// save thumbnail at 1/3 duration
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"vframes": "1",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			"vf":      fmt.Sprintf("scale=-2:%d", THUMBNAIL_HEIGHT),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	videoDuration := time.Duration(duration) * time.Second
	media.VideoDuration = int32(videoDuration.Seconds())

	image, err := vips.NewImageFromFile(outputFile.Name(), nil)
	if err != nil {
		return fmt.Errorf("unable to read image file with vips: %w", err)
	}

	media.ThumbnailHeight = THUMBNAIL_HEIGHT
	media.ThumbnailWidth = int32((THUMBNAIL_HEIGHT * image.Width()) / image.Height())

	if media.Thumbnail == "" || media.Thumbnail == media.Original {
		media.Thumbnail = s3.CreateAssetKey("webp")
	}

	_, err = s3.PutObject(
		ctx,
		media.Thumbnail,
		outputFile,
		-1,
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)

	return nil
}

func processVideoPreview(
	ctx context.Context, media *types.Media,
	originalFile *os.File, info Probe,
) error {
	slog.Info("process video preview", slog.Any("id", media.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}
	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer os.Remove(outputFile.Name())

	duration, err := strconv.ParseFloat(info.Format.Duration, 10)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	// save preview at 1/3 duration, 5 seconds-long in 5 fps.
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"t":       "5",
			"loop":    "0",
			"quality": fmt.Sprintf("%d", THUMBNAIL_QUALITY),
			"vf":      fmt.Sprintf("fps=5,scale=-2:%d", THUMBNAIL_HEIGHT),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	if media.Preview == "" || media.Preview == media.Original {
		media.Preview = s3.CreateAssetKey("webp")
	}
	_, err = s3.PutObject(
		ctx,
		media.Preview,
		outputFile,
		-1,
		minio.PutObjectOptions{
			ContentType: "image/webp",
		},
	)
	return nil
}
