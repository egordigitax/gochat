// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package fbchat

import "strconv"

type ActionType byte

const (
	ActionTypeUNKNOWN      ActionType = 0
	ActionTypeGET_MESSAGE  ActionType = 1
	ActionTypeSEND_MESSAGE ActionType = 2
	ActionTypeGET_CHATS    ActionType = 3
)

var EnumNamesActionType = map[ActionType]string{
	ActionTypeUNKNOWN:      "UNKNOWN",
	ActionTypeGET_MESSAGE:  "GET_MESSAGE",
	ActionTypeSEND_MESSAGE: "SEND_MESSAGE",
	ActionTypeGET_CHATS:    "GET_CHATS",
}

var EnumValuesActionType = map[string]ActionType{
	"UNKNOWN":      ActionTypeUNKNOWN,
	"GET_MESSAGE":  ActionTypeGET_MESSAGE,
	"SEND_MESSAGE": ActionTypeSEND_MESSAGE,
	"GET_CHATS":    ActionTypeGET_CHATS,
}

func (v ActionType) String() string {
	if s, ok := EnumNamesActionType[v]; ok {
		return s
	}
	return "ActionType(" + strconv.FormatInt(int64(v), 10) + ")"
}
