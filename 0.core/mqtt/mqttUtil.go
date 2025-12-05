package mqtt

import (
	"batchLog/0.core/global"
	"batchLog/0.core/logafa"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SubTopic(client mqtt.Client, topic string, handler mqtt.MessageHandler) error {
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
	if topic == "" {
		return fmt.Errorf("發布主題不可為空")
	}
	if msg == "" {
		return fmt.Errorf("發布內容不可為空")
	}
	token := global.GlobalBroker.Publish(topic, 1, false, msg)
	success := token.WaitTimeout(3 * time.Second)
	if !success {
		return fmt.Errorf("發送訊息逾時")
	}
	if token.Error() != nil {
		logafa.Error("❌ 發送訊息失敗", "error", token.Error())
		return token.Error()
	}
	logafa.Debug("✅", "topic", topic, "msg", msg)
	return nil
}
