package http

import (
	"github.com/gin-gonic/gin"
	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	"gorm.io/gorm"
	"net/http"
)

func NewRouter(db *gorm.DB, redis *cache.Redis, worker *application.ViewWorker) *gin.Engine {
	repo := repository.NewPostgres(db)
	h := handler{submissions: application.NewSubmissionUseCase(repo, repo, redis), links: application.NewLinkUseCase(repo, repo, redis, worker)}
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusNoContent) })
	r.GET("/api/v1/shares/:shortCode", h.resolve)
	p := r.Group("/api/v1", requireActor())
	p.POST("/submissions", h.createSubmission)
	p.PATCH("/submissions/:submissionID/sharing", h.setSharing)
	p.POST("/submissions/:submissionID/share-links", h.createLink)
	p.PATCH("/share-links/:shareLinkID", h.updateLink)
	p.DELETE("/share-links/:shareLinkID", h.revokeLink)
	p.GET("/share-links/:shareLinkID/analytics", h.analytics)
	return r
}
