package chat

import (
	"fmt"
	"log"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// SlackSMClient is
type SlackSMClient struct {
	api      *slack.Client
	sm       *socketmode.Client
	handlers map[string]func(string) string
}

// NewSlackSMClient is
func NewSlackSMClient(appToken, botToken string, verbose bool, logger *log.Logger) Client {
	api := slack.New(
		botToken,
		slack.OptionDebug(verbose),
		slack.OptionLog(logger),
		slack.OptionAppLevelToken(appToken),
	)

	sm := socketmode.New(
		api,
		socketmode.OptionDebug(verbose),
		socketmode.OptionLog(logger),
	)

	return &SlackSMClient{
		api:      api,
		sm:       sm,
		handlers: map[string]func(string) string{},
	}
}

// ConnectAsync is
func (s *SlackSMClient) ConnectAsync(fail chan<- struct{}) {
	go func(sm *socketmode.Client, cn chan<- struct{}) {
		sm.Run()
		cn <- struct{}{}
	}(s.sm, fail)
}

// Disconnect is
func (s *SlackSMClient) Disconnect() error {
	return nil
}

// RegisterHandler is
func (s *SlackSMClient) RegisterHandler(event string, h func(string) string) {
	s.handlers[event] = h
}

// HandleEventsAsync is
func (s *SlackSMClient) HandleEventsAsync() {
	go s.startEventLoop()
}

func (s *SlackSMClient) startEventLoop() {
	for event := range s.sm.Events {
		switch event.Type {
		case socketmode.EventTypeEventsAPI:
			if ev, ok := event.Data.(slackevents.EventsAPIEvent); ok {
				s.sm.Ack(*event.Request)
				s.handleAPIEvent(ev)
			}
		case socketmode.EventTypeIncomingError, socketmode.EventTypeErrorWriteFailed, socketmode.EventTypeErrorBadMessage:
			s.handleEvent(EventOnError, "Something was wrong")
		case socketmode.EventTypeInvalidAuth:
			s.handleEvent(EventOnError, "Invalid Slack token")
		case socketmode.EventTypeConnected:
			s.handleEvent(EventOnConnection, "Connected to Slack via Socket Mode")
		case socketmode.EventTypeDisconnect:
			s.handleEvent(EventOnConnection, "Disconnected from Slack")
		default:
			continue
		}
	}
}

func (s *SlackSMClient) handleAPIEvent(ev slackevents.EventsAPIEvent) {
	switch ev.Type {
	case slackevents.CallbackEvent:
		switch ev.InnerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			s.handleMessageEvent(ev.InnerEvent.Data.(*slackevents.MessageEvent))
		}
	}
}

func (s *SlackSMClient) handleMessageEvent(ev *slackevents.MessageEvent) {
	if ev.SubType != "" || strings.HasPrefix(ev.Channel, "D") || ev.BotID != "" {
		return
	}

	reply := s.handleEvent(EventOnMessage, ev.Text)
	if reply == "" {
		return
	}

	_, _, err := s.api.PostMessage(ev.Channel,
		slack.MsgOptionText(reply, false), slack.MsgOptionTS(ev.ThreadTimeStamp))
	if err != nil {
		s.handleEvent(EventOnError, fmt.Sprintf("Failed to reply: %s", reply))
	}
}

func (s *SlackSMClient) handleEvent(eventKey, text string) string {
	if h, ok := s.handlers[eventKey]; ok {
		return h(text)
	}

	return ""
}
