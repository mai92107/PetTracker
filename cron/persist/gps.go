package persist

import (
	"batchLog/core/global"
	gormTable "batchLog/core/gorm"
	"batchLog/core/logafa"
	"batchLog/core/redis"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func SaveGpsFmRedisToMaria(){
	logafa.Info("開始執行 GPS DATA 持久化...")
	deviceKeyPattern := "device:*"

	keys,err := redis.KeyScan(deviceKeyPattern)
	if err != nil{
		logafa.Error("取得 redis device key 值發生錯誤, error: %+v",err)
	}
	logafa.Debug("取得 %v 筆裝置資料, 開始讀取",len(keys))
	// 取得過去30分中的資料
	end := time.Now().UTC()
	start := end.Add(-30 * time.Minute)

	var records []gormTable.DeviceLocation

	for _, key := range keys {
		datas,err := redis.ZRangeByScore(key,start.UnixMilli(),end.UnixMilli())
		if err != nil{
			logafa.Error("取得 redis device data 發生錯誤, key: %s, error: %+v",key,err)
			continue
		}
		if len(datas) == 0{
			logafa.Debug("讀取到%v筆資料",len(datas))
			continue
		}
		logafa.Debug("準備寫入資料庫...")
	
		device, _ := strings.CutPrefix(key, "device:")
	
		for _, jsonStr := range datas {
			// 解出 lng 及 lat 及 time
			data := gormTable.GPS{}
			jsoniter.UnmarshalFromString(jsonStr,&data)

			record := gormTable.DeviceLocation{}
			record.UUID = uuid.NewString()
			record.Device = device
			record.Lat = data.Latitude
			record.Lng = data.Longitude
			record.RecordedAt = data.RequestTime
			record.CreatedAt = time.Now().UTC()

			records = append(records, record)
		}

		if err = saveToDB(records);err != nil{
			logafa.Error("批次寫入資料至 DB 失敗, error: %+v", err)
			continue
		}
	
		// 清除已寫入資料
		if err := redis.ZRemRangeByScore(key, start.UnixMilli(), end.UnixMilli()); err != nil {
			logafa.Error("刪除 redis 資料失敗 key: %s, error: %+v", key, err)
		}
	}
	
	// 🧠 如果 records 不為空 → 批次寫入

}

func saveToDB(records []gormTable.DeviceLocation)error{
	if len(records) < 1 {
		return fmt.Errorf("無有效紀錄可存入資料庫, 傳入值為 %+v", records)
	}
	logafa.Debug("批次寫入 %d 筆資料至 DB...", len(records))
	if err := global.Repository.DB.Writing.Table("device_location").CreateInBatches(&records,500).Error; err != nil {
		return err
	}
	logafa.Debug("資料成功批次寫入 DB")
	return nil
}