package handler

import (
	"regexp"
	"strings"

	"github.com/supercaracal/tkkz-bot/internal/command"
	"github.com/supercaracal/tkkz-bot/internal/shared"
)

var (
	regexpForWhiteSpace = regexp.MustCompile(`[\s　]+`)
	regexpForMention    = regexp.MustCompile(`(?P<mention><@U[0-9A-Z]+>)`)
)

// EventHandler is
type EventHandler struct {
	ctx *shared.BotContext
}

// NewEventHandler is
func NewEventHandler(ctx *shared.BotContext) *EventHandler {
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
		return command.GetDefaultReply(h.ctx.Config.BrainURL, strings.Join(cmd, " "))
	}
}

func extractMentionIDsAndTokens(text string) ([]string, []string) {
	txt := regexpForMention.ReplaceAllString(text, " ${mention} ")
	txt = regexpForWhiteSpace.ReplaceAllString(txt, " ")
	txt = strings.TrimSpace(txt)
	words := strings.Split(txt, " ")
	ids := make([]string, 0, len(words))
	tokens := make([]string, 0, len(words))
	for _, word := range words {
		if strings.HasPrefix(word, "<@U") && strings.HasSuffix(word, ">") {
			ids = append(ids, string([]rune(word)[2:len(word)-1]))
		} else {
			tokens = append(tokens, word)
		}
	}

	return ids, tokens
}
