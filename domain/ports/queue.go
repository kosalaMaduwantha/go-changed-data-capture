package ports

type QueueSp interface {
	SendMessage(message string, qName string) error
	ReceiveMessage() (string, error)
}