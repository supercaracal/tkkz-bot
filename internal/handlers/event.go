package handlers

import (
	"regexp"
	"strings"

	"github.com/supercaracal/tkkz-bot/internal/commands"
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
	mentionIds, command := extractMentionIDsAndTokens(text)
	for _, id := range mentionIds {
		if id == h.ctx.Config.BotID {
			return h.doTask(command)
		}
	}
	return ""
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

func (h *EventHandler) doTask(command []string) string {
	if len(command) == 0 {
		return "Hi,"
	}

	switch strings.ToLower(command[0]) {
	case "ping":
		return commands.GetPingReply()
	default:
		return commands.GetDefaultReply()
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
