package main

import (
	"context"
	"encoding/json/v2"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Marlliton/slogpretty"
	"github.com/kouhin/envflag"
	"github.com/redis/go-redis/v9"
	"github.com/wutipong/albums2/worker/processor"
	"github.com/wutipong/albums2/worker/processor/media"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/x/mongo/driver/connstring"
)

var (
	logLevel = flag.String("log-level", "info", "Set log level (debug, info, warn, error)")
	dev      = flag.Bool("dev", false, "Enable development mode")
	dbUri    = flag.String("db-connection", "mongodb://localhost:27017", "MongoDB connection URI")
	redisUrl = flag.String("redis-url", "redis://localhost:6379", "Redis server URL")
	workerId = flag.Int("worker-id", 0, "Worker ID, this must be unique if there are multiple workers")

	ErrDrainingInterrupted = errors.New("draining interrupted")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	envflag.Parse()

	if dev != nil && *dev {
		handler := slogpretty.New(os.Stdout, nil)
		slog.SetDefault(slog.New(handler))
		slog.Debug("Development mode enabled")
	} else {
		handler := slog.NewJSONHandler(os.Stdout, nil)
		slog.SetDefault(slog.New(handler))
	}

	var level slog.Level
	switch *logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		slog.Error("Invalid log level. Using 'info' as default.")
		level = slog.LevelInfo
	}

	slog.SetLogLoggerLevel(level)

	if dbUri == nil {
		slog.Error("database connection string is not set.")
		os.Exit(1)
	}

	cs, err := connstring.Parse(*dbUri)
	if err != nil {
		slog.Error("Error parsing MongoDB connection string", "error", err)
		return
	}

	client, err := mongo.Connect(options.Client().
		ApplyURI(*dbUri))
	if err != nil {
		slog.Error("Failed to connect to MongoDB", "error", err)
		os.Exit(1)
	}

	slog.Info("Connected to MongoDB", "uri", *dbUri)
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			slog.Error("Error disconnecting from MongoDB", "error", err)
		}
	}()

	if redisUrl == nil || *redisUrl == "" {
		slog.Error("Redis URL is not set.")
		os.Exit(1)
	}

	redisOptions, err := redis.ParseURL(*redisUrl)
	if err != nil {
		slog.Error("Failed to parse Redis URL", "error", err)
		os.Exit(1)
	}
	redisClient := redis.NewClient(redisOptions)

	if redisClient == nil {
		slog.Error("Failed to create Redis client")
		os.Exit(1)
	}

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err = redisClient.Close(); err != nil {
			slog.Error("Error disconnecting from Redis", "error", err)
		}
	}()

	slog.Info("Connected to Redis", slog.String("url", *redisUrl))

	processingList := fmt.Sprintf("processing_%d", *workerId)

	slog.Info("Worker", "id", *workerId)

	processor.RegisterProcessor(&media.Processor{
		MongoClient: client,
		Database:    cs.Database,
	})

	slog.Info("Draing worker list")
	err = DrainExistings(ctx, redisClient, processingList)
	if err != nil {
		slog.Error("draining existing items fails", "error", err)
		return
	}

	slog.Info("Start taking new jobs")

	err = ProcessTasks(ctx, redisClient, processingList)
	if err != nil {
		slog.Error("processing taks fails", "error", err)
	}
	slog.Info("Process Loop terminated")
}

func ProcessTasks(ctx context.Context, redisClient *redis.Client, processingList string) error {
	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down gracefully...")
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
					slog.Error("Task process error", "error", err)
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
			slog.Info("Shutting down gracefully")
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
					slog.Error("Task process error", "error", err)
					redisClient.LPush(ctx, "errors", str)
					continue
				}
				redisClient.LPush(ctx, "done", str)
			}
		}
	}
}

func Process(ctx context.Context, redisClient *redis.Client, processingList string, str string) (err error) {
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
