package media

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/worker/processor"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	THUMBNAIL_QUALITY = 75
	THUMBNAIL_HEIGHT  = 600
)

type Payload struct {
	MediaID string `json:"asset_id"`
}

type Processor struct {
	MongoClient *mongo.Client
	Database    string
	MinioClient *minio.Client
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

	slog.Info("processing asset", "id", payload.MediaID)

	objID, err := primitive.ObjectIDFromHex(payload.MediaID)

	if err != nil {
		return fmt.Errorf("invalid media id: %w", err)
	}

	result := p.MongoClient.Database(p.Database).
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
		processErr = ProcessVideoMedia(ctx, p.MinioClient, &media)
	}

	if err == nil {
		media.ProcessStatus = "processed"
	} else {
		media.ProcessStatus = "failed"
	}

	_, err = p.MongoClient.Database(p.Database).
		Collection("media").
		UpdateByID(ctx, objID, media)

	if err != nil {
		return fmt.Errorf("unable to update media information: %w", err)
	}

	if processErr != nil {
		return fmt.Errorf("process media fails: %w", processErr)
	}

	return nil
}

func createAssetKey(extension string) string {
	u, _ := uuid.NewV7()
	return fmt.Sprintf("public/%s.%s", u.String(), extension)
}
