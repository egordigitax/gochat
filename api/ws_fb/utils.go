package ws_fb

import (
	"chat-service/gen/fbchat"

	flatbuffers "github.com/google/flatbuffers/go"
)

func PackRootMessage(
	actionType fbchat.ActionType,
	payloadType fbchat.RootMessagePayload,
	payload any,
) []byte {

	builder := flatbuffers.NewBuilder(1024)

	rootMessage := fbchat.RootMessageT{
		ActionType: actionType,
		Payload: &fbchat.RootMessagePayloadT{
			Type:  payloadType,
			Value: payload,
		},
	}

    offset := rootMessage.Pack(builder)
    builder.Finish(offset)
	return builder.FinishedBytes()
}
