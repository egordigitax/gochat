package ports

type ClientTransport interface {
	Close() error
	ReadJSON(v interface{}) error
	WriteJSON(v interface{}) error
}
