package chat

import (
	"bufio"
	"fmt"
	"os"
)

// LocalClient is
type LocalClient struct {
	handlers map[string]func(string) string
}

// NewLocalClient is
func NewLocalClient() Client {
	return &LocalClient{
		handlers: map[string]func(string) string{},
	}
}

// ConnectAsync is
func (c *LocalClient) ConnectAsync(fail chan<- struct{}) {
}

// Disconnect is
func (c *LocalClient) Disconnect() error {
	return nil
}

// RegisterHandler is
func (c *LocalClient) RegisterHandler(event string, h func(string) string) {
	c.handlers[event] = h
}

// HandleEventsAsync is
func (c *LocalClient) HandleEventsAsync() {
	sc := bufio.NewScanner(os.Stdin)
	go func(sc *bufio.Scanner) {
		for sc.Scan() {
			input := sc.Text()
			if output := c.handleEvent("onMessage", input); output != "" {
				fmt.Println(output)
			}
		}
	}(sc)
}

func (c *LocalClient) handleEvent(eventKey, text string) string {
	if h, ok := c.handlers[eventKey]; ok {
		return h(text)
	}

	return ""
}
