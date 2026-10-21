package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	"github.com/vucongthanh92/shortbit/internal/share/service"
	"gorm.io/gorm"
	"net/http"
)

func NewRouter(db *gorm.DB, redis *cache.Redis) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	h := handler{service: service.New(repository.NewPostgres(db), redis)}
	r.GET("/api/v1/shares/:shortCode", h.resolve)
	return r
}
