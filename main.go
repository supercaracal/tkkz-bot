package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/supercaracal/tkkz-bot/internal/chat"
	"github.com/supercaracal/tkkz-bot/internal/config"
	"github.com/supercaracal/tkkz-bot/internal/handler"
)

func waitUntil(
	chatCli chat.Client,
	fail <-chan struct{},
	sign <-chan os.Signal,
) error {

	select {
	case <-sign:
		return chatCli.Disconnect()
	case <-fail:
		return fmt.Errorf("Failed to manage connection")
	}
}

func main() {
	godotenv.Load()

	logger := config.NewBotLogger()
	opt, err := config.NewBotOption()
	if err != nil {
		logger.Err.Fatalln(err)
	}
	ctx := config.NewBotContext(opt, logger)

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&ctx.Option.Verbose, "verbose", false, "verbose log")
	f.BoolVar(&ctx.Option.Debug, "debug", false, "debug with stdin/stdout")
	f.Parse(os.Args[1:])

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	chatCli := chat.NewSlackClient(
		ctx.Option.SlackToken,
		ctx.Option.Verbose,
		ctx.Logger.Info,
	)
	if ctx.Option.Debug {
		chatCli = chat.NewLocalClient()
	}

	h := handler.NewEventHandler(ctx)
	chatCli.RegisterHandler(chat.EventOnConnection, h.LogAsInfo)
	chatCli.RegisterHandler(chat.EventOnMessage, h.RespondToContact)
	chatCli.RegisterHandler(chat.EventOnError, h.LogAsErr)

	fail := make(chan struct{})
	chatCli.ConnectAsync(fail)
	chatCli.HandleEventsAsync()

	err = waitUntil(chatCli, fail, sign)
	if err != nil {
		ctx.Logger.Err.Fatalln(err)
	}

	ctx.Logger.Info.Println("Shutting Down")
	os.Exit(0)
}
