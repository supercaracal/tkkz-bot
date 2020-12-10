package handler

import (
	"strings"

	"github.com/supercaracal/tkkz-bot/internal/command"
	"github.com/supercaracal/tkkz-bot/internal/config"
)

// EventHandler is
type EventHandler struct {
	ctx *config.BotContext
}

// NewEventHandler is
func NewEventHandler(ctx *config.BotContext) *EventHandler {
	return &EventHandler{ctx: ctx}
}

// RespondToContact is
func (h *EventHandler) RespondToContact(text string) string {
	_, cmd := extractMentionIDsAndTokens(text)
	return h.doTask(cmd)
}

// LogAsInfo is
func (h *EventHandler) LogAsInfo(text string) string {
	h.ctx.Logger.Info.Println(text)
	return ""
}

// LogAsErr is
func (h *EventHandler) LogAsErr(text string) string {
	h.ctx.Logger.Err.Println(text)
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
		return command.GetDefaultReply(h.ctx.Option.BrainURL, strings.Join(cmd, " "))
	}
}
