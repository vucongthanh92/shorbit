package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
	"github.com/vucongthanh92/shortbit/internal/share/service"
	"net/http"
)

type handler struct{ service *service.Service }

func (h handler) resolve(c *gin.Context) {
	link, err := h.service.Resolve(c.Request.Context(), c.Param("shortCode"))
	if errors.Is(err, repository.ErrNotFound) || errors.Is(err, service.ErrUnavailable) {
		c.JSON(http.StatusNotFound, gin.H{"error": "share link not found or unavailable"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resolve share link"})
		return
	}
	// View persistence is deliberately deferred to the planned event/outbox worker.
	c.JSON(http.StatusOK, gin.H{"shortCode": link.ShortCode, "submission": gin.H{"id": link.Submission.ID, "title": link.Submission.Title, "content": link.Submission.Content}})
}
