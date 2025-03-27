package types

import (
	"errors"
	"slices"
)

type User struct {
	Uid      string `json:"uid" db:"uid"`
	Nickname string `json:"nickname" db:"nickname"`
	MediaUrl string `json:"media_url" db:"media_url"`
}

func (u User) DeleteMessage(msg Message) error {
	if msg.UserUid != u.Uid {
		return errors.New("not your message")
	}

	return nil
}

func (u User) JoinChat(chat Chat) (ChatUser, error) {
	if !slices.Contains(chat.UsersUids, u.Uid) {
		return ChatUser{}, errors.New("no access")
	}

	return ChatUser{
		Chat: chat,
		User: u,
	}, nil
}
