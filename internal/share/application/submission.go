package application

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"strings"
	"time"
)

type SubmissionUseCase struct {
	repo  domain.SubmissionRepository
	links domain.ShareLinkRepository
	cache domain.LinkCache
	now   func() time.Time
}

func NewSubmissionUseCase(repo domain.SubmissionRepository, links domain.ShareLinkRepository, cache domain.LinkCache) *SubmissionUseCase {
	return &SubmissionUseCase{repo: repo, links: links, cache: cache, now: time.Now}
}

type CreateSubmissionInput struct {
	OwnerID      uuid.UUID
	Title        string
	Content      json.RawMessage
	ShareEnabled bool
}

func (u *SubmissionUseCase) Create(ctx context.Context, in CreateSubmissionInput) (domain.Submission, error) {
	if in.OwnerID == uuid.Nil || strings.TrimSpace(in.Title) == "" {
		return domain.Submission{}, domain.ErrForbidden
	}
	if len(in.Content) == 0 {
		in.Content = json.RawMessage(`{}`)
	}
	s := domain.Submission{ID: uuid.New(), OwnerUserID: in.OwnerID, Title: strings.TrimSpace(in.Title), Content: in.Content, ShareEnabled: in.ShareEnabled, CreatedAt: u.now(), UpdatedAt: u.now()}
	return s, u.repo.Create(ctx, s)
}
func (u *SubmissionUseCase) SetSharing(ctx context.Context, actor, id uuid.UUID, enabled bool) error {
	s, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if s.OwnerUserID != actor {
		return domain.ErrForbidden
	}
	if err = u.repo.UpdateShareEnabled(ctx, id, enabled); err != nil {
		return err
	}
	links, err := u.links.FindBySubmissionID(ctx, id)
	if err != nil {
		return err
	}
	for _, link := range links {
		_ = u.cache.Delete(ctx, "share-link:"+link.ShortCode)
	}
	return nil
}
