package domain

import (
	"github.com/google/uuid"
	"time"
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
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerUserID          uuid.UUID `gorm:"type:uuid;not null"`
	Title                string
	Content              []byte `gorm:"type:jsonb"`
	ShareEnabled         bool
	DeletedAt            *time.Time
	CreatedAt, UpdatedAt time.Time
}

func (Submission) TableName() string { return "submissions" }

type ShareLink struct {
	ID                   uuid.UUID   `gorm:"type:uuid;primaryKey"`
	SubmissionID         uuid.UUID   `gorm:"type:uuid;not null"`
	ShortCode            string      `gorm:"uniqueIndex;size:32"`
	Status               LinkStatus  `gorm:"size:20"`
	AccessLevel          AccessLevel `gorm:"size:20"`
	ViewCount            int64
	ExpiresAt, RevokedAt *time.Time
	CreatedBy            uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt, UpdatedAt time.Time
	Submission           Submission
}

func (ShareLink) TableName() string { return "share_links" }

func (s ShareLink) IsAccessible(now time.Time) bool {
	return s.Status == LinkStatusActive && s.RevokedAt == nil && (s.ExpiresAt == nil || s.ExpiresAt.After(now)) && s.Submission.ShareEnabled && s.Submission.DeletedAt == nil
}
