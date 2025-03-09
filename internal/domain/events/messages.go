package events

type MessageBrokerAdaptor interface {
	Subscribe(topics ...string) (chan string, error)
	Publish(topic, message string) error
}
