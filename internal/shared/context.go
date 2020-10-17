package shared

import (
	"github.com/supercaracal/tkkz-bot/internal/config"
)

// BotContext is
type BotContext struct {
	Option *config.BotOption
	Logger *config.BotLogger
}

// NewBotContext is
func NewBotContext(o *config.BotOption, l *config.BotLogger) *BotContext {
	return &BotContext{Option: o, Logger: l}
}
