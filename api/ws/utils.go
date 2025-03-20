package ws_api

import (
	"chat-service/internal/application/schema/resources"
	"encoding/json"
)

func PackToRootMessage(jsonData json.RawMessage, dto resources.IAction) ([]byte, error) {
	msg := RootMessage{
		ActionType: ActionType(dto.GetActionType()),
		RawPayload: jsonData,
	}

	return json.Marshal(msg)
}
