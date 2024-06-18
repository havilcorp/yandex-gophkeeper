package entity

// Item модель для хранимой информации
type Item struct {
	ID     int    `json:"id"`
	UserId int    `json:"user_id"`
	Data   []byte `json:"data"`
	Meta   string `json:"meta"`
}
