package mediaembedding

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/service"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"github.com/wutipong/albums2/worker/vips"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Payload struct {
	ID string `json:"id"`
}

type Processor struct {
}

func (p *Processor) GetType() string {
	return "media.embedding"
}

func (p *Processor) Process(ctx context.Context, req processor.TaskRequest) error {
	payload := Payload{}
	err := json.Unmarshal(req.Payload, &payload)
	if err != nil {
		return fmt.Errorf("unable to unmarshal the payload: %w", err)
	}

	objID, err := bson.ObjectIDFromHex(payload.ID)

	if err != nil {
		return fmt.Errorf("invalid media id: %w", err)
	}

	zlog.Info().Str("id", objID.String()).Msg("populating embedding")

	result := db.MongoDB().
		Collection("media").
		FindOne(ctx, bson.D{{Key: "_id", Value: objID}})

	if result.Err() != nil {
		return fmt.Errorf("unable to retrive media information: %w", result.Err())
	}

	media := types.Media{}
	err = result.Decode(&media)

	if err != nil {
		return fmt.Errorf("unable to decode media information: %w", err)
	}

	zlog.Info().
		Str("id", media.ID.Hex()).
		Str("name", media.Name).
		Str("type", media.Type).
		Str("original", media.Original).
		Msg("Media information")

	if media.ProcessStatus != "processed" {
		return fmt.Errorf("unable to populate embedding for unprocessed media")
	}

	spec, err := service.GetImageSpec(ctx)
	if err != nil {
		return fmt.Errorf("unable to retrive ML image spec: %w", err)
	}

	object, err := s3.GetObject(ctx, media.Thumbnail, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("unable to get thumbnail object: %w", err)
	}
	defer func() {
		err = object.Close()
		if err != nil {
			zlog.Error().Err(err).Msg("error closing S3 object")
		}
	}()

	src := vips.NewSource(io.NopCloser(object))
	defer src.Close()

	img, err := vips.NewImageFromSource(src, vips.DefaultLoadOptions())
	if err != nil {
		return fmt.Errorf("failed to open thumbnail image: %w", err)
	}
	defer img.Close()

	err = img.ThumbnailImage(int(spec.GetWidth()), &vips.ThumbnailImageOptions{
		Height: int(spec.GetHeight()),
	})

	if err != nil {
		return fmt.Errorf("failed to resize image: %w", err)
	}

	buffer, err := img.WebpsaveBuffer(vips.DefaultWebpsaveBufferOptions())
	if err != nil {
		return fmt.Errorf("unable to create image data: %w", err)
	}

	resp, err := service.EncodeImage(ctx, buffer)
	if err != nil {
		return fmt.Errorf("failed to get embedding data: %w", err)
	}

	_, err = db.MongoDB().
		Collection("media").
		UpdateByID(ctx, media.ID, bson.M{
			"$set": bson.M{
				"embedding": resp.Embedding,
			},
		})

	if err != nil {
		return fmt.Errorf("unable to update media information: %w", err)
	}

	zlog.Info().Str("id", objID.String()).Msg("process media embedding complete")

	return nil
}
