package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("share link not found")

type Postgres struct{ db *gorm.DB }

func NewPostgres(db *gorm.DB) *Postgres { return &Postgres{db: db} }
func (r *Postgres) FindByCode(ctx context.Context, code string) (domain.ShareLink, error) {
	var link domain.ShareLink
	err := r.db.WithContext(ctx).Preload("Submission").Where("short_code = ?", code).First(&link).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return link, ErrNotFound
	}
	return link, err
}
func (r *Postgres) IncrementView(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.ShareLink{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
