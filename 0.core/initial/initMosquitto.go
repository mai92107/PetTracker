package initial

import (
	jsonModal "batchLog/0.config"
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	mqttUtil "batchLog/0.core/mqtt"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 配置 MQTT 客戶端
func InitMosquitto(setting jsonModal.MosquittoConfig) mqtt.Client {

	vagueTopic := setting.VagueTopic

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%s", setting.BrokerHost, setting.BrokerPort)). // Mosquitto 伺服器地址
		SetClientID(setting.ClientID). // 客戶端 ID
		SetAutoReconnect(true).// 啟用自動重連
		SetDefaultPublishHandler(mqttUtil.OnMessageReceived).
		SetConnectionLostHandler(onConnectionLost).
		SetOnConnectHandler(func (client mqtt.Client)  {
			subscribeVagueTopic(client,vagueTopic)
		})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logafa.Error(" ❌ Mosquitto 連線失敗：%v", token.Error())
		panic(token.Error())
	}
	logafa.Debug(" ✅ 已連接到 Mosquitto 伺服器")
	global.IsConnected.Store(true)

	// 初始訂閱 裝置主題
	subscribeVagueTopic(client,vagueTopic)
	return client
}

func subscribeVagueTopic(client mqtt.Client, vagueTopic string){
	if err := mqttUtil.SubTopic(client,vagueTopic);err != nil {
		logafa.Error(" ❌ 系統裝置追蹤失敗： %+v", err)
	} else {
		logafa.Info(" ✅ 系統開始追蹤裝置主題: %s", vagueTopic)
	}
}

func onConnectionLost(client mqtt.Client, err error) {
	fmt.Printf("mosquitto 伺服器連線斷開: %v\n", err)
	global.IsConnected.Store(false)
}