package media

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/db"
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
	return "media.process"
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

	zlog.Info().Str("id", objID.String()).Msg("processing asset")

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
		UpdateByID(ctx, media.ID, bson.M{
			"$set": bson.M{
				"view":            media.View,
				"viewWidth":       media.ViewWidth,
				"viewHeight":      media.ViewHeight,
				"preview":         media.Preview,
				"thumbnail":       media.Thumbnail,
				"thumbnailWidth":  media.ThumbnailWidth,
				"thumbnailHeight": media.ThumbnailHeight,
				"imageFrames":     media.ImageFrames,
				"videoDuration":   media.VideoDuration,
				"processStatus":   media.ProcessStatus,

				"updatedAt": time.Now(),
			},
		})

	if err != nil {
		return fmt.Errorf("unable to update media information: %w", err)
	}

	if processErr != nil {
		return fmt.Errorf("process media fails: %w", processErr)
	}

	embeddingTask, err := json.Marshal(processor.TaskRequest{
		Type:    "embedding",
		Payload: req.Payload,
	})
	if err != nil {
		return fmt.Errorf("unable to marshal embedding task: %w", err)
	}

	if err = db.Redis().LPush(ctx, "tasks", string(embeddingTask)).Err(); err != nil {
		return fmt.Errorf("unable to queue embedding task: %w", err)
	}

	zlog.Info().Str("id", objID.String()).Msg("process asset complete")

	return nil
}
