package model

type User struct {
	ID        string `json:"id"`
	PublicKey string `json:"public_key"`
	CreatedAt int64  `json:"created_at"`
}
