package entities

import (
	"encoding/json"
	"errors"
)

type Message struct {
	Uid         string `json:"uid" db:"uid"`
	ChatUid     string `json:"chat_uid" db:"chat_uid"`
	UserUid     string `json:"user_uid" db:"user_uid"`
	Text        string `json:"text" db:"message"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	UserInfo    User   `json:"user_info" db:"user_info"`
	MessageType string `json:"message_type" db:"message_type"`
}

func (m *Message) ToJSON() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func NewMessageFromJson(jsonMsg string) (Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(jsonMsg), &msg)
	return msg, err
}

func (m *Message) Validate() error {
    if m.ChatUid == "" {
        return errors.New("ChatUid not found")
    }
    
    if m.UserUid == "" {
        return errors.New("UserUid not found")
    }
    
    if m.Text == "" {
        return errors.New("Text not found")
    }

    return nil
}
