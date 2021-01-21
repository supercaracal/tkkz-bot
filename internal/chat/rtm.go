package chat

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
)

// SlackRTMClient is
type SlackRTMClient struct {
	rtm      *slack.RTM
	handlers map[string]func(string) string
}

// NewSlackRTMClient is
func NewSlackRTMClient(token string, verbose bool, logger *log.Logger) Client {
	api := slack.New(token, slack.OptionDebug(verbose), slack.OptionLog(logger))

	return &SlackRTMClient{
		rtm:      api.NewRTM(),
		handlers: map[string]func(string) string{},
	}
}

// ConnectAsync is
func (s *SlackRTMClient) ConnectAsync(fail chan<- struct{}) {
	go func(rtm *slack.RTM, cn chan<- struct{}) {
		rtm.ManageConnection()
		cn <- struct{}{}
	}(s.rtm, fail)
}

// Disconnect is
func (s *SlackRTMClient) Disconnect() error {
	return s.rtm.Disconnect()
}

// RegisterHandler is
func (s *SlackRTMClient) RegisterHandler(event string, h func(string) string) {
	s.handlers[event] = h
}

// HandleEventsAsync is
func (s *SlackRTMClient) HandleEventsAsync() {
	go s.startEventLoop()
}

func (s *SlackRTMClient) startEventLoop() {
	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			s.handleMessageEvent(ev)
		case *slack.RTMError:
			reply := fmt.Sprintf("CODE=%d %s", ev.Code, ev.Error())
			s.handleEvent(EventOnError, reply)
		case *slack.InvalidAuthEvent:
			s.handleEvent(EventOnError, "Invalid Slack token")
		case *slack.ConnectedEvent:
			s.handleEvent(EventOnConnection, "Connected to Slack via RTM")
		case *slack.DisconnectedEvent:
			reply := fmt.Sprintf("Disconnected from Slack: intentionally=%t %s", ev.Intentional, ev.Cause.Error())
			s.handleEvent(EventOnConnection, reply)
		default:
			continue
		}
	}
}

func (s *SlackRTMClient) handleEvent(eventKey, text string) string {
	if h, ok := s.handlers[eventKey]; ok {
		return h(text)
	}

	return ""
}

func (s *SlackRTMClient) handleMessageEvent(ev *slack.MessageEvent) {
	if ev.SubType != "" || strings.HasPrefix(ev.Channel, "D") || ev.BotID != "" {
		return
	}

	reply := s.handleEvent(EventOnMessage, ev.Text)
	if reply == "" {
		return
	}

	resp := s.rtm.NewOutgoingMessage(reply, ev.Channel)
	resp.ThreadTimestamp = ev.ThreadTimestamp
	s.rtm.SendMessage(resp)
}
