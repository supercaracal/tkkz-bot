package chat

// Client is
type Client interface {
	ConnectAsync(chan<- struct{})
	Disconnect() error
	RegisterHandler(string, func(string) string)
	HandleEventsAsync()
}
