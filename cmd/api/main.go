package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vucongthanh92/shortbit/internal/config"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/platform/database"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	sharehttp "github.com/vucongthanh92/shortbit/internal/share/transport/http"
)

func main() {
	cfg := config.Load()
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	redisCache := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	defer redisCache.Close()

	worker := application.NewViewWorker(repository.NewPostgres(db), 1_000)
	workerContext, stopWorker := context.WithCancel(context.Background())
	defer stopWorker()
	worker.Start(workerContext)

	router := sharehttp.NewRouter(db, redisCache, worker)
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: router, ReadHeaderTimeout: 5 * time.Second}

	go func() {
		log.Printf("shortbit API listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
