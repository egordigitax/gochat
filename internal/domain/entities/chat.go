package entities

import "chat-service/internal/domain/base"

type ChatType int

func (c ChatType) String() string {
	return [...]string{"Personal", "Group"}[c]
}

type Chat struct {
	base.AggregateRoot
	Id          int      `db:"id"`
	Uid         string   `db:"uid"`
	Title       string   `db:"title"`
	MediaURL    string   `db:"media_url"`
	UpdatedAt   string   `db:"updated_at"`
	ChatType    ChatType `db:"chat_type"`
	LastMessage Message  `db:"message"`
}

type ChatRole int

func (c ChatRole) String() string {
	return [...]string{"Personal", "Group"}[c]
}
