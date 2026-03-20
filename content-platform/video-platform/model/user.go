package model

type User struct {
	ID         string `json:"id"`
	PublicKey  string `json:"public_key"`
	Nickname   string `json:"nickname"`
	AvatarURL  string `json:"avatar_url"`
	Bio        string `json:"bio"`
	CreatedAt  int64  `json:"created_at"`
}
