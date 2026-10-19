package embedding

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/gopkg/vips"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/service"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Payload struct {
	ID string `json:"id"`
}

type Processor struct {
}

func (p *Processor) GetType() string {
	return "embedding"
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

	slog.Info("populating embedding", "id", objID.String())

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

	slog.Info("Media information",
		"id", media.ID,
		"name", media.Name,
		"type", media.Type,
		"original", media.Original,
	)

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
	defer object.Close()

	src := vips.NewSource(object)
	defer src.Close()

	img, err := vips.NewImageFromSource(src, vips.DefaultLoadOptions())
	if err != nil {
		return fmt.Errorf("failed to open thumbnail image: %w", err)
	}

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

	/// TODO: remove object after processed.
	/*
		if media.Original != media.View &&
			media.Original != media.Preview &&
			media.Original != media.Thumbnail {

			err = s3.RemoveObject(ctx,
				media.Original,
				minio.RemoveObjectOptions{},
			)
		}
	*/

	slog.Info("process media embedding complete", "id", objID.String())

	return nil
}
