package global

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/model"
	"sync"
	"sync/atomic"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	ConfigSetting jsonModal.Config
	Repository    *model.Repo
)

var (
	ActiveDevices     = make(map[string]model.DeviceStatus) // 儲存所有裝置和 狀態
	ActiveDevicesLock sync.Mutex                            // 互斥鎖確保併發安全
	GlobalBroker      mqtt.Client                           // 全域 MQTT 客戶端
	IsConnected       atomic.Bool                           // 確認目前連線狀態
)

var (
	PriorWorkerPool  chan struct{}
	NormalWorkerPool chan struct{}
)
