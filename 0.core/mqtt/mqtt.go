package mqttUtil

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SubTopic(client mqtt.Client, topic string, handler mqtt.MessageHandler) error {
	token := client.Subscribe(topic, 1, handler)
	if token.Wait() && token.Error() != nil {
		logafa.Error("訂閱失敗: %+v", token.Error())
		return token.Error()
	}
	logafa.Info("已訂閱主題: %s", topic)
	return nil
}

func PubMsgToTopic(topic, msg string) error {
	token := global.GlobalBroker.Publish(topic, 1, false, msg)
	token.Wait()
	if token.Error() != nil {
		logafa.Error("❌ 發送訊息失敗: %v", token.Error())
		return token.Error()
	}
	logafa.Debug("✅ topic: %s 已發送訊息: %s", topic, msg)
	return nil
}
