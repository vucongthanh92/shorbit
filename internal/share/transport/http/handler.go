package http

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/application"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"net/http"
	"time"
)

type handler struct {
	submissions *application.SubmissionUseCase
	links       *application.LinkUseCase
}
type createSubmissionRequest struct {
	Title        string          `json:"title" binding:"required,max=255"`
	Content      json.RawMessage `json:"content"`
	ShareEnabled bool            `json:"shareEnabled"`
}
type setSharingRequest struct {
	Enabled bool `json:"enabled"`
}
type createLinkRequest struct {
	AccessLevel domain.AccessLevel `json:"accessLevel"`
	ExpiresAt   *time.Time         `json:"expiresAt"`
}
type updateLinkRequest struct {
	AccessLevel *domain.AccessLevel `json:"accessLevel"`
	ExpiresAt   *time.Time          `json:"expiresAt"`
}

func (h handler) createSubmission(c *gin.Context) {
	var r createSubmissionRequest
	if c.ShouldBindJSON(&r) != nil {
		bad(c)
		return
	}
	s, e := h.submissions.Create(c, application.CreateSubmissionInput{OwnerID: actor(c), Title: r.Title, Content: r.Content, ShareEnabled: r.ShareEnabled})
	if fail(c, e) {
		return
	}
	c.JSON(http.StatusCreated, submission(s))
}
func (h handler) setSharing(c *gin.Context) {
	id, ok := id(c, "submissionID")
	if !ok {
		return
	}
	var r setSharingRequest
	if c.ShouldBindJSON(&r) != nil {
		bad(c)
		return
	}
	if fail(c, h.submissions.SetSharing(c, actor(c), id, r.Enabled)) {
		return
	}
	c.Status(http.StatusNoContent)
}
func (h handler) createLink(c *gin.Context) {
	sid, ok := id(c, "submissionID")
	if !ok {
		return
	}
	var r createLinkRequest
	if c.ShouldBindJSON(&r) != nil {
		bad(c)
		return
	}
	link, e := h.links.Create(c, application.CreateLinkInput{ActorID: actor(c), SubmissionID: sid, AccessLevel: r.AccessLevel, ExpiresAt: r.ExpiresAt})
	if fail(c, e) {
		return
	}
	c.JSON(http.StatusCreated, linkJSON(link))
}
func (h handler) updateLink(c *gin.Context) {
	lid, ok := id(c, "shareLinkID")
	if !ok {
		return
	}
	var r updateLinkRequest
	if c.ShouldBindJSON(&r) != nil {
		bad(c)
		return
	}
	var expiry **time.Time
	if r.ExpiresAt != nil {
		expiry = &r.ExpiresAt
	}
	link, e := h.links.Update(c, application.UpdateLinkInput{ActorID: actor(c), LinkID: lid, AccessLevel: r.AccessLevel, ExpiresAt: expiry})
	if fail(c, e) {
		return
	}
	c.JSON(http.StatusOK, linkJSON(link))
}
func (h handler) revokeLink(c *gin.Context) {
	lid, ok := id(c, "shareLinkID")
	if !ok {
		return
	}
	if fail(c, h.links.Revoke(c, actor(c), lid)) {
		return
	}
	c.Status(http.StatusNoContent)
}
func (h handler) resolve(c *gin.Context) {
	link, e := h.links.Resolve(c, c.Param("shortCode"), optionalActor(c))
	if fail(c, e) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"shortCode": link.ShortCode, "submission": submission(link.Submission)})
}
func (h handler) analytics(c *gin.Context) {
	lid, ok := id(c, "shareLinkID")
	if !ok {
		return
	}
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)
	if raw := c.Query("from"); raw != "" {
		v, e := time.Parse("2006-01-02", raw)
		if e != nil {
			bad(c)
			return
		}
		from = v
	}
	if raw := c.Query("to"); raw != "" {
		v, e := time.Parse("2006-01-02", raw)
		if e != nil {
			bad(c)
			return
		}
		to = v
	}
	link, rows, e := h.links.Analytics(c, actor(c), lid, from, to)
	if fail(c, e) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"shareLink": linkJSON(link), "dailyViews": rows, "note": "uniqueViews is not calculated in MVP because no visitor/event table exists"})
}
func id(c *gin.Context, key string) (uuid.UUID, bool) {
	v, e := uuid.Parse(c.Param(key))
	if e != nil {
		bad(c)
		return uuid.Nil, false
	}
	return v, true
}
func submission(s domain.Submission) gin.H {
	return gin.H{"id": s.ID, "title": s.Title, "content": s.Content, "shareEnabled": s.ShareEnabled}
}
func linkJSON(l domain.ShareLink) gin.H {
	return gin.H{"id": l.ID, "submissionId": l.SubmissionID, "shortCode": l.ShortCode, "status": l.Status, "accessLevel": l.AccessLevel, "viewCount": l.ViewCount, "expiresAt": l.ExpiresAt, "createdAt": l.CreatedAt}
}
func bad(c *gin.Context) { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"}) }
func fail(c *gin.Context, e error) bool {
	if e == nil {
		return false
	}
	switch {
	case errors.Is(e, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.Is(e, domain.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(e, domain.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	return true
}
