package db

type Record struct {
	Id     int     `json:"id"`
	Player string  `json:"player"`
	Score  float32 `json:"score"`
}
