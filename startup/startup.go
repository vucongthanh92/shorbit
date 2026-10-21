package startup

import (
	"context"
	"github.com/vucongthanh92/shortbit/internal"
	"github.com/vucongthanh92/shortbit/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Execute() {
	cfg := config.Load()
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
