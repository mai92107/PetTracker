package gormTable

import "time"

type GPS struct{
	DeviceCode	string		`json:"device"`
	Longitude	string		`json:"lng"`
	Latitude	string		`json:"lat"`
	RequestTime string		`json:"time"`
}

type DeviceLocation struct {
	UUID        string    `gorm:"type:char(36);primaryKey" json:"uuid"`
	Device      string    `gorm:"type:varchar(32)" json:"device"`
	Lat         string    `gorm:"type:varchar(32)" json:"lat"`
	Lng         string    `gorm:"type:varchar(32)" json:"lng"`
	RecordedAt  string	  `gorm:"column:recorded_at" json:"recorded_at"` // GPS 傳送時間
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"` // DB 寫入時間
}
