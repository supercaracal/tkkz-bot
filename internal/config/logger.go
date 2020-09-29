package config

import (
	"log"
	"os"
)

// BotLogger is
type BotLogger struct {
	Info *log.Logger
	Err  *log.Logger
}

// NewBotLogger is
func NewBotLogger() *BotLogger {
	return &BotLogger{
		Info: log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		Err:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}
