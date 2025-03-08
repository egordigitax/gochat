package entities

type User struct {
    Uid string `json:"uid" db:"uid"`
	Nickname string `json:"nickname" db:"nickname"`
	MediaUrl string `json:"media_url" db:"media_url"`
}

type UserUid string

