package handlers

import (
	"regexp"
	"strings"

	"github.com/supercaracal/tkkz-bot/internal/commands"
	"github.com/supercaracal/tkkz-bot/internal/configurations"
)

var (
	regexpForWhiteSpace = regexp.MustCompile(`[\s　]+`)
	regexpForMention    = regexp.MustCompile(`(?P<mention><@U[0-9A-Z]+>)`)
)

// EventHandler is
type EventHandler struct {
	cfg *configurations.MyConfig
}

// NewEventHandler is
func NewEventHandler(cfg *configurations.MyConfig) *EventHandler {
	return &EventHandler{cfg: cfg}
}

// RespondToContact is
func (h *EventHandler) RespondToContact(text string) string {
	mentionIds, command := extractMentionIDsAndTokens(text)
	for _, id := range mentionIds {
		if id == h.cfg.BotID {
			return doTask(command, h.cfg)
		}
	}
	return ""
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

func doTask(command []string, cfg *configurations.MyConfig) string {
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
