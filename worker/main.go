package main

import (
	"context"
	"encoding/json/v2"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kouhin/envflag"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/processor/albumsync"
	"github.com/wutipong/albums2/worker/processor/collectionsync"
	"github.com/wutipong/albums2/worker/processor/mediaalbumcover"
	"github.com/wutipong/albums2/worker/processor/mediacollectioncover"
	"github.com/wutipong/albums2/worker/processor/mediaembedding"
	"github.com/wutipong/albums2/worker/util"
	"github.com/wutipong/albums2/worker/util/db"
	"github.com/wutipong/albums2/worker/util/s3"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

//go:generate protoc --go_out=. --go-grpc_out=. -I/workspaces/grpc mlservice.proto
//go:generate vipsgen -out ./vips

var (
	appEnv = flag.String(
		"app-env", "production", "application environment (production, develpment)",
	)
	logLevel = flag.String("log-level", "info", "Set log level (debug, info, warn, error)")
	dbUri    = flag.String("db-connection", "mongodb://localhost:27017", "MongoDB connection URI")
	redisUrl = flag.String("redis-url", "redis://localhost:6379", "Redis server URL")
	workerId = flag.Int(
		"worker-id", 0, "Worker ID, this must be unique if there are multiple workers",
	)
	accessKeyId = flag.String("aws-access-key-id", "", "AWS/S3 Access Key ID")
	secret      = flag.String("aws-secret-access-key", "", "AWS/S3 Secret Key")
	awsEndpoint = flag.String("aws-endpoint-url", "", "AWS/S3 Endpoint URL")
	s3Bucket    = flag.String("s3-bucket", "", "S3 Bucket Name")
	awsRegion   = flag.String("aws-default-region", "", "AWS/S3 Region")

	ErrDrainingInterrupted = errors.New("draining interrupted")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := envflag.Parse()
	if err != nil {
		zlog.Error().Err(err).Msg("Error parsing parameters")
		os.Exit(-1)
	}

	var level zerolog.Level
	invalidLogLevel := false
	switch *logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	default:
		invalidLogLevel = true
		level = zerolog.InfoLevel
	}

	dev := *appEnv == "development"

	zlog.Logger = newLogger(os.Stdout, dev, level)
	if invalidLogLevel {
		zlog.Error().Msg("Invalid log level. Using 'info' as default.")
	}
	if dev {
		zlog.Debug().Msg("Development mode enabled")
	}

	if dbUri == nil {
		zlog.Error().Msg("database connection string is not set.")
		os.Exit(1)
	}

	cs, err := connstring.Parse(*dbUri)
	if err != nil {
		zlog.Error().Err(err).Msg("Error parsing MongoDB connection string")
		return
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(*dbUri))
	if err != nil {
		zlog.Error().Err(err).Msg("Failed to connect to MongoDB")
		os.Exit(1)
	}

	zlog.Info().
		Interface("hosts", cs.Hosts).
		Str("database", cs.Database).
		Msg("Connected to MongoDB")
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			zlog.Error().Err(err).Msg("Error disconnecting from MongoDB")
		}
	}()

	if redisUrl == nil || *redisUrl == "" {
		zlog.Error().Msg("Redis URL is not set.")
		os.Exit(1)
	}

	redisOptions, err := redis.ParseURL(*redisUrl)
	if err != nil {
		zlog.Error().Err(err).Msg("Failed to parse Redis URL")
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOptions)

	if redisClient == nil {
		zlog.Error().Msg("Failed to create Redis client")
		os.Exit(1)
	}

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		zlog.Error().Err(err).Msg("Failed to connect to Redis")
		os.Exit(1)
	}

	defer func() {
		if err = redisClient.Close(); err != nil {
			zlog.Error().Err(err).Msg("Error disconnecting from Redis")
		}
	}()

	zlog.Info().Str("url", *redisUrl).Msg("Connected to Redis")

	processingList := fmt.Sprintf("processing_%d", *workerId)

	zlog.Info().Int("id", *workerId).Msg("Worker")

	db.Init(client, cs.Database, redisClient)

	err = processor.RegisterProcessor(&albumsync.Processor{})
	if err != nil {
		zlog.Error().Err(err).Msg("failed to register album sync processor")
		os.Exit(-1)
	}

	err = processor.RegisterProcessor(&collectionsync.Processor{})
	if err != nil {
		zlog.Error().Err(err).Msg("failed to register collection sync processor")
		os.Exit(-1)
	}

	err = processor.RegisterProcessor(&mediaalbumcover.Processor{})
	if err != nil {
		zlog.Error().Err(err).Msg("failed to register album cover processor")
		os.Exit(-1)
	}

	err = processor.RegisterProcessor(&mediacollectioncover.Processor{})
	if err != nil {
		zlog.Error().Err(err).Msg("failed to register collection cover processor")
		os.Exit(-1)
	}

	err = processor.RegisterProcessor(&mediaembedding.Processor{})
	if err != nil {
		zlog.Error().Err(err).Msg("failed to register embedding processor")
		os.Exit(-1)
	}

	util.Init(util.UtilOptions{
		S3Bucket:      *s3Bucket,
		MongoDatabase: cs.Database,
	})

	err = s3.Init(*awsEndpoint, *accessKeyId, *secret, *s3Bucket, *awsRegion)
	if err != nil {
		zlog.Error().Err(err).Msg("S3 client initialization fails")
		return
	}

	zlog.Info().Msg("Draing worker list")
	err = DrainExistings(ctx, redisClient, processingList)
	if err != nil {
		zlog.Error().Err(err).Msg("draining existing items fails")
		return
	}

	zlog.Info().Msg("Start taking new jobs")

	err = ProcessTasks(ctx, redisClient, processingList)
	if err != nil {
		zlog.Error().Err(err).Msg("processing taks fails")
	}
	zlog.Info().Msg("Process Loop terminated")
}

func newLogger(out io.Writer, devMode bool, level zerolog.Level) zerolog.Logger {
	if devMode {
		out = zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	}

	return zerolog.New(out).Level(level).With().Timestamp().Logger()
}

func ProcessTasks(ctx context.Context, redisClient *redis.Client, processingList string) error {
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Msg("Shutting down gracefully...")
			return nil

		default:
			{
				str, err := redisClient.BLMove(
					ctx, "tasks", processingList, "LEFT", "RIGHT", 5*time.Second,
				).Result()
				if err == redis.Nil {
					continue
				}
				if err != nil {
					return fmt.Errorf("processing taks fails: %w", err)
				}

				err = Process(ctx, redisClient, processingList, str)
				if err != nil {
					zlog.Error().Err(err).Msg("Task process error")
					redisClient.LMove(ctx, processingList, "errors", "LEFT", "RIGHT")
					continue
				}
				redisClient.LMove(ctx, processingList, "done", "LEFT", "RIGHT")
			}
		}
	}
}

func DrainExistings(ctx context.Context, redisClient *redis.Client, processingList string) error {
	for {
		select {
		case <-ctx.Done():
			zlog.Info().Msg("Shutting down gracefully")
			return ErrDrainingInterrupted
		default:
			{
				str, err := redisClient.LPop(ctx, processingList).Result()
				if err == redis.Nil {
					return nil
				}

				if err != nil {
					return fmt.Errorf("draining fails: %w", err)
				}

				err = Process(ctx, redisClient, processingList, str)
				if err != nil {
					zlog.Error().Err(err).Msg("Task process error")
					redisClient.RPush(ctx, "errors", str)
					continue
				}
				redisClient.RPush(ctx, "done", str)
			}
		}
	}
}

func Process(
	ctx context.Context,
	redisClient *redis.Client,
	processingList string,
	str string,
) (err error) {
	req := processor.TaskRequest{}
	err = json.Unmarshal([]byte(str), &req)
	if err != nil {
		err = fmt.Errorf("unable to parse task request: %w", err)
		return
	}
	err = processor.Process(ctx, req)
	if err != nil {
		err = fmt.Errorf("task processed with error: %w", err)
		return
	}
	return
}
