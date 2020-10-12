package chat

// Client is
type Client interface {
	ConnectAsync(chan<- struct{})
	Disconnect() error
	RegisterHandler(string, func(string) string)
	HandleEventsAsync()
}

// Event keys for handler
const (
	EventOnConnection = "onConnection"
	EventOnMessage    = "onMessage"
	EventOnError      = "onError"
)
