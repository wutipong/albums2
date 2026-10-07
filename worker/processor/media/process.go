package media

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/wutipong/albums2/gopkg/types"
	"github.com/wutipong/albums2/worker/processor"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Payload struct {
	MediaID string `json:"asset_id"`
}

type Processor struct {
	MongoClient *mongo.Client
	Database    string
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

	result := p.MongoClient.Database(p.Database).
		Collection("media").
		FindOne(ctx, bson.D{{"_id", payload.MediaID}})

	if result.Err() != nil {
		return fmt.Errorf("unable to retrive asset information: %w", result.Err())
	}

	media := types.Media{}
	err = result.Decode(&media)

	if err != nil {
		return fmt.Errorf("unable to decode media information: %w", err)
	}
	return nil
}
