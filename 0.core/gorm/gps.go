package gormTable

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GPS struct {
	DeviceId   string  `json:"deviceId"`
	Longitude  float64 `json:"lng"`
	Latitude   float64 `json:"lat"`
	RecordTime time.Time  `json:"time"`
	DataRef    string  `json:"dataRef"`
}

// GeoJSONPoint 地理位置點
type GeoJSONPoint struct {
	Type        string     `bson:"type" json:"type"`               // 固定為 "Point"
	Coordinates [2]float64 `bson:"coordinates" json:"coordinates"` // [經度, 緯度] - 必須是 float64!
}

// DeviceLocation MongoDB 裝置位置記錄
type DeviceLocation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeviceID   string             `bson:"device_id" json:"device_id"`
	Location   GeoJSONPoint       `bson:"location" json:"location"`
	RecordedAt time.Time          `bson:"recorded_at" json:"recorded_at"` // 改用 time.Time
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	DataRef    string             `bson:"data_ref" json:"data_ref"`
}

// NewGeoJSONPoint 建立 GeoJSON Point
func NewGeoJSONPoint(lng, lat float64) GeoJSONPoint {
	return GeoJSONPoint{
		Type:        "Point",
		Coordinates: [2]float64{lng, lat}, // MongoDB GeoJSON 順序: [經度, 緯度]
	}
}
