package ws

type MessageHandler interface {
	Handle(client *Client, data []byte) error
}
