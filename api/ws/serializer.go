package ws_api

import (
	"chat-service/internal/application/schema/dto"
	"chat-service/internal/application/schema/resources"
	"encoding/json"
	"errors"
)


func PackToRootMessage(jsonData []byte, dto resources.IAction) (json.RawMessage, error) {
	msg := RootMessage{
		ActionType: ActionType(dto.GetActionType()),
		RawPayload: jsonData,
	}

	return json.Marshal(msg)
}

func Serialize(action resources.Action) (json.RawMessage, error) {
	switch action.Action {
	case resources.REQUEST_CHATS:
		data, ok := action.Data.(dto.RequestUserChatsPayload)
		if !ok {
			return nil, errors.New("unable serialize RequestUserChatsPayload")
		}

		items := make([]Chat, len(data.Items))
		for i, item := range data.Items {
			items[i] = Chat{
				Title:       item.Title,
				UnreadCount: item.UnreadCount,
				LastMessage: item.LastMessage.Text,
				LastAuthor:  item.LastMessage.Username,
				MediaUrl:    item.MediaUrl,
			}
		}
		res := GetChatsResponse{
			Items: items,
		}
		jsonRes, err := json.Marshal(res)
		if err != nil {
			return nil, err
		}

		return PackToRootMessage(jsonRes, data)

	case resources.REQUEST_MESSAGE:
		data, _ := action.Data.(dto.RequestMessagePayload)

		payload := SendMessageToClientResponse{
			Text:      data.Text,
			AuthorId:  data.AuthorUid,
			Nickname:  "implement me",
			CreatedAt: data.CreatedAt,
		}

		jsonRes, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		return PackToRootMessage(jsonRes, data)
	}

	return nil, nil
}
