package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LinkStatus string

const (
	LinkStatusActive  LinkStatus = "active"
	LinkStatusRevoked LinkStatus = "revoked"
)

type AccessLevel string

const (
	AccessPublic        AccessLevel = "public"
	AccessAuthenticated AccessLevel = "authenticated"
	AccessRestricted    AccessLevel = "restricted"
)

type Submission struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey"`
	OwnerUserID  uuid.UUID       `gorm:"type:uuid;not null"`
	Title        string          `gorm:"size:255;not null"`
	Content      json.RawMessage `gorm:"type:jsonb;not null"`
	ShareEnabled bool            `gorm:"not null"`
	DeletedAt    *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Submission) TableName() string { return "shortbit.submissions" }

type ShareLink struct {
	ID           uuid.UUID   `gorm:"type:uuid;primaryKey"`
	SubmissionID uuid.UUID   `gorm:"type:uuid;not null"`
	ShortCode    string      `gorm:"uniqueIndex;size:32;not null"`
	Status       LinkStatus  `gorm:"size:20;not null"`
	AccessLevel  AccessLevel `gorm:"size:20;not null"`
	ViewCount    int64       `gorm:"not null"`
	ExpiresAt    *time.Time
	RevokedAt    *time.Time
	CreatedBy    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Submission   Submission
}

func (ShareLink) TableName() string { return "shortbit.share_links" }
func (s ShareLink) IsAvailable(now time.Time, viewerID *uuid.UUID) bool {
	return s.Status == LinkStatusActive && s.RevokedAt == nil && (s.ExpiresAt == nil || s.ExpiresAt.After(now)) && s.Submission.ShareEnabled && s.Submission.DeletedAt == nil && s.AccessLevel != AccessRestricted && (s.AccessLevel != AccessAuthenticated || viewerID != nil)
}

type DailyViews struct {
	ShareLinkID uuid.UUID `gorm:"type:uuid;primaryKey"`
	ViewedDate  time.Time `gorm:"type:date;primaryKey"`
	TotalViews  int64
	UniqueViews int64
}

func (DailyViews) TableName() string { return "shortbit.share_link_view_daily" }
