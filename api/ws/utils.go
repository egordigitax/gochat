package ws_api

import (
	"chat-service/internal/types"
	"encoding/json"
)

func PackToRootMessage(jsonData json.RawMessage, dto types.IAction) ([]byte, error) {
	msg := RootMessage{
		ActionType: ActionType(dto.GetActionType()),
		RawPayload: jsonData,
	}

	return json.Marshal(msg)
}
