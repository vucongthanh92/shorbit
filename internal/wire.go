//go:build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/vucongthanh92/shortbit/config"
	"github.com/vucongthanh92/shortbit/internal/api"
	apiHTTP "github.com/vucongthanh92/shortbit/internal/api/http"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/platform/database"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	"gorm.io/gorm"
	"time"
)

func InitializeContainer(cfg *config.AppConfig) (*api.ApiContainer, error) {
	wire.Build(provideDB, provideRedis, repository.NewPostgres, provideViewWorker, apiHTTP.NewServer, api.NewApiContainer)
	return nil, nil
}
func provideDB(cfg *config.AppConfig) (*gorm.DB, error) {
	return database.Open(cfg.Database.WriteDbCfg.ConnectionString)
}
func provideRedis(cfg *config.AppConfig) *cache.Redis {
	return cache.NewRedis(cfg.Redis.Addrs[0], cfg.Redis.Password, cfg.Redis.DB, time.Duration(cfg.CacheTTLSeconds)*time.Second)
}
func provideViewWorker(repo *repository.Postgres) *application.ViewWorker {
	return application.NewViewWorker(repo, 1000)
}
