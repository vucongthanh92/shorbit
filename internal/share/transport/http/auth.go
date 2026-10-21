package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

const actorKey = "actorID"

func requireActor() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, e := uuid.Parse(c.GetHeader("X-User-ID"))
		if e != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-ID must be a UUID"})
			c.Abort()
			return
		}
		c.Set(actorKey, id)
		c.Next()
	}
}
func actor(c *gin.Context) uuid.UUID { return c.MustGet(actorKey).(uuid.UUID) }
func optionalActor(c *gin.Context) *uuid.UUID {
	id, e := uuid.Parse(c.GetHeader("X-User-ID"))
	if e != nil {
		return nil
	}
	return &id
}
