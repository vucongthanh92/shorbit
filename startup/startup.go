package startup

import (
	"context"
	"github.com/vucongthanh92/go-base-utils/command"
	"github.com/vucongthanh92/shortbit/config"
	"github.com/vucongthanh92/shortbit/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var cfg *config.AppConfig

func Execute() {
	cfg = &config.AppConfig{}
	command.UseCommands(command.WithStartCommand(start, cfg))
}
func start() {
	container, err := internal.InitializeContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	container.ViewWorker.Start(ctx)
	if err := container.HttpServer.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
