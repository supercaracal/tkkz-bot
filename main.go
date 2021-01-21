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
	clients []chat.Client,
	fail <-chan struct{},
	sign <-chan os.Signal,
) error {

	select {
	case <-sign:
		for _, cli := range clients {
			_ = cli.Disconnect()
		}
		return nil
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

	fail := make(chan struct{})
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

	clients := make([]chat.Client, 0, 3)
	if ctx.Option.Debug {
		clients = append(clients, chat.NewLocalClient())
	} else {
		c1 := chat.NewSlackRTMClient(ctx.Option.SlackToken, ctx.Option.Verbose, ctx.Logger.Info)
		c2 := chat.NewSlackSMClient(ctx.Option.SlackAppToken, ctx.Option.SlackBotToken, ctx.Option.Verbose, ctx.Logger.Info)
		clients = append(clients, c1, c2)
	}

	h := handler.NewEventHandler(ctx)
	for _, cli := range clients {
		cli.RegisterHandler(chat.EventOnConnection, h.LogAsInfo)
		cli.RegisterHandler(chat.EventOnMessage, h.RespondToContact)
		cli.RegisterHandler(chat.EventOnError, h.LogAsErr)
		cli.ConnectAsync(fail)
		cli.HandleEventsAsync()
	}

	err = waitUntil(clients, fail, sign)
	if err != nil {
		ctx.Logger.Err.Fatalln(err)
	}

	ctx.Logger.Info.Println("Shutting Down")
	os.Exit(0)
}
