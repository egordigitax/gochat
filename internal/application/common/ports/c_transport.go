package ports

type ClientTransport interface {
	Close() error
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, data []byte, err error)
}
