package chats

// ChatClient is
type ChatClient interface {
	ConnectAsync(chan<- struct{})
	Disconnect() error
	RegisterHandler(string, func(string) string)
	HandleEventsAsync()
}
