package albumsync

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
	"github.com/wutipong/albums2/worker/processor/mediaembedding"
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
	return "album.sync"
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

	zlog.Info().Str("album", payload.ID).Msg("processing album sync")

	albumId, err := bson.ObjectIDFromHex(payload.ID)
	if err != nil {
		return fmt.Errorf("unable to parse collection id: %w", err)
	}

	r := db.MongoDB().
		Collection("albums").
		FindOne(ctx, bson.D{{Key: "_id", Value: albumId}})

	if r.Err() != nil {
		return fmt.Errorf("unable to get album information: %w", r.Err())
	}
	album := types.Album{}
	err = r.Decode(&album)
	if err != nil {
		return fmt.Errorf("unable to parse album information: %w", err)
	}

	collection := types.Collection{}
	r = db.MongoDB().
		Collection("collections").
		FindOne(ctx, bson.D{{Key: "_id", Value: album.CollectionID}})

	if r.Err() != nil {
		return fmt.Errorf("unable to get collection information: %w", r.Err())
	}

	err = r.Decode(&collection)
	if err != nil {
		return fmt.Errorf("unable to parse collection information: %w", err)
	}

	zlog.Debug().
		Str("id", album.ID.Hex()).
		Str("name", album.Name).
		Str("collection", collection.Name).
		Str("prefix", collection.Prefix).
		Msg("collection details")

	err = populateAlbumMedia(ctx, collection, album)
	if err != nil {
		return fmt.Errorf("failed to populate album media: %w", err)
	}

	err = postMediaEmbedding(ctx, album)
	if err != nil {
		return fmt.Errorf("failed to post media embedding populate request: %w", err)
	}

	err = postMediaAlbumCover(ctx, album)
	if err != nil {
		return fmt.Errorf("failed to post media album cover populate request: %w", err)
	}

	zlog.Info().
		Str("name", album.Name).
		Msg("album sync completed")

	return nil
}

func postMediaEmbedding(ctx context.Context, album types.Album) error {
	pageSize := int64(100)
	skip := int64(0)

	for {
		opts := options.Find().SetLimit(pageSize).SetSkip(skip)
		cursor, err := db.MongoDB().Collection("media").
			Find(
				ctx,
				bson.D{{Key: "albumId", Value: album.ID}},
				opts,
			)

		if err != nil {
			return fmt.Errorf("unable to get media: %w", err)
		}

		var medias []types.Media
		if err := cursor.All(ctx, &medias); err != nil {
			cursor.Close(ctx)
			return fmt.Errorf("unable to parse albums: %w", err)
		}
		cursor.Close(ctx)

		if len(medias) == 0 {
			break
		}

		skip = skip + int64(len(medias))

		zlog.Info().
			Str("album", album.Name).
			Msg("posting media embedding generate command")

		tasks := make([]string, 0)

		for _, m := range medias {
			task := processor.PostTaskRequest{
				Type: "album.sync",
				Payload: mediaembedding.Payload{
					ID: m.ID.Hex(),
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

func postMediaAlbumCover(ctx context.Context, album types.Album) error {
	r := db.MongoDB().Collection("media").FindOne(ctx, bson.D{{Key: "albumId", Value: album.ID}})
	if r.Err() != nil {
		return fmt.Errorf("failed to retrieve media data: %w", r.Err())
	}

	m := types.Media{}
	err := r.Decode(&m)
	if err != nil {
		return fmt.Errorf("failed to decode media data")
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
	db.Redis().RPush(ctx, "tasks", taskStr)

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

		data, err := s3.GetObject(ctx, obj.Key, minio.GetObjectOptions{})

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

		errProc := processMediaAsset(ctx, &media, data)
		if errProc != nil {
			media.ProcessStatus = "failed"
		} else {
			media.ProcessStatus = "processed"
		}

		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to insert media")
			continue
		}

		if errProc != nil {
			zlog.Error().Err(errProc).
				Str("name", media.Name).
				Str("album", album.Name).
				Str("collection", collection.Name).
				Msg("failed processing media asset.")
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

		name := path[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:         name,
			CollectionID: collection.ID,
			AlbumID:      album.ID,
			// Original:      assetKey,
			Path:          path,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}

		errProc := processMediaAsset(ctx, &media, file)
		if errProc != nil {
			media.ProcessStatus = "failed"
		} else {
			media.ProcessStatus = "processed"
		}

		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to insert media.")
			return nil
		}

		if errProc != nil {
			return fmt.Errorf("failed to process media asset: %w", errProc)
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

		name := objPath[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:         name,
			CollectionID: collection.ID,
			AlbumID:      album.ID,
			// Original:      assetKey,
			Path:          objPath,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}

		errProc := processMediaAsset(ctx, &media, file)
		if errProc != nil {
			media.ProcessStatus = "failed"
		} else {
			media.ProcessStatus = "processed"
		}

		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to insert media")
			return nil
		}

		if errProc != nil {
			zlog.Error().Err(errProc).
				Str("name", media.Name).
				Str("album", album.Name).
				Str("collection", collection.Name).
				Msg("failed processing media asset.")
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
			zlog.Debug().Str("name", decodedName).
				Msg("media already exists, skipping")
			return nil
		}

		mt := mime.TypeByExtension(filepath.Ext(f.Name()))
		file, err := f.Open()
		if err != nil {
			zlog.Error().Err(err).Msg("unable to open file in archive")
			return nil
		}
		defer file.Close()

		name := path[len(album.ObjectKey):] // remove album prefix
		if name[0] == ':' {
			name = name[1:] // remove leading colon
		}

		media := types.Media{
			Name:         name,
			CollectionID: collection.ID,
			AlbumID:      album.ID,
			// Original:      assetKey,
			Path:          path,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Type:          mt[:strings.Index(mt, "/")],
			ProcessStatus: "pending",
		}

		errProc := processMediaAsset(ctx, &media, file)
		if errProc != nil {
			media.ProcessStatus = "failed"
		} else {
			media.ProcessStatus = "processed"
		}

		_, err = db.MongoDB().Collection("media").InsertOne(ctx, media)
		if err != nil {
			zlog.Error().Err(err).Msg("unable to insert media")
			return nil
		}

		if errProc != nil {
			zlog.Error().Err(errProc).
				Str("name", media.Name).
				Str("album", album.Name).
				Str("collection", collection.Name).
				Msg("failed processing media asset.")
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

func processMediaAsset(ctx context.Context,
	media *types.Media,
	reader io.ReadCloser,
) error {
	switch media.Type {
	case "image":
		return ProcessImageMedia(ctx, media, reader)
	case "video":
		return ProcessVideoMedia(ctx, media, reader)
	default:
		return fmt.Errorf("unable to process media: invalid media type")
	}
}
