package domain

type User struct {
	Nickname string `json:"nickname" db:"nickname"`
	MediaUrl string `json:"media_url" db:"media_url"`
}
