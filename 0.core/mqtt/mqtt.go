package mqttUtil

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"batchLog/0.core/model"
	router "batchLog/1.router"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
)

// 處理接收到的訊息
func OnMessageReceived(client mqtt.Client, msg mqtt.Message) {
	logafa.Debug("📥 收到 MQTT 訊息！")
	logafa.Debug("主題: %s", msg.Topic())
	logafa.Debug("內容: %s", string(msg.Payload()))

	payloads := []model.MqttPayload{}
	err := jsoniter.Unmarshal(msg.Payload(), &payloads)
	if err != nil{
		logafa.Error("解析mqtt payload錯誤, payload: %s, err: %+v",msg.Payload(),err)
		return
	}
	for _,payload := range payloads{
		router.RouteFunction(msg.Topic(), payload, msg.Qos())
	}
}

func SubTopic(client mqtt.Client,topic string)error{
	token := client.Subscribe(topic, 1, OnMessageReceived)
	if token.Wait() && token.Error() != nil {
		logafa.Error("訂閱失敗: %+v", token.Error())
		return token.Error()
	}
	logafa.Info("已訂閱主題: %s", topic)
	return nil
}

func PubMsgToTopic(msg string, topic string)error{
	// 發送訊息
	token := global.GlobalBroker.Publish(topic, 1, false, msg)
	token.Wait()
	logafa.Debug("已發送訊息: %s\n", msg)
	return nil
}