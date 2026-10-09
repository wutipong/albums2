package media

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	THUMBNAIL_QUALITY = 75
	THUMBNAIL_HEIGHT  = 600
)

type Payload struct {
	ID string `json:"id"`
}

type Processor struct {
}

func (p *Processor) GetType() string {
	return "media"
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

	slog.Info("processing asset", "id", objID.String())

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

	var processErr error
	switch media.Type {
	case "video":
		processErr = ProcessVideoMedia(ctx, &media)
	case "image":
		processErr = ProcessImageMedia(ctx, &media)
	}

	if processErr == nil {
		media.ProcessStatus = "processed"
	} else {
		media.ProcessStatus = "failed"
	}

	_, err = db.MongoDB().
		Collection("media").
		ReplaceOne(ctx, bson.D{{Key: "_id", Value: objID}}, media)

	if err != nil {
		return fmt.Errorf("unable to update media information: %w", err)
	}

	if processErr != nil {
		return fmt.Errorf("process media fails: %w", processErr)
	}

	if media.Original != media.View &&
		media.Original != media.Preview &&
		media.Original != media.Thumbnail {

		err = s3.RemoveObject(ctx,
			media.Original,
			minio.RemoveObjectOptions{},
		)
	}

	slog.Info("process asset complete", "id", objID.String())

	return nil
}
