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

// 輸出結構
type TripSummary struct {
	DataRef         string    `bson:"data_ref"`
	DeviceID        string    `bson:"device_id"`
	StartTime       time.Time `bson:"start_time"`
	EndTime         time.Time `bson:"end_time"`
	DurationMinutes float64   `bson:"duration_minutes"`
	PointCount      int       `bson:"point_count"`
	DistanceKM      float64   `bson:"distance_km"`
}

// 單純只撈資料，距離讓 Go 算（最穩最快）
func GetLastDayTripsWithDistance() {
	ctx := context.Background()
	coll := global.Repository.DB.MongoDb.Reading.Collection("pettrack")

	oneDayAgo := time.Now().UTC().Add(-time.Hour * 24)

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
		var temp struct {
			DataRef   string      `bson:"data_ref"`
			DeviceID  string      `bson:"device_id"`
			StartTime time.Time   `bson:"start_time"`
			EndTime   time.Time   `bson:"end_time"`
			Coords    [][]float64 `bson:"coords"` // [[lng, lat], [lng, lat], ...]
		}
		if err := cursor.Decode(&temp); err != nil {
			log.Printf("decode error: %v", err)
			continue
		}

		// Go 端用 Haversine 算距離（超快超穩）
		distance := 0.0
		for i := 1; i < len(temp.Coords); i++ {
			distance += haversine(
				temp.Coords[i-1][1], temp.Coords[i-1][0], // lat1, lng1
				temp.Coords[i][1], temp.Coords[i][0], // lat2, lng2
			)
		}

		results = append(results, TripSummary{
			DataRef:         temp.DataRef,
			DeviceID:        temp.DeviceID,
			StartTime:       temp.StartTime,
			EndTime:         temp.EndTime,
			DurationMinutes: temp.EndTime.Sub(temp.StartTime).Minutes(),
			PointCount:      len(temp.Coords),
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

	for _, t := range results {
		log.Printf("行程 %s | %s | %.3f km | %v 開始 | %v 結束 | 總耗時 %.1f 分鐘",
			t.DataRef, t.DeviceID, t.DistanceKM, t.StartTime, t.EndTime, t.DurationMinutes)
	}
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
