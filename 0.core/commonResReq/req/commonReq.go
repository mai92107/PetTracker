package request

type PageInfo struct{
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	OrderBy   string `json:"orderBy"`
	Direction string `json:"direction"`
}