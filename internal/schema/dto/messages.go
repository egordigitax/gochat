package dto

import (
	"chat-service/internal/domain/entities"
	"chat-service/internal/schema/resources"
	"errors"

	"github.com/google/uuid"
)

type SendMessageToChatPayload struct {
	MsgType string            `json:"msg_type"`
	Msg     resources.Message `json:"msg"`
}

func (m *SendMessageToChatPayload) Validate() error {
	if m.Msg.ChatUid == "" {
		return errors.New("ChatUid cannot be empty")
	}
	return nil
}

func (m *SendMessageToChatPayload) ToEntity() entities.Message {
	return entities.Message{
		Uid:     uuid.New().String(),
		ChatUid: m.Msg.ChatUid,
		UserUid: m.Msg.AuthorUid,
		Text:    m.Msg.Text,
		UserInfo: entities.User{
			Uid: m.Msg.AuthorUid,
		},
		MessageType: m.MsgType,
	}
}

func BuildGetMessageFromChatPayloadFromEntity(
	m entities.Message,
) GetMessageFromChatPayload {
	return GetMessageFromChatPayload{
		MsgType: m.MessageType,
		Msg: resources.Message{
			Username:  m.UserInfo.Nickname,
			AuthorUid: m.UserUid,
			ChatUid:   m.ChatUid,
			Text:      m.Text,
			CreatedAt: m.CreatedAt,
		},
	}
}

type GetMessageFromChatPayload struct {
	MsgType string            `json:"msg_type"`
	Msg     resources.Message `json:"msg"`
}
