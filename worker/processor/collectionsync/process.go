package collectionsync

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mholt/archives"
	"github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/processor/albumsync"
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

	zlog.Info().Str("collection", payload.ID).Msg("processing collection sync")

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

	zlog.Debug().
		Str("id", collection.ID.Hex()).
		Str("name", collection.Name).
		Str("prefix", collection.Prefix).
		Msg("collection details")

	syncAt := time.Now()
	err = populateAlbums(ctx, collection, syncAt)
	if err != nil {
		return fmt.Errorf("unable to populate albums: %w", err)
	}

	err = purgeAlbums(ctx, collection, syncAt)
	if err != nil {
		return fmt.Errorf("failed to purge missing albums: %w", err)
	}

	err = postAlbumsSync(ctx, collection)

	zlog.Info().
		Str("name", collection.Name).
		Msg("collection sync completed")

	return nil
}

func purgeAlbums(
	ctx context.Context,
	collection types.Collection,
	syncAt time.Time,
) error {
	albumsCol := db.MongoDB().Collection("albums")
	pageSize := int64(100)
	skip := int64(0)

	for {
		opts := options.Find().SetLimit(pageSize).SetSkip(skip)
		cursor, err := albumsCol.
			Find(
				ctx,
				bson.D{
					{Key: "collectionId", Value: collection.ID},
					{Key: "lastSyncAt", Value: bson.M{"$ne": syncAt}},
				},
				opts,
			)

		if err != nil {
			return fmt.Errorf("unable to get albums: %w", err)
		}

		var albums []types.Album
		if err := cursor.All(ctx, &albums); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("unable to parse albums: %w", err)
		}
		cursor.Close(ctx)

		if len(albums) == 0 {
			break
		}

		skip = skip + int64(len(albums))

		zlog.Info().
			Str("collectionName", collection.Name).
			Msg("deleting missing album media")

		for _, album := range albums {
			err = deleteMediaCacheObjects(ctx, album)
			if err != nil {
				zlog.Error().Err(err).
					Str("album", album.Name).
					Msg("unable to delete media objects.")
			}
			r, err := db.MongoDB().Collection("media").
				DeleteMany(ctx, bson.D{
					{Key: "albumId", Value: album.ID},
				})
			if err != nil {
				zlog.Error().Err(err).
					Str("album", album.Name).
					Msg("unable to delete album media")
				continue
			}

			zlog.Info().
				Str("album", album.Name).
				Int64("count", r.DeletedCount).
				Bool("acknowledged", r.Acknowledged).
				Msg("album media deleted.")
		}
	}

	r, err := db.MongoDB().
		Collection("albums").
		DeleteMany(ctx, bson.D{
			{Key: "collectionId", Value: collection.ID},
			{Key: "lastSyncAt", Value: bson.M{"$ne": syncAt}},
		})
	if err != nil {
		return fmt.Errorf("failed to delete albums: %w", err)
	}

	zlog.Info().
		Int64("count", r.DeletedCount).
		Bool("acknownledged", r.Acknowledged).
		Msg("albums have been deleted")

	return nil
}

func deleteMediaCacheObjects(
	ctx context.Context,
	album types.Album,
) error {
	pageSize := int64(100)
	skip := int64(0)
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for {
			opts := options.Find().SetLimit(pageSize).SetSkip(skip)
			cursor, err := db.MongoDB().Collection("media").
				Find(
					ctx,
					bson.D{
						{Key: "albumId", Value: album.ID},
					},
					opts,
				)

			if err != nil {
				zlog.Error().Err(err).Msg("unable to get media")
				return
			}

			var medias []types.Media
			if err := cursor.All(ctx, &medias); err != nil {
				cursor.Close(ctx)
				zlog.Error().Err(err).Msg("unable to parse media")

				return
			}
			cursor.Close(ctx)

			if len(medias) == 0 {
				break
			}

			skip = skip + int64(len(medias))

			zlog.Info().
				Str("album", album.Name).
				Msg("deleting missing album media")

			for _, m := range medias {
				if strings.HasPrefix(m.View, "_cache") {
					zlog.Debug().Str("key", m.View).Msg("object to be delete.")
					objectsCh <- minio.ObjectInfo{Key: m.View}
				}

				if strings.HasPrefix(m.Preview, "_cache") {
					zlog.Debug().Str("key", m.Preview).Msg("object to be delete.")
					objectsCh <- minio.ObjectInfo{Key: m.Preview}
				}

				if strings.HasPrefix(m.Thumbnail, "_cache") {
					zlog.Debug().Str("key", m.Thumbnail).Msg("object to be delete.")
					objectsCh <- minio.ObjectInfo{Key: m.Thumbnail}
				}
			}
		}
	}()
	return nil
}

func populateAlbums(
	ctx context.Context,
	collection types.Collection,
	syncAt time.Time,
) error {
	albumsCol := db.MongoDB().Collection("albums")

	s3Client := s3.Client()
	bucket := s3.Bucket()

	for obj := range s3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix: collection.Prefix,
	}) {
		zlog.Info().Str("key", obj.Key).Msg("found object in S3")
		objectKey := obj.Key

		name := strings.TrimPrefix(obj.Key, collection.Prefix)
		if strings.HasPrefix(name, ".") {
			zlog.Debug().Str("key", obj.Key).Msg("skipping hidden file")
			continue
		}

		if strings.HasSuffix(name, "/") {
			zlog.Debug().Str("key", obj.Key).Msg("found directory")
			name = name[:len(name)-1] // remove trailing slash
		} else if _, _, err := archives.Identify(ctx, obj.Key, nil); err != nil {
			zlog.Debug().Str("key", obj.Key).Msg("archive not found, skipping")
			continue
		}

		filter := bson.M{"name": name, "collectionId": collection.ID}

		update := bson.M{
			"$setOnInsert": bson.M{
				"name":         name,
				"collectionId": collection.ID,
				"objectKey":    objectKey,
				"createdAt":    time.Now(),
			},
			"$set": bson.M{
				"updatedAt":  time.Now(),
				"lastSyncAt": syncAt,
			},
		}

		opts := options.UpdateOne().SetUpsert(true)

		results, err := albumsCol.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			zlog.Error().
				Str("name", name).
				Str("collectionId", collection.ID.Hex()).
				Err(err).
				Msg("upsert album")
		} else {
			zlog.Debug().
				Interface("acknowledged", results.Acknowledged).
				Interface("matchedCount", results.MatchedCount).
				Interface("modifiedCount", results.ModifiedCount).
				Interface("upsertedCount", results.UpsertedCount).
				Interface("upsertedID", results.UpsertedID).
				Msg("upsert result")
		}
	}

	//TODO: Create tasks for removing non-existance albums.

	return nil
}

func postAlbumsSync(ctx context.Context, collection types.Collection) error {
	pageSize := int64(100)
	skip := int64(0)

	for {
		opts := options.Find().SetLimit(pageSize).SetSkip(skip)
		cursor, err := db.MongoDB().Collection("albums").
			Find(
				ctx,
				bson.D{{Key: "collectionId", Value: collection.ID}},
				opts,
			)

		if err != nil {
			return fmt.Errorf("unable to get albums: %w", err)
		}

		var albums []types.Album
		if err := cursor.All(ctx, &albums); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("unable to parse albums: %w", err)
		}
		cursor.Close(ctx)

		if len(albums) == 0 {
			break
		}

		skip = skip + int64(len(albums))

		zlog.Info().
			Str("collectionName", collection.Name).
			Msg("posting album sync command")

		tasks := make([]string, 0)

		for _, album := range albums {
			task := processor.PostTaskRequest{
				Type: "album.sync",
				Payload: albumsync.Payload{
					ID: album.ID.Hex(),
				},
			}

			taskStr, err := json.Marshal(task)
			if err != nil {
				continue
			}
			tasks = append(tasks, string(taskStr))
		}
		db.Redis().RPush(ctx, "tasks", tasks)
	}

	return nil
}

var mediaExtensions = []string{
	".jpg",
	".jpeg",
	".png",
	".gif",
	".svg",
	".tiff",
	".webp",

	// movie
	".3gp",
	".avi",
	".m4v",
	".mkv",
	".mov",
	".mp4",
	".webm",
}

func IsMediaFile(path string) bool {
	return slices.Contains(mediaExtensions, strings.ToLower(filepath.Ext(path)))
}
