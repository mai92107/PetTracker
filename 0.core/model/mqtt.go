package model

import "time"

type Subject int

const(
	LOGIN		Subject	 = iota
	LOCATION
	OFFLINE
)

func (s Subject)ToString()string{
	switch s{
	case LOGIN:
		return "LOGIN"
	case LOCATION:
		return "LOCATION"
	case OFFLINE:
		return "OFFLINE"
	default:
		return "UNKNOWN"
	}
	
}

type MqttPayload struct{
	Subject		string		`json:"subject"`
	Payload		string		`json:"payload"`
	CurrentTime	time.Time	`json:"currentTime"`
}

type DeviceStatus struct {
	LastSeen   string	`json:"lastSeen"`
}