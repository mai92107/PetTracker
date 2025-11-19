package persist

import (
	"context"
	"log"
	"math"
	"time"

	"batchLog/0.core/global"
	"batchLog/0.core/logafa"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripSummary struct {
	DataRef         string    `bson:"data_ref"`
	DeviceID        string    `bson:"device_id"`
	StartTime       time.Time `bson:"start_time"`
	EndTime         time.Time `bson:"end_time"`
	DurationMinutes float64   `bson:"duration_minutes"`
	PointCount      int       `bson:"point_count"`
	DistanceKM      float64   `bson:"distance_km"`
}

type rawData struct {
	DataRef   string      `bson:"data_ref"`
	DeviceID  string      `bson:"device_id"`
	StartTime time.Time   `bson:"start_time"`
	EndTime   time.Time   `bson:"end_time"`
	Coords    [][]float64 `bson:"coords"` // [[lng, lat], [lng, lat], ...]
}

// 計算近一日每趟行程資訊
func GetLastDayTripsWithDistance() {
	ctx := context.Background()
	coll := global.Repository.DB.MongoDb.Reading.Collection("pettrack")

	oneDayAgo := time.Now().UTC().Add(-time.Hour * 36) // 增加重疊部分 以36小時為基準

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "recorded_at", Value: bson.D{{Key: "$gte", Value: oneDayAgo}}},
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
		logafa.Error("Mongo 資料讀取錯誤, error: %+v", err)
		return
	}
	defer cursor.Close(ctx)

	var results []TripSummary

	for cursor.Next(ctx) {
		rawData := getRawData(cursor)
		distance := getDistance(*rawData)

		results = append(results, TripSummary{
			DataRef:         rawData.DataRef,
			DeviceID:        rawData.DeviceID,
			StartTime:       rawData.StartTime,
			EndTime:         rawData.EndTime,
			DurationMinutes: rawData.EndTime.Sub(rawData.StartTime).Minutes(),
			PointCount:      len(rawData.Coords),
			DistanceKM:      math.Round(distance*1000) / 1000, // 保留3位
		})
	}

	// 依時間排序
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].StartTime.After(results[j].StartTime) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	// 逐項處理
	// TODO: 若有重複DataRef 則更新DB原本資料
	for _, t := range results {
		logafa.Info("行程 %s | %s | %.3f km | %v 開始 | %v 結束 | 總耗時 %.1f 分鐘",
			t.DataRef, t.DeviceID, t.DistanceKM, t.StartTime, t.EndTime, t.DurationMinutes)
	}
}

func getDistance(rawData rawData)float64{
	distance := 0.0
	for i := 1; i < len(rawData.Coords); i++ {
		distance += haversine(
			rawData.Coords[i-1][1], rawData.Coords[i-1][0], // lat1, lng1
			rawData.Coords[i][1], rawData.Coords[i][0], // lat2, lng2
		)
	}
	return distance
}

func getRawData(cursor *mongo.Cursor)*rawData{
	var temp rawData
	if err := cursor.Decode(&temp); err != nil {
		log.Printf("decode error: %v", err)
		return nil
	}
	return &temp
}
// 經典 Haversine 公式（精確到公尺）
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
