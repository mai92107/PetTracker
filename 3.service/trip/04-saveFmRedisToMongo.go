package tripService

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

// 持久化 過去時間(分鐘) 定位資料
func FlushGpsFmRdsToMongo(ctx context.Context, deviceId *string, duration int) {
	logafa.Info("開始執行 GPS DATA 持久化...")
	deviceKeys := []string{}
	if deviceId == nil {
		deviceKeyPattern := "device:*"
		keys, err := redis.KeyScan(ctx, deviceKeyPattern)
		if err != nil {
			logafa.Error("取得 redis device key 值發生錯誤", "error", err)
			return
		}

		if len(keys) == 0 {
			logafa.Debug("取無裝置資料, 罷工回家睡覺")
			return
		}
		logafa.Debug("取得裝置資料, 開始讀取", "count", len(keys))
		deviceKeys = append(deviceKeys, keys...)
	} else {
		deviceKeys = []string{fmt.Sprintf("device:%s", *deviceId)}
	}

	end := time.Now().UTC()
	start := end.Add(time.Minute * time.Duration(-1*duration))

	var records []gormTable.DeviceLocation
	for _, key := range deviceKeys {
		readAndOrganizeKeyData(ctx, key, start.UnixMilli(), end.UnixMilli(), &records)
	}

	batchSize := 5000
	for i := 0; i < len(records); i += batchSize {
		end := i + batchSize
		if end > len(records) {
			end = len(records)
		}
		if err := saveLocationToMongoDB(ctx, records[i:end]); err != nil {
			logafa.Error("批次寫入失敗", "from", i, "to", end, "error", err)
			continue
		}
	}

	// 只有成功寫入才刪除
	if err := redis.ZRemRangeByScore(ctx, deviceKeys, start.UnixMilli(), end.UnixMilli()); err != nil {
		logafa.Error("⚠️ 刪除 redis 資料失敗", "error", err)
		// TODO: 觸發告警或記錄到監控系統
	}
}

func readAndOrganizeKeyData(ctx context.Context, key string, startTs int64, endTs int64, records *[]gormTable.DeviceLocation) {
	datas, err := redis.ZRangeByScore(ctx, key, startTs, endTs)
	if err != nil {
		logafa.Error("取得 redis device data 發生錯誤", "key", key, "error", err)
		return
	}

	if len(datas) == 0 {
		logafa.Debug("從讀取無資料", "key", key)
		return
	}

	logafa.Debug("準備寫入資料庫...", "key", key, "count", len(datas))

	for _, jsonStr := range datas {
		data := gormTable.GPS{}
		if err := jsoniter.UnmarshalFromString(jsonStr, &data); err != nil {
			logafa.Error("解析 GPS JSON 失敗", "jsonStr", jsonStr, "error", err)
			continue
		}
		record := gormTable.DeviceLocation{
			DeviceID:   data.DeviceId,
			Location:   gormTable.NewGeoJSONPoint(data.Longitude, data.Latitude),
			RecordedAt: data.RecordTime,
			DataRef:    data.DataRef,
			CreatedAt:  time.Now().UTC(),
		}
		*records = append(*records, record)
	}
}

func saveLocationToMongoDB(ctx context.Context, records []gormTable.DeviceLocation) error {
	start := global.GetNow()
	if len(records) < 1 {
		return fmt.Errorf("無有效紀錄可存入資料庫")
	}

	seen := make(map[string]bool)
	deviceIds := []string{}
	for _, r := range records {
		if !seen[r.DeviceID] {
			seen[r.DeviceID] = true
			deviceIds = append(deviceIds, r.DeviceID)
		}
	}

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

	logafa.Debug("批次寫入DB...", "count", len(records))

	result, err := collection.BulkWrite(ctx, operations)
	if err != nil {
		logafa.Error("MongoDB 批次寫入失敗", "error", err)
		return err
	}

	logafa.Info("GPS 持久化完成",
		"device_count", len(deviceIds),
		"records", len(records),
		"upserted", result.UpsertedCount,
		"duration_minutes", time.Since(start).Round(time.Second))
	return nil
}
