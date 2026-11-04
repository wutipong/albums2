package albumsync

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/s3"
	"github.com/wutipong/albums2/worker/util/video"
	"github.com/wutipong/albums2/worker/vips"
)

const VIDEO_WIDTH = 1280
const VIDEO_HEIGHT = 720

func ProcessVideoMedia(ctx context.Context,
	media *types.Media,
	reader io.ReadCloser,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("process video asset")

	err := ctx.Err()
	if err != nil {
		zlog.Info().Str("error", err.Error()).Msg("context.")
		return fmt.Errorf("context cancelled: %w", err)
	}

	originalFile, ok := reader.(*os.File)
	if !ok {
		originalFile, err = os.CreateTemp("",
			fmt.Sprintf("album_archive_*%s", filepath.Ext(media.Name)),
		)

		if err != nil {
			return fmt.Errorf("unable to create temp file for original file: %w", err)
		}

		defer func() {
			err = os.Remove(originalFile.Name())
			if err != nil {
				zlog.Error().Err(err).Msg("error removing temporary source file: %w")
			}
		}()

		_, err = io.Copy(originalFile, reader)
	}

	if err != nil {
		return fmt.Errorf("unable to copy data from S3 to temporary file: %w", err)
	}

	probe, err := video.ReadProbe(originalFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	err = processVideoThumbnail(ctx, media, originalFile, probe)
	if err != nil {
		return fmt.Errorf("unable to process video asset thumbnail: %w", err)
	}

	err = processVideoPreview(ctx, media, originalFile, probe)
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
	zlog.Info().
		Str("name", media.Name).
		Msg("process video asset view media")

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	if media.View == "" || media.View == media.Original {
		media.View = s3.CreateAssetKey(".mp4")
	}
	outputFile, err := os.CreateTemp("", "*view.mp4")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}

	defer func() {
		err = os.Remove(outputFile.Name())
		if err != nil {
			zlog.Error().Err(err).Msg("error deleting temprary source file")
		}
	}()

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

	probe, err := video.ReadProbe(outputFile.Name())
	if err != nil {
		return fmt.Errorf("unable to probe original video: %w", err)
	}

	viewVideoStream, err := probe.Video()
	if err != nil {
		return fmt.Errorf("unable to get video stream from video asset: %w", err)
	}
	media.ViewWidth = int32(viewVideoStream.Width)
	media.ViewHeight = int32(viewVideoStream.Height)

	_, err = outputFile.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("fails to reset file back to the begining: %w", err)
	}

	_, err = s3.PutObject(
		ctx,
		media.View,
		outputFile,
		-1,
		minio.PutObjectOptions{
			ContentType: "video/mp4",
		},
	)

	if err != nil {
		return fmt.Errorf("error putting object to S3: %w", err)
	}

	return nil
}

func processVideoThumbnail(
	ctx context.Context,
	media *types.Media,
	originalFile *os.File,
	info video.Probe,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("process video asset thumbnail")
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	duration, err := strconv.ParseFloat(info.Format.Duration, 32)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}

	defer func() {
		err = os.Remove(outputFile.Name())
		if err != nil {
			zlog.Error().Err(err).Msg("error deleting temporary output file")
		}
	}()

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
	defer image.Close()

	media.ThumbnailHeight = THUMBNAIL_HEIGHT
	media.ThumbnailWidth = int32((THUMBNAIL_HEIGHT * image.Width()) / image.Height())

	if media.Thumbnail == "" || media.Thumbnail == media.Original {
		media.Thumbnail = s3.CreateAssetKey(".webp")
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

	if err != nil {
		return fmt.Errorf("failed to put object to S3: %w", err)
	}

	return nil
}

func processVideoPreview(
	ctx context.Context, media *types.Media,
	originalFile *os.File, info video.Probe,
) error {
	zlog.Info().
		Str("name", media.Name).
		Msg("process video preview")

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}
	outputFile, err := os.CreateTemp("", "*preview.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer func() {
		err = os.Remove(outputFile.Name())
		if err != nil {
			zlog.Error().Err(err).Msg("error deleting temporary output file")
		}
	}()

	duration, err := strconv.ParseFloat(info.Format.Duration, 32)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	err = video.ExtractAnimateScreenshotToFile(
		originalFile.Name(),
		outputFile.Name(),
		duration/3,
		5,
		0,
		THUMBNAIL_QUALITY,
		THUMBNAIL_HEIGHT,
	)

	if err != nil {
		return fmt.Errorf("unable to create thumbnail asset for video asset: %w", err)
	}

	if media.Preview == "" || media.Preview == media.Original {
		media.Preview = s3.CreateAssetKey(".webp")
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

	if err != nil {
		return fmt.Errorf("unable to upload to S3: %w", err)
	}
	return nil
}
