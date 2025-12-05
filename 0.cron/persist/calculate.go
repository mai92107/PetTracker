package persist

import (
	"context"
	"fmt"
	"math"
	"time"

	common "batchLog/0.core/commonFunction"
	"batchLog/0.core/global"
	gormTable "batchLog/0.core/gorm"
	"batchLog/0.core/logafa"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type rawData struct {
	DataRef   string      `bson:"data_ref"`
	DeviceID  string      `bson:"device_id"`
	StartTime time.Time   `bson:"start_time"`
	EndTime   time.Time   `bson:"end_time"`
	Coords    [][]float64 `bson:"coords"` // [[lng, lat], [lng, lat], ...]
}

func SaveTripFmMongoToMaria(ctx context.Context) {
	FlushTripFmMongoToMaria(ctx, 30)
}

// 計算近 30min 每趟行程資訊
func FlushTripFmMongoToMaria(ctx context.Context, timeDuration int) {
	coll := global.Repository.DB.MongoDb.Reading.Collection("pettrack")

	duration := time.Now().UTC().Add(time.Minute * -(time.Duration(timeDuration)))
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "recorded_at", Value: bson.D{{Key: "$gte", Value: duration}}},
			{Key: "location", Value: bson.D{{Key: "$ne", Value: nil}}},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "data_ref", Value: 1},
			{Key: "recorded_at", Value: 1},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$data_ref"},
			{Key: "device_id", Value: bson.D{{Key: "$first", Value: "$device_id"}}},
			{Key: "start_time", Value: bson.D{{Key: "$min", Value: "$recorded_at"}}},
			{Key: "end_time", Value: bson.D{{Key: "$max", Value: "$recorded_at"}}},
			{Key: "coords", Value: bson.D{{Key: "$push", Value: "$location.coordinates"}}}, // [lng, lat]
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "data_ref", Value: "$_id"},
			{Key: "device_id", Value: 1},
			{Key: "start_time", Value: 1},
			{Key: "end_time", Value: 1},
			{Key: "coords", Value: 1},
		}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		logafa.Error("Mongo 資料讀取錯誤", "error", err)
		return
	}
	defer cursor.Close(ctx)

	var results []gormTable.TripSummary

	for cursor.Next(ctx) {
		rawData := decodeRawData(cursor)
		distance := getDistance(*rawData)
		duration := rawData.EndTime.Sub(rawData.StartTime).Minutes()
		now := time.Now()

		results = append(results, gormTable.TripSummary{
			DataRef:         rawData.DataRef,
			DeviceID:        rawData.DeviceID,
			StartTime:       common.ToUtcTime(rawData.StartTime),
			EndTime:         common.ToUtcTime(rawData.EndTime),
			DurationMinutes: math.Round(duration*100) / 100, // 保留2位
			PointCount:      len(rawData.Coords),
			DistanceKM:      math.Round(distance*1000) / 1000, // 保留3位
			CreatedAt:       common.ToUtcTime(now),
			UpdatedAt:       common.ToUtcTime(now),
		})
	}

	err = saveTripSummaries(results)
	if err != nil {
		logafa.Error("%w", err)
	}
}

func saveTripSummaries(results []gormTable.TripSummary) error {
	tx := global.Repository.DB.MariaDb.Reading.Begin()
	if err := tx.Error; err != nil {
		return fmt.Errorf("開始交易失敗, %w", err)
	}

	// 確保一定會 rollback（除非我們明確 commit）
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			logafa.Error("交易 panic，已 rollback", "error", r)
		}
	}()

	for i, t := range results {
		if err := saveTripToDB(tx, &t); err != nil {
			logafa.Error("某筆儲存失敗，將 rollback 整批", "no", i+1, "data_ref", t.DataRef, "error", err)
			return fmt.Errorf("儲存失敗: %w", err) // 觸發 rollback
		}
	}

	if err := tx.Commit().Error; err != nil {
		logafa.Error("交易提交失敗", "error", err)
		return fmt.Errorf("commit 失敗: %w", err)
	}

	logafa.Info("全部行程摘要寫入成功！", "count", len(results))
	return nil
}

func saveTripToDB(tx *gorm.DB, trip *gormTable.TripSummary) error {
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "data_ref"}},
		DoUpdates: []clause.Assignment{
			{Column: clause.Column{Name: "updated_at"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(updated_at), updated_at)")},
			{Column: clause.Column{Name: "end_time"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(end_time), end_time)")},
			{Column: clause.Column{Name: "distance_km"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(distance_km), distance_km)")},
			{Column: clause.Column{Name: "duration_minutes"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(duration_minutes), duration_minutes)")},
			{Column: clause.Column{Name: "point_count"}, Value: gorm.Expr("IF(VALUES(point_count) > point_count, VALUES(point_count), point_count)")},
		},
	}).Create(trip).Error
}

func getDistance(rawData rawData) float64 {
	distance := 0.0
	for i := 1; i < len(rawData.Coords); i++ {
		distance += haversine(
			rawData.Coords[i-1][1], rawData.Coords[i-1][0], // lat1, lng1
			rawData.Coords[i][1], rawData.Coords[i][0], // lat2, lng2
		)
	}
	return distance
}

func decodeRawData(cursor *mongo.Cursor) *rawData {
	var temp rawData
	if err := cursor.Decode(&temp); err != nil {
		logafa.Error("decode error", "error", err)
		return nil
	}
	return &temp
}
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371000 // 地球半徑（公尺）
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c / 1000 // 回傳公里
}
