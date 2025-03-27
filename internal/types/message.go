package types

import (
	"encoding/json"
	"github.com/google/uuid"
)

type Message struct {
	Uid       string `json:"uid" db:"uid"`
	ChatUid   string `json:"chat_uid" db:"chat_uid"`
	UserUid   string `json:"user_uid" db:"user_uid"`
	Text      string `json:"text" db:"text"`
	CreatedAt string `json:"created_at" db:"created_at"`
	UserInfo  User   `json:"user_info" db:"user_info"`
}

func NewMessage(
	text string,
	userUid string,
	chatUid string,
) Message {
	return Message{
		Uid:     uuid.New().String(),
		ChatUid: chatUid,
		UserUid: userUid,
		Text:    text,
	}
}

func (m Message) ToJSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func NewMessageFromJson(jsonMsg string) (Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(jsonMsg), &msg)
	return msg, err
}
