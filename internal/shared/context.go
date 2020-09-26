package shared

import (
	"github.com/supercaracal/tkkz-bot/internal/configs"
)

// BotContext is
type BotContext struct {
	Config *configs.BotConfig
	Logger *configs.BotLogger
}

// NewBotContext is
func NewBotContext(c *configs.BotConfig, l *configs.BotLogger) *BotContext {
	return &BotContext{Config: c, Logger: l}
}
