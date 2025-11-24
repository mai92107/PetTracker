package model

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Repo struct {
	DB    *DataBase
	Cache *Cache
}
type DataBase struct {
	MariaDb *SqlDB
	MongoDb *NoSqlDB
}
type SqlDB struct {
	Reading *gorm.DB
	Writing *gorm.DB
}
type NoSqlDB struct {
	Reading *mongo.Database
	Writing *mongo.Database
}
type Cache struct {
	Reading *redis.Client
	Writing *redis.Client
	CTX     context.Context
}
