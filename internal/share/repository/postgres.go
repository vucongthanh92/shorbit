package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Postgres is the infrastructure adapter implementing domain repository ports.
type Postgres struct{ db *gorm.DB }

func NewPostgres(db *gorm.DB) *Postgres { return &Postgres{db: db} }

func (r *Postgres) Create(ctx context.Context, submission domain.Submission) error {
	return r.db.WithContext(ctx).Create(&submission).Error
}
func (r *Postgres) FindByID(ctx context.Context, id uuid.UUID) (domain.Submission, error) {
	var s domain.Submission
	return s, normalize(r.db.WithContext(ctx).First(&s, "id = ?", id).Error)
}
func (r *Postgres) UpdateShareEnabled(ctx context.Context, id uuid.UUID, enabled bool) error {
	return normalize(r.db.WithContext(ctx).Model(&domain.Submission{}).Where("id = ?", id).Update("share_enabled", enabled).Error)
}

func (r *Postgres) CreateLink(ctx context.Context, link domain.ShareLink) error {
	err := r.db.WithContext(ctx).Create(&link).Error
	if unique(err) {
		return domain.ErrConflict
	}
	return err
}
func (r *Postgres) FindByIDLink(ctx context.Context, id uuid.UUID) (domain.ShareLink, error) {
	var link domain.ShareLink
	return link, normalize(r.db.WithContext(ctx).Preload("Submission").First(&link, "id = ?", id).Error)
}
func (r *Postgres) FindByCode(ctx context.Context, code string) (domain.ShareLink, error) {
	var link domain.ShareLink
	return link, normalize(r.db.WithContext(ctx).Preload("Submission").First(&link, "short_code = ?", code).Error)
}
func (r *Postgres) FindBySubmissionID(ctx context.Context, submissionID uuid.UUID) ([]domain.ShareLink, error) {
	var links []domain.ShareLink
	return links, r.db.WithContext(ctx).Where("submission_id = ?", submissionID).Find(&links).Error
}
func (r *Postgres) Update(ctx context.Context, link domain.ShareLink) error {
	return normalize(r.db.WithContext(ctx).Model(&domain.ShareLink{}).Where("id = ?", link.ID).Updates(map[string]any{"access_level": link.AccessLevel, "expires_at": link.ExpiresAt}).Error)
}
func (r *Postgres) Revoke(ctx context.Context, id uuid.UUID, at time.Time) error {
	return normalize(r.db.WithContext(ctx).Model(&domain.ShareLink{}).Where("id = ?", id).Updates(map[string]any{"status": domain.LinkStatusRevoked, "revoked_at": at}).Error)
}

func (r *Postgres) RecordView(ctx context.Context, id uuid.UUID, viewedAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.ShareLink{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error; err != nil {
			return err
		}
		day := time.Date(viewedAt.Year(), viewedAt.Month(), viewedAt.Day(), 0, 0, 0, 0, time.UTC)
		return tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "share_link_id"}, {Name: "viewed_date"}}, DoUpdates: clause.Assignments(map[string]any{"total_views": gorm.Expr("shortbit.share_link_view_daily.total_views + 1")})}).Create(&domain.DailyViews{ShareLinkID: id, ViewedDate: day, TotalViews: 1}).Error
	})
}
func (r *Postgres) GetDailyViews(ctx context.Context, id uuid.UUID, from, to time.Time) ([]domain.DailyViews, error) {
	var views []domain.DailyViews
	return views, r.db.WithContext(ctx).Where("share_link_id = ? AND viewed_date >= ? AND viewed_date <= ?", id, from, to).Order("viewed_date").Find(&views).Error
}
func normalize(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}
func unique(err error) bool {
	return err != nil && (errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key"))
}
