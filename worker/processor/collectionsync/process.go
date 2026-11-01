package collectionsync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	COVER_WIDTH   = 400
	COVER_HEIGHT  = 600
	COVER_QUALITY = 75
)

type Processor struct {
}

func (p *Processor) GetType() string {
	return "collection.sync"
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

	slog.Info("processing collection sync", "collection", payload.ID)

	collectionId, err := bson.ObjectIDFromHex(payload.ID)
	if err != nil {
		return fmt.Errorf("unable to parse collection id: %w", err)
	}

	r := db.MongoDB().
		Collection("collections").
		FindOne(ctx, bson.D{{Key: "_id", Value: collectionId}})

	if r.Err() != nil {
		return fmt.Errorf("unable to get collection information: %w", r.Err())
	}
	collection := types.Collection{}
	err = r.Decode(&collection)
	if err != nil {
		return fmt.Errorf("unable to parse collection information: %w", err)
	}

	slog.Debug("collection information",
		"id", collection.ID.Hex(),
		"name", collection.Name,
		"prefix", collection.Prefix,
	)

	s3Client := s3.Client()
	bucket := s3.Bucket()

	albumsCollection := db.MongoDB().Collection("albums")

	for obj := range s3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix: collection.Prefix,
	}) {
		slog.Info("found object in S3", "key", obj.Key)
		objectKey := obj.Key

		name := strings.TrimPrefix(obj.Key, collection.Prefix)
		if strings.HasPrefix(name, ".") {
			slog.Debug("skipping hidden file", "key", obj.Key)
			continue
		}

		filter := bson.M{"name": name, "collectionId": collection.ID}

		update := bson.M{
			"$set": bson.M{
				"name":         name,
				"collectionId": collection.ID,
				"objectKey":    objectKey,
			},
		}

		opts := options.UpdateOne().SetUpsert(true)

		results, err := albumsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return fmt.Errorf("unable to upsert album: %w", err)
		}

		slog.Debug("upsert result",
			"acknowledged", results.Acknowledged,
			"matchedCount", results.MatchedCount,
			"modifiedCount", results.ModifiedCount,
			"upsertedCount", results.UpsertedCount,
			"upsertedID", results.UpsertedID,
		)

	}

	return nil
}
