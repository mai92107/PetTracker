package persist

import (
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"
	"batchLog/0.core/redis"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	jsoniter "github.com/json-iterator/go"
)

func SaveGpsFmRedisToMaria() {
	logafa.Info("開始執行 GPS DATA 持久化...")

	deviceKeyPattern := "device:*"
	keys, err := redis.KeyScan(deviceKeyPattern)
	if err != nil {
		logafa.Error("取得 redis device key 值發生錯誤, error: %+v", err)
		return
	}

	if len(keys) == 0 {
		logafa.Debug("取得 %v 筆裝置資料, 罷工回家睡覺", len(keys))
		return
	}

	logafa.Debug("取得 %v 筆裝置資料, 開始讀取", len(keys))

	end := time.Now().UTC()
	start := end.Add(-30 * time.Minute)

	for _, key := range keys {
		datas, err := redis.ZRangeByScore(key, start.UnixMilli(), end.UnixMilli())
		if err != nil {
			logafa.Error("取得 redis device data 發生錯誤, key: %s, error: %+v", key, err)
			continue
		}

		if len(datas) == 0 {
			logafa.Debug("從 key %s 讀取到 0 筆資料", key)
			continue
		}

		logafa.Debug("從 key %s 讀取到 %d 筆資料，準備寫入資料庫...", key, len(datas))

		var records []gormTable.DeviceLocation
		for _, jsonStr := range datas {
			data := gormTable.GPS{}
			if err := jsoniter.UnmarshalFromString(jsonStr, &data); err != nil {
				logafa.Error("解析 GPS JSON 失敗, jsonStr: %s, error: %+v", jsonStr, err)
				continue
			}
			recordTime, _ := time.Parse("2006-01-02 15:04:05", data.RecordTime)
			record := gormTable.DeviceLocation{
				DeviceID:   data.DeviceId,
				Location:   gormTable.NewGeoJSONPoint(data.Longitude, data.Latitude),
				RecordedAt: recordTime,
				DataRef:    data.DataRef,
				CreatedAt:  time.Now().UTC(),
			}
			records = append(records, record)
		}

		if err = saveToDB(records); err != nil {
			logafa.Error("批次寫入資料至 DB 失敗, error: %+v", err)
			continue
		}

		// 只有成功寫入才刪除
		if err := redis.ZRemRangeByScore(
			global.Repository.Cache.Writing,
			key,
			start.UnixMilli(),
			end.UnixMilli(),
		); err != nil {
			logafa.Error("⚠️ 刪除 redis 資料失敗 key: %s, error: %+v", key, err)
			// TODO: 觸發告警或記錄到監控系統
		}
	}
}

func saveToDB(records []gormTable.DeviceLocation) error {
	if len(records) < 1 {
		return fmt.Errorf("無有效紀錄可存入資料庫")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := global.Repository.DB.MongoDb.Writing.
		Collection("pettrack")

	// 使用 BulkWrite 進行 upsert,防止重複資料
	var operations []mongo.WriteModel
	for _, record := range records {
		filter := bson.M{
			"device_id":   record.DeviceID,
			"recorded_at": record.RecordedAt,
		}

		update := bson.M{
			"$setOnInsert": bson.M{
				"device_id":   record.DeviceID,
				"location":    record.Location,
				"recorded_at": record.RecordedAt,
				"data_ref":    record.DataRef,
				"created_at":  record.CreatedAt,
			},
		}

		operation := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, operation)
	}

	logafa.Debug("批次寫入 %d 筆資料至 DB...", len(records))

	result, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		logafa.Error("MongoDB 批次寫入失敗: %+v", err)
		return err
	}

	logafa.Debug("資料成功批次寫入 DB, Upsert: %d 筆", result.UpsertedCount)
	return nil
}
