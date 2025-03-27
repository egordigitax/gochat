package types

const (
	REQUEST_MESSAGE ActionType = "REQUEST_MESSAGE"
	SEND_MESSAGE    ActionType = "SEND_MESSAGE"
	REQUEST_CHATS   ActionType = "REQUEST_CHATS"
)

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
