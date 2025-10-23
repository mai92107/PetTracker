package mqttUtil

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SubTopic(client mqtt.Client, topic string, handler mqtt.MessageHandler) error {
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("❌ 訂閱主題時發生恐慌: %v", r)
		}
	}()
	if !global.GlobalBroker.IsConnected() {
		return fmt.Errorf("MQTT broker 未連線")
	}
	if topic == "" {
		return fmt.Errorf("訂閱主題不可為空")
	}
	token := client.Subscribe(topic, 1, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func PubMsgToTopic(topic, msg string) error {
	defer func() {
		if r := recover(); r != nil {
			logafa.Error("❌ 發送主題時發生恐慌: %v", r)
		}
	}()
	if !global.GlobalBroker.IsConnected() {
		return fmt.Errorf("MQTT broker 未連線")
	}
	if topic == "" {
		return fmt.Errorf("發布主題不可為空")
	}
	if msg == "" {
		return fmt.Errorf("發布內容不可為空")
	}
	token := global.GlobalBroker.Publish(topic, 1, false, msg)
	success := token.WaitTimeout(3 * time.Second)
	if !success {
		println("發布後不成功")
		return fmt.Errorf("發送訊息逾時")
	}
	println("發布後成功")
	if token.Error() != nil {
		logafa.Error("❌ 發送訊息失敗: %v", token.Error())
		return token.Error()
	}
	logafa.Debug("✅ topic: %s 已發送訊息: %s", topic, msg)
	return nil
}
