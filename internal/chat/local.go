package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		fmt.Print("> ")

		for sc.Scan() {
			input := sc.Text()
			lowerInput := strings.ToLower(input)
			if lowerInput == "exit" || lowerInput == "quit" {
				break
			}

			if output := c.handleEvent(EventOnMessage, input); output != "" {
				fmt.Println(output)
			}

			fmt.Print("> ")
		}

		if err := sc.Err(); err != nil {
			c.handleEvent(EventOnError, err.Error())
		}

		if err := c.stopChat(); err != nil {
			c.handleEvent(EventOnError, err.Error())
		}
	}(sc)
}

func (c *LocalClient) handleEvent(eventKey, text string) string {
	if h, ok := c.handlers[eventKey]; ok {
		return h(text)
	}

	return ""
}

func (c *LocalClient) stopChat() error {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}

	if err := proc.Signal(os.Interrupt); err != nil {
		return err
	}

	return nil
}
