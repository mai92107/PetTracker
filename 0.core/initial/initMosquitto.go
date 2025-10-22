package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	mqttUtil "batchLog/0.core/mqtt"
	router "batchLog/1.router"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 初始化 MQTT 連線
func InitMosquitto(setting jsonModal.MosquittoConfig) mqtt.Client {
	vagueTopic := setting.VagueTopic

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", setting.BrokerHost, setting.BrokerPort)).
		SetClientID(setting.ClientID).
		SetCleanSession(false).
		SetDefaultPublishHandler(router.OnMessageReceived).
		SetConnectionLostHandler(onConnectionLost).
		SetOnConnectHandler(func(c mqtt.Client) {
			logafa.Info("🔄 重新連上 MQTT broker 成功！")
			global.IsConnected.Store(true)
			// 重新訂閱主題
			subscribeVagueTopic(c, vagueTopic)
		})

	client := mqtt.NewClient(opts)

	// 初次連線
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logafa.Error("❌ Mosquitto 初始連線失敗：%v", token.Error())
		panic(token.Error())
	}

	logafa.Info("✅ 已連接到 Mosquitto 伺服器")
	global.IsConnected.Store(true)

	// 啟動監測 Goroutine
	go monitorMqttConnection(client, vagueTopic)

	return client
}

// 監測 MQTT 連線狀態
func monitorMqttConnection(client mqtt.Client, vagueTopic []string) {
	for {
		if !client.IsConnected() {
			logafa.Warn("⚠️ MQTT 已斷線，嘗試重新連線中...")
			token := client.Connect()
			if token.Wait() && token.Error() == nil {
				logafa.Info("🔁 MQTT 重新連線成功！")
				global.IsConnected.Store(true)
				subscribeVagueTopic(client, vagueTopic)
			} else {
				logafa.Warn("❌ 重新連線失敗：%v", token.Error())
			}
		}
		time.Sleep(10 * time.Second)
	}
}

// 訂閱主題
func subscribeVagueTopic(client mqtt.Client, vagueTopic []string) {
	for _, topic := range vagueTopic {
		if err := mqttUtil.SubTopic(client, topic, nil); err != nil {
			logafa.Error(" ❌ 主題:%s, 訂閱失敗： %+v", topic, err)
		} else {
			logafa.Info(" ✅ 系統開始追蹤裝置主題: %s", topic)
		}
	}
}

// 當連線中斷時觸發
func onConnectionLost(client mqtt.Client, err error) {
	fmt.Printf("🚫 Mosquitto 伺服器連線斷開: %v\n", err)
	global.IsConnected.Store(false)
}
