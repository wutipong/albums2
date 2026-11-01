package collectionsync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mime"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mholt/archives"
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

	err = populateAlbums(ctx, s3Client, bucket, collection)
	if err != nil {
		return fmt.Errorf("unable to populate albums: %w", err)
	}

	err = populateAlbumsMedia(ctx, s3Client, bucket, collection)
	if err != nil {
		return fmt.Errorf("unable to populate albums media: %w", err)
	}

	return nil
}

func populateAlbums(
	ctx context.Context,
	s3Client *minio.Client,
	bucket string,
	collection types.Collection,
) error {
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

		if strings.HasSuffix(name, "/") {
			slog.Debug("found directory", "key", obj.Key)
			name = name[:len(name)-1] // remove trailing slash
		} else if _, _, err := archives.Identify(ctx, obj.Key, nil); err != nil {
			slog.Debug("archive not found, skipping", "key", obj.Key)
			continue
		}

		filter := bson.M{"name": name, "collectionId": collection.ID}

		update := bson.M{
			"$set": bson.M{
				"name":         name,
				"collectionId": collection.ID,
				"objectKey":    objectKey,
				"createdAt":    time.Now(),
				"updatedAt":    time.Now(),
			},
		}

		opts := options.UpdateOne().SetUpsert(true)

		results, err := albumsCollection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			slog.Error("upsert album",
				"name", name,
				"collectionId", collection.ID.Hex(),
				"error", err)
		} else {

			slog.Debug("upsert result",
				"acknowledged", results.Acknowledged,
				"matchedCount", results.MatchedCount,
				"modifiedCount", results.ModifiedCount,
				"upsertedCount", results.UpsertedCount,
				"upsertedID", results.UpsertedID,
			)
		}
	}
	return nil
}

func populateAlbumsMedia(ctx context.Context, s3Client *minio.Client, bucket string, collection types.Collection) error {
	cursor, err := db.MongoDB().Collection("albums").
		Find(ctx, bson.D{{Key: "collectionId", Value: collection.ID}})

	if err != nil {
		return fmt.Errorf("unable to get albums: %w", err)
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		album := types.Album{}
		err := cursor.Decode(&album)
		if err != nil {
			slog.Error("unable to decode album", "error", err)
			continue
		}

		err = processAlbumMedia(ctx, s3Client, bucket, collection, album)
		if err != nil {
			slog.Error("unable to process album media",
				"name", album.Name,
				"collectionId", collection.ID.Hex(),
				"error", err)
			continue
		}
	}
	return nil
}

func processAlbumMedia(ctx context.Context, s3Client *minio.Client, bucket string, collection types.Collection, album types.Album) error {
	slog.Info("processing album media", "name", album.Name, "collectionId", collection.ID.Hex())

	if strings.HasSuffix(album.ObjectKey, "/") {
		processDirectoryAlbum(ctx, s3Client, bucket, collection, album)
		return nil
	} else {
		processArchiveAlbum(ctx, s3Client, bucket, collection, album)
	}
	return nil
}

func processDirectoryAlbum(ctx context.Context, s3Client *minio.Client, bucket string, collection types.Collection, album types.Album) error {
	slog.Info("processing directory album", "name", album.Name, "collectionId", collection.ID.Hex())

	for obj := range s3Client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    album.ObjectKey,
		Recursive: true,
	}) {
		slog.Info("found object in S3", "key", obj.Key)
		if _, _, err := archives.Identify(ctx, obj.Key, nil); err == nil {
			slog.Debug("skipping archive file", "key", obj.Key)
			continue
		} else if IsMediaFile(obj.Key) {
			slog.Debug("found media file", "key", obj.Key)
		} else {
			slog.Debug("skipping non-media file", "key", obj.Key)
			continue
		}

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "objectKey", Value: obj.Key}}, opts)

		if err != nil {
			slog.Error("unable to count media", "error", err)
			continue
		}

		if count > 0 {
			slog.Debug("media already exists, skipping", "key", obj.Key)
			continue
		}

		mt := mime.TypeByExtension(filepath.Ext(obj.Key))

		media := types.Media{
			Name:          obj.Key[len(album.ObjectKey):], // remove album prefix
			CollectionID:  collection.ID,
			AlbumID:       album.ID,
			Original:      obj.Key,
			Path:          obj.Key,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}
		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			slog.Error("unable to insert media", "error", err)
			continue
		}
	}

	return nil
}

func processArchiveAlbum(ctx context.Context, s3Client *minio.Client, bucket string, collection types.Collection, album types.Album) error {
	slog.Info("processing archive album", "name", album.Name, "collectionId", collection.ID.Hex())

	obj, err := s3Client.GetObject(ctx, bucket, album.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		slog.Error("unable to get object from S3", "error", err)
		return fmt.Errorf("unable to get object from S3: %w", err)
	}
	defer obj.Close()

	format, _, err := archives.Identify(ctx, album.ObjectKey, obj)
	if err != nil {
		slog.Error("unable to identify archive", "error", err)
		return fmt.Errorf("unable to identify archive: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		slog.Error("unable to extract archive", "error", err)
		return fmt.Errorf("unable to extract archive: %w", err)
	}

	slog.Info("found archive", "format", format.MediaType())

	ex.Extract(ctx, obj, func(ctx context.Context, f archives.FileInfo) error {
		slog.Info("found file in archive", "name", f.Name(), "size", f.Size())

		if IsMediaFile(f.Name()) {
			slog.Debug("found media file", "name", f.Name())
		} else {
			slog.Debug("skipping non-media file", "name", f.Name())
			return nil
		}

		path := album.ObjectKey + ":" + f.Name()

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "path", Value: path}}, opts)

		if err != nil {
			slog.Error("unable to count media", "error", err)
			return nil
		}

		if count > 0 {
			slog.Debug("media already exists, skipping", "name", f.Name)
			return nil
		}

		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		file, err := f.Open()
		if err != nil {
			slog.Error("unable to open file in archive", "error", err)
			return nil
		}
		defer file.Close()

		assetKey := s3.CreateAssetKey(filepath.Ext(f.Name()))

		slog.Info("putting object to S3", "key", assetKey, "name", f.Name())
		s3.PutObject(ctx, assetKey, file, f.Size(), minio.PutObjectOptions{
			ContentType: mt,
		})

		media := types.Media{
			Name:          f.Name(),
			CollectionID:  collection.ID,
			AlbumID:       album.ID,
			Original:      assetKey,
			Path:          path,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}
		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			slog.Error("unable to insert media", "error", err)
			return nil
		}

		return nil
	})

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
