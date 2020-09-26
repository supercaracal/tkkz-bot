package chats

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

// SlackClient is
type SlackClient struct {
	rtm      *slack.RTM
	handlers map[string]func(string) string
}

// NewSlackClient is
func NewSlackClient(token string, verbose bool, logger *log.Logger) ChatClient {
	api := slack.New(token, slack.OptionDebug(verbose), slack.OptionLog(logger))

	return &SlackClient{
		rtm:      api.NewRTM(),
		handlers: map[string]func(string) string{},
	}
}

// ConnectAsync is
func (s *SlackClient) ConnectAsync(fail chan<- struct{}) {
	go func(rtm *slack.RTM, cn chan<- struct{}) {
		rtm.ManageConnection()
		cn <- struct{}{}
	}(s.rtm, fail)
}

// Disconnect is
func (s *SlackClient) Disconnect() error {
	return s.rtm.Disconnect()
}

// RegisterHandler is
func (s *SlackClient) RegisterHandler(event string, h func(string) string) {
	s.handlers[event] = h
}

// HandleEventsAsync is
func (s *SlackClient) HandleEventsAsync() {
	go s.startEventLoop()
}

func (s *SlackClient) startEventLoop() {
	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			s.handleMessageEvent(ev)
		case *slack.RTMError:
			reply := fmt.Sprintf("CODE=%d %s", ev.Code, ev.Error())
			s.handleEvent("onMessagingError", reply)
		case *slack.InvalidAuthEvent:
			s.handleEvent("onAuthenticationError", "Invalid Slack token")
		case *slack.ConnectedEvent:
			s.handleEvent("onConnected", "Connected to Slack")
		case *slack.DisconnectedEvent:
			reply := fmt.Sprintf("Disconnected from Slack: intentionally=%t %s", ev.Intentional, ev.Cause.Error())
			s.handleEvent("onDisconnected", reply)
		default:
			continue
		}
	}
}

func (s *SlackClient) handleEvent(eventKey, text string) string {
	if h, ok := s.handlers[eventKey]; ok {
		return h(text)
	}

	return ""
}

func (s *SlackClient) handleMessageEvent(ev *slack.MessageEvent) {
	if ev.SubType != "" || strings.HasPrefix(ev.Channel, "D") || ev.BotID != "" {
		return
	}

	reply := s.handleEvent("onMessage", ev.Text)
	if reply == "" {
		return
	}

	resp := s.rtm.NewOutgoingMessage(reply, ev.Channel)
	resp.ThreadTimestamp = ev.ThreadTimestamp
	s.rtm.SendMessage(resp)
}
