package global

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/model"
	"context"
	"sync"
	"sync/atomic"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ConfigSetting 	jsonModal.Config
	Repository		*Repo
)

var (
    ActiveDevices = make(map[string]model.DeviceStatus) 		// 儲存所有裝置和 狀態
    ActiveDevicesLock   	sync.Mutex          				// 互斥鎖確保併發安全
    GlobalBroker 			mqtt.Client         				// 全域 MQTT 客戶端
	IsConnected				atomic.Bool							// 確認目前連線狀態
)

type Repo struct{
	DB		*DB
	Cache	*Cache
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
