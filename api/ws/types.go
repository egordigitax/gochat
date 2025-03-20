package ws_api

import (
	"chat-service/internal/application/use_cases/messages"
	"context"
)

type HandlerFunc func(
	ctx context.Context,
	data interface{},
	client *messages.MessageClient,
) error
