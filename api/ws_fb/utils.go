package ws_fb

import (
	"chat-service/gen/fbchat"

	flatbuffers "github.com/google/flatbuffers/go"
)

type Serializable interface {
	Pack(builder *flatbuffers.Builder) flatbuffers.UOffsetT
}

func PackRootMessage(
	actionType fbchat.ActionType,
	payload Serializable,
) []byte {

	payloadBuilder := flatbuffers.NewBuilder(256)
	payloadOffset := payload.Pack(payloadBuilder)
	payloadBuilder.Finish(payloadOffset)
	payloadBytes := payloadBuilder.FinishedBytes()

	builder := flatbuffers.NewBuilder(1024)

	rootMessage := fbchat.RootMessageT{
		ActionType: actionType,
		Payload:    payloadBytes,
	}

	offset := rootMessage.Pack(builder)
	builder.Finish(offset)
	return builder.FinishedBytes()
}
