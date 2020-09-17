package main

// Chat is
type Chat interface {
	ConnectAsync(chan<- struct{})
	Disconnect() error
	RegisterHandler(string, func(string) string)
	HandleEventsAsync()
}
