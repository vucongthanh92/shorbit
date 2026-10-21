//go:build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/vucongthanh92/shortbit/internal/api"
	apiHTTP "github.com/vucongthanh92/shortbit/internal/api/http"
	"github.com/vucongthanh92/shortbit/internal/config"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/platform/database"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	"gorm.io/gorm"
)

func InitializeContainer(cfg config.Config) (*api.ApiContainer, error) {
	wire.Build(provideDB, provideRedis, repository.NewPostgres, provideViewWorker, apiHTTP.NewServer, api.NewApiContainer)
	return nil, nil
}
func provideDB(cfg config.Config) (*gorm.DB, error) { return database.Open(cfg.DatabaseURL) }
func provideRedis(cfg config.Config) *cache.Redis {
	return cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
}
func provideViewWorker(repo *repository.Postgres) *application.ViewWorker {
	return application.NewViewWorker(repo, 1000)
}
