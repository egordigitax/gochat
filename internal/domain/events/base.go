package events

import (
	"context"
)

type BrokerBaseAdaptor interface {
	Subscribe(ctx context.Context, topics ...string) (chan string, error)
	Publish(ctx context.Context, topic, message string) error
}
