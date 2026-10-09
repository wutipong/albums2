package db

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	mongoDB     *mongo.Database
	redisClient *redis.Client
)

func Init(m *mongo.Client, db string, r *redis.Client) {
	mongoDB = m.Database(db)
	redisClient = r
}

func MongoDB() *mongo.Database {
	return mongoDB
}

func Redis() *redis.Client {
	return redisClient
}
