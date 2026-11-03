package collectionsync

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/mholt/archives"
	"github.com/minio/minio-go/v7"
	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/processor/media"
	"github.com/wutipong/albums2/worker/types"
	"github.com/wutipong/albums2/worker/util/charset"
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

	s3Client := s3.Client()
	bucket := s3.Bucket()

	err = populateAlbums(ctx, s3Client, bucket, collection)
	if err != nil {
		return fmt.Errorf("unable to populate albums: %w", err)
	}

	err = populateAlbumsMedia(ctx, collection)
	if err != nil {
		return fmt.Errorf("unable to populate albums media: %w", err)
	}

	err = processMedia(ctx, collection)
	if err != nil {
		return fmt.Errorf("unable to process media: %w", err)
	}

	zlog.Info().
		Str("name", collection.Name).
		Msg("collection sync completed")

	return nil
}

type mediaProcessRequest struct {
	Type    string        `json:"type"`
	Payload media.Payload `json:"payload"`
}

func processMedia(ctx context.Context, collection types.Collection) error {
	cursor, err := db.MongoDB().Collection("media").
		Find(ctx, bson.M{
			"processStatus": "pending",
			"collectionId":  collection.ID,
		})

	if err != nil {
		return fmt.Errorf("unable to get media: %w", err)
	}

	request := []string{}
	// Process the media items
	for cursor.Next(ctx) {
		m := types.Media{}
		err := cursor.Decode(&m)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to decode media information")
			continue
		}

		zlog.Info().
			Str("name", m.Name).
			Str("albumId", m.AlbumID.Hex()).
			Str("collectionId", collection.ID.Hex()).
			Msg("queue processing media job")

		r, err := json.Marshal(mediaProcessRequest{
			Type: "media.process",
			Payload: media.Payload{
				ID: m.ID.Hex(),
			},
		})
		if err != nil {
			zlog.Error().Err(err).Msg("unable to marshal media process request")
			continue
		}
		// Process each media item
		request = append(request, string(r))
	}

	db.Redis().LPush(ctx, "tasks", request)
	return nil
}

func populateAlbums(
	ctx context.Context,
	s3Client *minio.Client,
	bucket string,
	collection types.Collection,
) error {
	albumsCol := db.MongoDB().Collection("albums")

	_, err := albumsCol.UpdateMany(ctx,
		bson.M{"collectionId": collection.ID},
		bson.M{
			"$set": bson.M{
				"existed": false,
			},
		})

	if err != nil {
		return fmt.Errorf("unable to reset existed flags: %w", err)
	}

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
				"existed":   true,
				"updatedAt": time.Now(),
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

func populateAlbumsMedia(ctx context.Context, collection types.Collection) error {
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

		for _, album := range albums {
			zlog.Info().
				Str("name", album.Name).
				Str("collectionId", collection.ID.Hex()).
				Msg("processing album media")

			err = populateAlbumMedia(ctx, collection, album)
			if err != nil {
				zlog.Error().
					Str("name", album.Name).
					Str("collectionId", collection.ID.Hex()).
					Err(err).
					Msg("failed to populate media for album.")
				continue
			}
		}
	}

	return nil
}

func populateAlbumMedia(ctx context.Context, collection types.Collection, album types.Album) error {
	if strings.HasSuffix(album.ObjectKey, "/") {
		return populateDirectoryMedia(ctx, collection, album)

	} else {
		return populateArchiveMedia(ctx, collection, album)
	}
}

func populateDirectoryMedia(
	ctx context.Context,
	collection types.Collection,
	album types.Album,
) error {
	zlog.Info().
		Str("name", album.Name).
		Str("collectionId", collection.ID.Hex()).
		Msg("processing directory album")

	for obj := range s3.Client().ListObjects(ctx, s3.Bucket(), minio.ListObjectsOptions{
		Prefix:    album.ObjectKey,
		Recursive: true,
	}) {
		zlog.Info().Str("key", obj.Key).Msg("found object in S3")
		if _, _, err := archives.Identify(ctx, obj.Key, nil); err == nil {
			zlog.Debug().Str("key", obj.Key).Msg("process archive file")
			if err := processArchiveInDirectory(ctx, obj, collection, album); err != nil {
				zlog.Error().Str("key", obj.Key).Err(err).Msg("unable to process archive object")
			}
			continue
		} else if !IsMediaFile(obj.Key) {
			zlog.Debug().Str("key", obj.Key).Msg("skipping non-media file")
			continue
		}

		zlog.Debug().Str("key", obj.Key).Msg("found media file")

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "objectKey", Value: obj.Key}}, opts)

		if err != nil {
			zlog.Error().Err(err).Msg("unable to count media")
			continue
		}

		if count > 0 {
			zlog.Debug().Str("key", obj.Key).Msg("media already exists, skipping")
			continue
		}

		mt := mime.TypeByExtension(filepath.Ext(obj.Key))
		name := obj.Key[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:          name,
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
			zlog.Error().Err(err).Msg("unable to insert media")
			continue
		}
	}

	return nil
}

func processArchiveInDirectory(
	ctx context.Context,
	obj minio.ObjectInfo,
	collection types.Collection,
	album types.Album,
) error {
	o, err := s3.GetObject(ctx, obj.Key, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("unable to get object from S3: %w", err)
	}
	defer o.Close()

	format, _, err := archives.Identify(ctx, obj.Key, o)
	if err != nil {
		return fmt.Errorf("unable to identify archive: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unable to extract archive: %w", err)
	}

	zlog.Info().Str("format", format.MediaType()).Msg("found archive")

	ex.Extract(ctx, o, func(ctx context.Context, f archives.FileInfo) error {
		zlog.Info().Str("name", f.Name()).Int64("size", f.Size()).Msg("found file in archive")

		decodedName := charset.Decode(f.NameInArchive)
		path := obj.Key + ":" + decodedName

		if _, _, err := archives.Identify(ctx, f.Name(), nil); err == nil {
			zlog.Debug().Str("name", f.Name()).Msg("found archive file in archive")
			if err = processArchiveInArchive(ctx, path, f, collection, album); err != nil {
				zlog.Error().
					Err(err).
					Str("name", decodedName).
					Str("album", album.Name).
					Str("collection", collection.Name).
					Msg("error processing nested archive file.")
			}

			return nil
		}

		if !IsMediaFile(f.Name()) {
			zlog.Debug().Str("name", f.Name()).Msg("skipping non-media file")
			return nil
		}

		zlog.Debug().Str("name", f.Name()).Msg("found media file")

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "path", Value: path}}, opts)

		if err != nil {
			zlog.Error().Err(err).Msg("unable to count media")
			return nil
		}

		if count > 0 {
			zlog.Debug().Str("name", decodedName).Msg("media already exists, skipping")
			return nil
		}

		mt := mime.TypeByExtension(filepath.Ext(decodedName))
		file, err := f.Open()
		if err != nil {
			zlog.Error().Err(err).Msg("unable to open file in archive")
			return nil
		}
		defer file.Close()

		assetKey := s3.CreateAssetKey(filepath.Ext(f.Name()))

		s3.PutObject(ctx, assetKey, file, f.Size(), minio.PutObjectOptions{
			ContentType: mt,
		})

		name := path[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:          name,
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
			zlog.Error().Err(err).Msg("unable to insert media")
			return nil
		}

		return nil
	})

	return nil

}

func processArchiveInArchive(
	ctx context.Context,
	path string,
	f archives.FileInfo,
	collection types.Collection,
	album types.Album,
) error {
	decodedName := charset.Decode(f.NameInArchive)
	zlog.Info().
		Str("name", decodedName).
		Str("album", album.Name).
		Str("collection", collection.Name).
		Msg("processing directory album archive")

	path = strings.Join([]string{path, decodedName}, ":")

	reader, err := f.Open()
	if err != nil {
		zlog.Error().Err(err).Msg("unable to open file in archive")
		return fmt.Errorf("unable to open file in archive: %w", err)
	}
	defer reader.Close()

	file, err := os.CreateTemp("", fmt.Sprintf("album_archive_*%s", filepath.Ext(decodedName)))
	if err != nil {
		zlog.Error().Err(err).Msg("unable to create temporary file")
		return fmt.Errorf("unable to create temporary file: %w", err)
	}
	defer file.Close()

	zlog.Info().Str("name", file.Name()).Str("source", decodedName).Msg("Creating a temp file")

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("unable to copy file to temporary file: %w", err)
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("unable to seek temporary file: %w", err)
	}

	format, _, err := archives.Identify(ctx, f.Name(), file)
	if err != nil {
		return fmt.Errorf("unable to identify archive: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unable to extract archive: %w", err)
	}

	zlog.Info().Str("format", format.MediaType()).Msg("found archive")

	err = ex.Extract(ctx, file, func(ctx context.Context, f archives.FileInfo) error {
		decodedName = charset.Decode(f.NameInArchive)
		zlog.Info().Str("name", decodedName).Int64("size", f.Size()).Msg("found file in archive")

		if _, _, err := archives.Identify(ctx, f.Name(), nil); err == nil {
			zlog.Debug().Str("name", decodedName).Msg("found archive file in archive")
			err = processArchiveInArchive(ctx, path, f, collection, album)
			if err != nil {
				zlog.Error().
					Err(err).
					Str("name", decodedName).
					Str("album", album.Name).
					Str("collection", collection.Name).
					Msg("error processing archive file.")
			}

			return nil
		}

		if !IsMediaFile(f.Name()) {
			zlog.Debug().Str("name", decodedName).Msg("skipping non-media file")
			return nil
		}

		zlog.Debug().Str("name", f.Name()).Msg("found media file")
		objPath := path + ":" + decodedName

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "path", Value: objPath}}, opts)

		if err != nil {
			zlog.Error().Err(err).Msg("unable to count media")
			return nil
		}

		if count > 0 {
			zlog.Debug().Str("name", decodedName).Msg("media already exists, skipping")
			return nil
		}

		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		file, err := f.Open()
		if err != nil {
			zlog.Error().Err(err).Msg("unable to open file in archive")
			return nil
		}
		defer file.Close()

		assetKey := s3.CreateAssetKey(filepath.Ext(f.Name()))

		s3.PutObject(ctx, assetKey, file, f.Size(), minio.PutObjectOptions{
			ContentType: mt,
		})

		name := objPath[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:          name,
			CollectionID:  collection.ID,
			AlbumID:       album.ID,
			Original:      assetKey,
			Path:          objPath,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}
		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to insert media")
			return nil
		}

		return nil
	})

	if err != nil {
		zlog.Error().Err(err).Msg("unable to extract archive")
		return fmt.Errorf("unable to extract archive: %w", err)
	}

	return nil
}

func populateArchiveMedia(
	ctx context.Context,
	collection types.Collection,
	album types.Album,
) error {
	zlog.Info().
		Str("name", album.Name).
		Str("collection", collection.Name).
		Msg("processing archive album")

	obj, err := s3.Client().GetObject(ctx, s3.Bucket(), album.ObjectKey, minio.GetObjectOptions{})
	if err != nil {
		return fmt.Errorf("unable to get object from S3: %w", err)
	}
	defer obj.Close()

	format, _, err := archives.Identify(ctx, album.ObjectKey, obj)
	if err != nil {
		return fmt.Errorf("unable to identify archive: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unable to extract archive: %w", err)
	}

	zlog.Info().Str("format", format.MediaType()).Msg("found archive")

	ex.Extract(ctx, obj, func(ctx context.Context, f archives.FileInfo) error {
		decodedName := charset.Decode(f.NameInArchive)

		zlog.Info().Str("name", decodedName).Int64("size", f.Size()).Msg("found file in archive")

		if _, _, err := archives.Identify(ctx, f.Name(), nil); err == nil {
			zlog.Debug().Str("name", f.Name()).Msg("found archive file in archive")
			return processArchiveInArchive(ctx, album.ObjectKey, f, collection, album)
		}

		if !IsMediaFile(f.Name()) {
			zlog.Debug().Str("name", decodedName).Msg("skipping non-media file")
			return nil
		}

		zlog.Debug().Str("name", decodedName).Msg("found media file")

		path := album.ObjectKey + ":" + decodedName

		opts := options.Count().SetLimit(1)
		count, err := db.MongoDB().Collection("media").
			CountDocuments(ctx, bson.D{{Key: "path", Value: path}}, opts)

		if err != nil {
			zlog.Error().Err(err).Msg("unable to count media")
			return nil
		}

		if count > 0 {
			zlog.Debug().Str("name", decodedName).Msg("media already exists, skipping")
			return nil
		}

		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		file, err := f.Open()
		if err != nil {
			zlog.Error().Err(err).Msg("unable to open file in archive")
			return nil
		}
		defer file.Close()

		assetKey := s3.CreateAssetKey(filepath.Ext(f.Name()))

		s3.PutObject(ctx, assetKey, file, f.Size(), minio.PutObjectOptions{
			ContentType: mt,
		})

		name := path[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:          name,
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
			zlog.Error().Err(err).Msg("unable to insert media")
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
