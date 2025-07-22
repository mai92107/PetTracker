package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/logafa"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 處理接收到的訊息
func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("收到訊息: %s\n主題: %s\n", string(message.Payload()), message.Topic())
}

func InitMosquitto(setting jsonModal.MosquittoConfig)*mqtt.Client{
	// 配置 MQTT 客戶端
	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", setting.BrokerHost,setting.BrokerPort)). // Mosquitto 伺服器地址
		SetClientID(setting.ClientID).    // 客戶端 ID
		SetDefaultPublishHandler(onMessageReceived)

	// 創建 MQTT 客戶端
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logafa.Debug(" ✅ 已連接到 Mosquitto 伺服器")
	return &client
}