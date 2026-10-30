package albumcover

import (
	"bytes"
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
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"github.com/wutipong/albums2/worker/util/video"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	COVER_WIDTH   = 400
	COVER_HEIGHT  = 600
	COVER_QUALITY = 75
)

type Processor struct {
}

func (p *Processor) GetType() string {
	return "albumcover"
}

type Payload struct {
	ID string `json:"id"`
}

func (p *Processor) Process(ctx context.Context, req processor.TaskRequest) error {
	payload := Payload{}
	err := json.Unmarshal(req.Payload, &payload)
	if err != nil {
		return fmt.Errorf("unable to unmarshal the payload: %w", err)
	}

	slog.Debug("Album Cover", "media", payload.ID)

	mediaId, err := bson.ObjectIDFromHex(payload.ID)
	if err != nil {
		return fmt.Errorf("unable to parse media id: %w", err)
	}

	r := db.MongoDB().
		Collection("media").
		FindOne(ctx, bson.D{{Key: "_id", Value: mediaId}})

	if r.Err() != nil {
		return fmt.Errorf("unable to get media information: %w", r.Err())
	}
	media := types.Media{}
	err = r.Decode(&media)
	if err != nil {
		return fmt.Errorf("unable to parse album information: %w", err)
	}

	slog.Info("cover media", "media", media)

	album := types.Album{}
	r = db.MongoDB().
		Collection("albums").
		FindOne(ctx, bson.D{{Key: "_id", Value: media.AlbumID}})

	if r.Err() != nil {
		return fmt.Errorf("unable to get album information: %w", r.Err())
	}
	err = r.Decode(&album)
	if err != nil {
		return fmt.Errorf("unable to parse album information: %w", err)
	}

	slog.Info("target album", "album", album)

	switch media.Type {
	case "image":
		err = p.ProcessImage(ctx, &media, &album)
	case "video":
		err = p.ProcessVideo(ctx, &media, &album)
	}

	if err != nil {
		return fmt.Errorf("failed to process media: %w", err)
	}

	_, err = db.MongoDB().
		Collection("albums").
		UpdateByID(ctx, album.ID, bson.M{
			"$set": bson.M{
				"cover":     album.Cover,
				"updatedAt": time.Now(),
			},
		})

	if err != nil {
		return fmt.Errorf("unable to update media information: %w", err)
	}

	return nil
}

func (p *Processor) ProcessImage(
	ctx context.Context,
	media *types.Media,
	album *types.Album,
) error {
	slog.Info("populating album cover", slog.String("id", media.ID.String()))

	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	object, err := s3.GetObject(
		ctx,
		media.View,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("unable to get source object: %w", err)
	}

	source := vips.NewSource(object)
	defer source.Close()

	slog.Info("read original image file.")

	original, err := vips.NewImageFromSource(source, nil)
	if err != nil {
		return fmt.Errorf("unable to read original image: %w", err)
	}
	defer original.Close()

	options := vips.DefaultThumbnailImageOptions()
	options.Height = COVER_HEIGHT
	options.Crop = vips.InterestingAttention

	err = original.ThumbnailImage(COVER_WIDTH, options)

	if err != nil {
		return fmt.Errorf("unable to create thumbnail: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = COVER_QUALITY

	buf, err := original.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if album.Cover == "" {
		album.Cover = s3.CreateAssetKey("webp")
	}

	_, err = s3.PutObject(
		ctx,
		album.Cover,
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

func (p *Processor) ProcessVideo(
	ctx context.Context, media *types.Media, album *types.Album,
) error {
	slog.Info("process video asset thumbnail", slog.Any("id", media.ID))
	err := ctx.Err()
	if err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	originalFile, err := os.CreateTemp("",
		fmt.Sprintf("*.%s", filepath.Base(media.Original)),
	)

	if err != nil {
		return fmt.Errorf("unable to create temp file for original file: %w", err)
	}

	defer func() {
		err = os.Remove(originalFile.Name())
		if err != nil {
			slog.Error("error deleting temporary source file", "error", err)
		}
	}()

	object, err := s3.GetObject(
		ctx,
		media.View,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("unable to get source object: %w", err)
	}
	defer func() {
		err = object.Close()
		if err != nil {
			slog.Error("error on closing S3 object", "error", err)
		}
	}()

	_, err = io.Copy(originalFile, object)
	if err != nil {
		return fmt.Errorf("fails to create copy image: %w", err)
	}

	probe, err := video.ReadProbe(originalFile.Name())
	if err != nil {
		return fmt.Errorf("unable to read proble: %w", err)
	}

	outputFile, err := os.CreateTemp("", "*view.webp")
	if err != nil {
		return fmt.Errorf("unable to create temp file to transcode: %w", err)
	}
	defer func() {
		err = os.Remove(outputFile.Name())
		if err != nil {
			slog.Error("error deleting temporary output file", "error", err)
		}
	}()

	duration, err := strconv.ParseFloat(probe.Format.Duration, 32)
	if err != nil {
		return fmt.Errorf("unable to parse duration: %w", err)
	}

	// save thumbnail at 1/3 duration
	err = ffmpeg.
		Input(originalFile.Name(), ffmpeg.KwArgs{
			"ss": fmt.Sprintf("%f", duration/3),
		}).
		Output(outputFile.Name(), ffmpeg.KwArgs{
			"c:v":     "libwebp",
			"vframes": "1",
			"quality": fmt.Sprintf("%d", COVER_QUALITY),
		}).OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return fmt.Errorf("unable to create album cover for video asset: %w", err)
	}

	image, err := vips.NewImageFromFile(outputFile.Name(), nil)
	if err != nil {
		return fmt.Errorf("unable to read image file with vips: %w", err)
	}
	defer image.Close()

	options := vips.DefaultThumbnailImageOptions()
	options.Height = COVER_HEIGHT
	options.Crop = vips.InterestingAttention

	err = image.ThumbnailImage(COVER_WIDTH, options)
	if err != nil {
		return fmt.Errorf("fail to resize image: %w", err)
	}

	params := vips.DefaultWebpsaveBufferOptions()
	params.Q = COVER_QUALITY

	buf, err := image.WebpsaveBuffer(params)
	if err != nil {
		return fmt.Errorf("unable to write preview image: %w", err)
	}

	if album.Cover != "" {
		album.Cover = s3.CreateAssetKey("webp")
	}
	_, err = s3.PutObject(
		ctx,
		album.Cover,
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
