package http

import (
	"context"
	"net/http"
	"time"

	"github.com/vucongthanh92/shortbit/config"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	shareHTTP "github.com/vucongthanh92/shortbit/internal/share/transport/http"
	"gorm.io/gorm"
)

type Server struct {
	cfg     *config.AppConfig
	handler http.Handler
}

func NewServer(cfg *config.AppConfig, db *gorm.DB, redis *cache.Redis, worker *application.ViewWorker) *Server {
	return &Server{cfg: cfg, handler: shareHTTP.NewRouter(db, redis, worker)}
}
func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{Addr: s.cfg.Http.Port, Handler: s.handler, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}
