package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/supercaracal/tkkz-bot/infra"
)

type myLogger struct {
	info *log.Logger
	err  *log.Logger
}

type myConfig struct {
	slackToken string
	botID      string
	verbose    bool
	logger     myLogger
}

var (
	regexpForWhiteSpace = regexp.MustCompile(`[\s　]+`)
	regexpForMention    = regexp.MustCompile(`(?P<mention><@U[0-9A-Z]+>)`)
)

func doTask(command []string, cfg *myConfig) string {
	if len(command) == 0 {
		return "Hi,"
	}

	switch strings.ToLower(command[0]) {
	case "ping":
		return "PONG"
	default:
		return "?"
	}
}

func (cfg *myConfig) respondToContact(text string) string {
	mentionIds, command := extractMentionIDsAndTokens(text)
	for _, id := range mentionIds {
		if id == cfg.botID {
			return doTask(command, cfg)
		}
	}
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

func (cfg *myConfig) logAsInfo(text string) string {
	cfg.logger.info.Println(text)
	return ""
}

func (cfg *myConfig) logAsErr(text string) string {
	cfg.logger.err.Println(text)
	return ""
}

func newMyConfig() (*myConfig, error) {
	cfg := myConfig{
		slackToken: os.Getenv("SLACK_TOKEN"),
		botID:      os.Getenv("BOT_ID"),
	}

	if cfg.slackToken == "" {
		return nil, fmt.Errorf("SLACK_TOKEN environment variable required")
	}

	if cfg.botID == "" {
		return nil, fmt.Errorf("BOT_ID environment variable required")
	}

	return &cfg, nil
}

func waitForShutdown(chat Chat, fail <-chan struct{}, sign <-chan os.Signal) error {
	select {
	case <-sign:
		return chat.Disconnect()
	case <-fail:
		return fmt.Errorf("Failed to manage connection")
	}
}

func main() {
	logger := myLogger{
		info: log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		err:  log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}

	godotenv.Load()

	cfg, err := newMyConfig()
	if err != nil {
		logger.err.Fatalln(err)
	}

	cfg.logger = logger

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&cfg.verbose, "verbose", false, "verbose log")
	f.Parse(os.Args[1:])

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	var chat Chat = infra.NewSlackClient(cfg.slackToken, cfg.verbose, cfg.logger.info)
	chat.RegisterHandler("onMessage", cfg.respondToContact)
	chat.RegisterHandler("onMessagingError", cfg.logAsErr)
	chat.RegisterHandler("onAuthenticationError", cfg.logAsErr)
	chat.RegisterHandler("onConnected", cfg.logAsInfo)
	chat.RegisterHandler("onDisconnected", cfg.logAsInfo)
	fail := make(chan struct{})
	chat.ConnectAsync(fail)
	chat.HandleEventsAsync()

	err = waitForShutdown(chat, fail, sign)
	if err != nil {
		cfg.logger.err.Fatalln(err)
	}

	cfg.logger.info.Println("Shutting Down")
	os.Exit(0)
}
