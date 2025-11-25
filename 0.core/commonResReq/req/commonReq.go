package request

type MqttReq struct{
	SubscribeTo string `json:"subscribeTo"`
}
type PageInfo struct{
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	OrderBy   string `json:"orderBy"`
	Direction string `json:"direction"`
}