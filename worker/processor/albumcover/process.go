package albumcover

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/processor/mediaembedding"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/db"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	zlog "github.com/rs/zerolog/log"
)

type Payload struct {
	ID string `json:"id"`
}

func (p *Processor) GetType() string {
	return "album.cover"
}

type Processor struct{}

func (p *Processor) Process(ctx context.Context, req processor.TaskRequest) error {
	payload := Payload{
		ID: "",
	}
	err := json.Unmarshal(req.Payload, &payload)
	if err != nil {
		return fmt.Errorf("unable to unmarshal the payload: %w", err)
	}

	zlog.Info().Str("album", payload.ID).Msg("processing album sync")

	albumID, err := bson.ObjectIDFromHex(payload.ID)
	if err != nil {
		return fmt.Errorf("unable to parse album ID: %w", err)
	}

	album := types.Album{}
	err = db.MongoDB().Collection("albums").
		FindOne(ctx, bson.D{{Key: "_id", Value: albumID}}).
		Decode(&album)
	if err != nil {
		return fmt.Errorf("unable to retrieve album: %w", err)
	}

	return postMediaAlbumCover(ctx, album)
}

func postMediaAlbumCover(ctx context.Context, album types.Album) error {
	if album.CoverOverriden {
		zlog.Info().Str("album", album.ID.Hex()).Msg("skipping album cover: cover is overridden")
		return nil
	}

	r := db.MongoDB().Collection("settings").FindOne(ctx, bson.D{})
	if r.Err() != nil {
		return fmt.Errorf("failed to retrieve settings: %w", r.Err())
	}

	settings := types.Settings{}
	err := r.Decode(&settings)
	if err != nil {
		return fmt.Errorf("failed to decode settings: %w", err)
	}

	r = db.MongoDB().Collection("media").FindOne(
		ctx,
		bson.D{
			{Key: "albumId", Value: album.ID},
			{Key: "name",
				Value: bson.D{
					{Key: "$nin", Value: settings.Album.ExcludedCoverNames},
				},
			},
		})
	if err := r.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
		return fmt.Errorf("failed to retrieve media data: %w", err)
	}

	m := types.Media{}
	err = r.Decode(&m)
	if err != nil {
		return fmt.Errorf("failed to decode media data: %w", err)
	}

	task := processor.PostTaskRequest{
		Type: "media.album_cover",
		Payload: mediaembedding.Payload{
			ID: m.ID.Hex(),
		},
	}
	taskStr, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	db.Redis().LPush(ctx, "tasks", taskStr)

	return nil
}
