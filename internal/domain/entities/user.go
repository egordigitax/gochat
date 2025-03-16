package entities

import "errors"

type User struct {
	Uid      string `json:"uid" db:"uid"`
	Nickname string `json:"nickname" db:"nickname"`
	MediaUrl string `json:"media_url" db:"media_url"`
}

func (u *User) DeleteMessage(msg Message) error {
	if msg.UserUid != u.Uid {
		return errors.New("Not your message")
	}

	return nil
}
