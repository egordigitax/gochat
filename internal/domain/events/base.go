package events

type BrokerBaseAdaptor interface {
	Subscribe(topics ...string) (chan string, error)
	Publish(topic, message string) error
}
