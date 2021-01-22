package handler

import (
	"strings"

	"github.com/supercaracal/tkkz-bot/internal/command"
	"github.com/supercaracal/tkkz-bot/internal/config"
)

// EventHandler is
type EventHandler struct {
	cfg *config.BotConfig
}

// NewEventHandler is
func NewEventHandler(cfg *config.BotConfig) *EventHandler {
	return &EventHandler{cfg: cfg}
}

// RespondToContact is
func (h *EventHandler) RespondToContact(text string) string {
	_, cmd := extractMentionIDsAndTokens(text)
	return h.doTask(cmd)
}

// LogAsInfo is
func (h *EventHandler) LogAsInfo(text string) string {
	h.cfg.Logger.Info.Println(text)
	return ""
}

// LogAsErr is
func (h *EventHandler) LogAsErr(text string) string {
	h.cfg.Logger.Err.Println(text)
	return ""
}

func (h *EventHandler) doTask(cmd []string) string {
	if len(cmd) == 0 {
		return "Hi,"
	}

	switch strings.ToLower(cmd[0]) {
	case "ping":
		return command.GetPingReply()
	case "pang":
		return command.GetPangReply()
	default:
		return command.GetDefaultReply(h.cfg.Option.BrainURL, strings.Join(cmd, " "))
	}
}
