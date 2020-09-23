package configurations

import (
	"fmt"
	"log"
	"os"
)

// MyLogger is
type MyLogger struct {
	Info *log.Logger
	Err  *log.Logger
}

// MyConfig is
type MyConfig struct {
	SlackToken string
	BotID      string
	Verbose    bool
	Logger     *MyLogger
}

// NewMyLogger is
func NewMyLogger() *MyLogger {
	return &MyLogger{
		Info: log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		Err:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// NewMyConfig is
func NewMyConfig() (*MyConfig, error) {
	cfg := MyConfig{
		SlackToken: os.Getenv("SLACK_TOKEN"),
		BotID:      os.Getenv("BOT_ID"),
	}

	if cfg.SlackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if cfg.BotID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	return &cfg, nil
}
