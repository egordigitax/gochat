package resources

import "chat-service/internal/domain/entities"

type Message struct {
	Username  string `json:"username"`
	AuthorUid string `json:"author_uid"`
	ChatUid   string `json:"chat_uid"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

func NewMessageFromEntity(msg entities.Message) Message {
	return Message{
		Username:  msg.UserInfo.Nickname,
		AuthorUid: msg.UserUid,
		ChatUid:   msg.ChatUid,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	}
}

type Action struct {
	Action ActionType  `json:"type"`
	Data   interface{} `json:"data"`
}

type IAction interface {
	GetActionType() ActionType
}

func BuildAction(action IAction) Action {
	return Action{
		Action: action.GetActionType(),
		Data:   action,
	}
}

type ActionType string
