package infrastructures

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

// MySlack is
type MySlack struct {
	rtm      *slack.RTM
	handlers map[string]func(string) string
}

// NewSlackClient is
func NewSlackClient(token string, verbose bool, logger *log.Logger) *MySlack {
	api := slack.New(token, slack.OptionDebug(verbose), slack.OptionLog(logger))

	return &MySlack{
		rtm: api.NewRTM(),
		handlers: map[string]func(string) string{
			"onMessage":             noop,
			"onMessagingError":      noop,
			"onAuthenticationError": noop,
			"onConnected":           noop,
			"onDisconnected":        noop,
		},
	}
}

// ConnectAsync is
func (s *MySlack) ConnectAsync(fail chan<- struct{}) {
	go func(rtm *slack.RTM, cn chan<- struct{}) {
		rtm.ManageConnection()
		cn <- struct{}{}
	}(s.rtm, fail)
}

// Disconnect is
func (s *MySlack) Disconnect() error {
	return s.rtm.Disconnect()
}

// RegisterHandler is
func (s *MySlack) RegisterHandler(event string, h func(string) string) {
	s.handlers[event] = h
}

// HandleEventsAsync is
func (s *MySlack) HandleEventsAsync() {
	go s.startEventLoop()
}

func (s *MySlack) startEventLoop() {
	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			s.handleMessageEvent(ev)
		case *slack.RTMError:
			reply := fmt.Sprintf("CODE=%d %s", ev.Code, ev.Error())
			s.handlers["onMessagingError"](reply)
		case *slack.InvalidAuthEvent:
			s.handlers["onAuthenticationError"]("Invalid Slack token")
		case *slack.ConnectedEvent:
			s.handlers["onConnected"]("Connected to Slack")
		case *slack.DisconnectedEvent:
			reply := fmt.Sprintf("Disconnected from Slack: intentionally=%t %s", ev.Intentional, ev.Cause.Error())
			s.handlers["onDisconnected"](reply)
		default:
			continue
		}
	}
}

func (s *MySlack) handleMessageEvent(ev *slack.MessageEvent) {
	if ev.SubType != "" || strings.HasPrefix(ev.Channel, "D") || ev.BotID != "" {
		return
	}

	reply := s.handlers["onMessage"](ev.Text)
	if reply == "" {
		return
	}

	resp := s.rtm.NewOutgoingMessage(reply, ev.Channel)
	resp.ThreadTimestamp = ev.ThreadTimestamp
	s.rtm.SendMessage(resp)
}

func noop(a string) string {
	return a
}
