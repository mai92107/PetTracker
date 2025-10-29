package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/logafa"
	mqttUtil "batchLog/0.core/mqtt"
	router "batchLog/1.router"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	subscriptionMutex sync.Mutex
	subscribedTopics  = make(map[string]bool)
)

// InitMosquitto 初始化 MQTT 連線
func InitMosquitto(setting jsonModal.MosquittoConfig) (mqtt.Client) {
	vagueTopic := setting.VagueTopic

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", setting.BrokerHostCloud, setting.BrokerPort)).
		SetClientID(fmt.Sprintf("%s-%d", setting.ClientID, time.Now().UnixNano())).
		SetUsername(setting.Username).
		SetPassword(setting.Password).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetDefaultPublishHandler(router.OnMessageReceived).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetMaxReconnectInterval(60 * time.Second).
		SetCleanSession(false).
		SetOnConnectHandler(func(c mqtt.Client) {
			logafa.Debug("✅ 已連接到 Mosquitto 伺服器")
			// 使用 goroutine 避免阻塞連線處理
			go subscribeVagueTopic(c, vagueTopic)
		}).
		SetConnectionLostHandler(onConnectionLost).
		SetReconnectingHandler(func(c mqtt.Client, opts *mqtt.ClientOptions) {
			logafa.Info("🔄 正在重新連接到 Mosquitto 伺服器...")
		})

	client := mqtt.NewClient(opts)

	// 初次連線
	logafa.Debug("🔌 正在連接到 MQTT Broker: %s:%s...", setting.BrokerHostLocal, setting.BrokerPort)
	token := client.Connect()
	
	// 等待連線完成,最多 30 秒
	if !token.WaitTimeout(30 * time.Second) {
		 logafa.Error("連線超時")
		 return nil
	}
	
	if token.Error() != nil {
		logafa.Error("❌ Mosquitto 初始連線失敗：%v", token.Error())
		return nil
	}

	logafa.Debug("✅ MQTT 客戶端初始化成功")
	return client
}

// subscribeVagueTopic 訂閱主題(支援重試和去重)
func subscribeVagueTopic(client mqtt.Client, vagueTopic []string) {
	// 等待連線就緒,最多等 10 秒
	for i := 0; i < 100; i++ {
		if client.IsConnected() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !client.IsConnected() {
		logafa.Error("❌ MQTT 未連線,無法訂閱主題")
		return
	}

	subscriptionMutex.Lock()
	defer subscriptionMutex.Unlock()

	for _, topic := range vagueTopic {
		// 檢查是否已訂閱
		if subscribedTopics[topic] {
			logafa.Info("ℹ️  主題 %s 已訂閱,跳過", topic)
			continue
		}

		// 重試機制:最多 3 次
		var err error
		for retry := 0; retry < 3; retry++ {
			if retry > 0 {
				logafa.Debug("🔄 重試訂閱主題 %s (第 %d 次)", topic, retry)
				time.Sleep(time.Second * time.Duration(retry))
			}

			err = mqttUtil.SubTopic(client, topic, nil)
			if err == nil {
				subscribedTopics[topic] = true
				logafa.Debug("✅ 系統開始追蹤裝置主題: %s", topic)
				break
			}

			logafa.Error("⚠️  主題 %s 訂閱失敗(嘗試 %d/3): %v", topic, retry+1, err)
		}

		// 最終失敗
		if err != nil {
			logafa.Error("❌ 主題 %s 訂閱失敗(已重試 3 次): %v", topic, err)
		}
	}
}

// onConnectionLost 當連線中斷時觸發
func onConnectionLost(client mqtt.Client, err error) {
	logafa.Error("🚫 Mosquitto 伺服器連線斷開: %v", err)
	
	// 清空訂閱記錄,重連後需要重新訂閱
	subscriptionMutex.Lock()
	subscribedTopics = make(map[string]bool)
	subscriptionMutex.Unlock()
}