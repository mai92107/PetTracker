package device

import (
	gormTable "batchLog/core/gorm"
	"batchLog/core/redis"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
)

func Tracking(data gormTable.GPS, now time.Time)error{
	// 存入 redis 臨時保存
	key := fmt.Sprintf("device:%v",data.DeviceCode)
	score := float64(now.UnixMilli())
	byteData, err := jsoniter.Marshal(data)
	if err != nil{
		return err
	}
	// 存入 redis
	redis.ZAddData(key, score, byteData)
	return nil
}