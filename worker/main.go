package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Marlliton/slogpretty"
	"github.com/kouhin/envflag"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	logLevel = flag.String("log-level", "info", "Set log level (debug, info, warn, error)")
	dev      = flag.Bool("dev", false, "Enable development mode")

	dbUri    = flag.String("db-connection", "mongodb://localhost:27017", "MongoDB connection URI")
	redisUrl = flag.String("redis-url", "redis://localhost:6379", "Redis server URL")
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

	if dbUri == nil || *dbUri == "" {
		slog.Error("database connection string is not set.")
		os.Exit(1)
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
	// This is a placeholder for the main function.
	// The actual implementation will depend on the specific requirements of the application.
}
