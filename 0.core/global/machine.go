package global

import (
	jsonModal "batchLog/0.config"
	"context"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ConfigSetting 	jsonModal.Config
	Repository		Repo

	Broker			*mqtt.Client
)

type Repo struct{
	DB		DB
	Cache	Cache
}
type DB struct{
	Reading		*gorm.DB
	Writing		*gorm.DB
}
type Cache struct{
	Reading		*redis.Client
	Writing		*redis.Client
	CTX			context.Context
}
func NewDBRepository(reading, writing *gorm.DB)*DB{
	return &DB{
		Reading: reading,
		Writing: writing,
	}
}
func NewCacheRepository(reading, writing *redis.Client)*Cache{
	return &Cache{
		Reading: reading,
		Writing: writing,
		CTX: context.Background(),
	}
}

