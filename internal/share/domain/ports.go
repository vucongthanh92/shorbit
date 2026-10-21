package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound  = errors.New("resource not found")
	ErrForbidden = errors.New("forbidden")
	ErrConflict  = errors.New("conflict")
)

type SubmissionRepository interface {
	Create(context.Context, Submission) error
	FindByID(context.Context, uuid.UUID) (Submission, error)
	UpdateShareEnabled(context.Context, uuid.UUID, bool) error
}
type ShareLinkRepository interface {
	CreateLink(context.Context, ShareLink) error
	FindByIDLink(context.Context, uuid.UUID) (ShareLink, error)
	FindByCode(context.Context, string) (ShareLink, error)
	FindBySubmissionID(context.Context, uuid.UUID) ([]ShareLink, error)
	Update(context.Context, ShareLink) error
	Revoke(context.Context, uuid.UUID, time.Time) error
	RecordView(context.Context, uuid.UUID, time.Time) error
	GetDailyViews(context.Context, uuid.UUID, time.Time, time.Time) ([]DailyViews, error)
}
type LinkCache interface {
	GetJSON(context.Context, string, any) (bool, error)
	SetJSON(context.Context, string, any) error
	Delete(context.Context, string) error
}
type ViewRecorder interface{ Record(uuid.UUID, time.Time) }
