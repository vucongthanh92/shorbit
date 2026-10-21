package application

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/google/uuid"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"time"
)

type LinkUseCase struct {
	submissions domain.SubmissionRepository
	links       domain.ShareLinkRepository
	cache       domain.LinkCache
	views       domain.ViewRecorder
	now         func() time.Time
}

func NewLinkUseCase(s domain.SubmissionRepository, l domain.ShareLinkRepository, c domain.LinkCache, views domain.ViewRecorder) *LinkUseCase {
	return &LinkUseCase{submissions: s, links: l, cache: c, views: views, now: time.Now}
}

type CreateLinkInput struct {
	ActorID, SubmissionID uuid.UUID
	AccessLevel           domain.AccessLevel
	ExpiresAt             *time.Time
}

func (u *LinkUseCase) Create(ctx context.Context, in CreateLinkInput) (domain.ShareLink, error) {
	s, err := u.submissions.FindByID(ctx, in.SubmissionID)
	if err != nil {
		return domain.ShareLink{}, err
	}
	if s.OwnerUserID != in.ActorID {
		return domain.ShareLink{}, domain.ErrForbidden
	}
	if in.AccessLevel == "" {
		in.AccessLevel = domain.AccessPublic
	}
	if !validAccess(in.AccessLevel) {
		return domain.ShareLink{}, domain.ErrConflict
	}
	for i := 0; i < 3; i++ {
		code, err := shortCode()
		if err != nil {
			return domain.ShareLink{}, err
		}
		link := domain.ShareLink{ID: uuid.New(), SubmissionID: s.ID, ShortCode: code, Status: domain.LinkStatusActive, AccessLevel: in.AccessLevel, ExpiresAt: in.ExpiresAt, CreatedBy: in.ActorID, CreatedAt: u.now(), UpdatedAt: u.now()}
		if err = u.links.CreateLink(ctx, link); err == nil {
			return link, nil
		}
		if err != domain.ErrConflict {
			return domain.ShareLink{}, err
		}
	}
	return domain.ShareLink{}, domain.ErrConflict
}

type UpdateLinkInput struct {
	ActorID, LinkID uuid.UUID
	AccessLevel     *domain.AccessLevel
	ExpiresAt       **time.Time
}

func (u *LinkUseCase) Update(ctx context.Context, in UpdateLinkInput) (domain.ShareLink, error) {
	link, err := u.links.FindByIDLink(ctx, in.LinkID)
	if err != nil {
		return link, err
	}
	if link.Submission.OwnerUserID != in.ActorID {
		return link, domain.ErrForbidden
	}
	if in.AccessLevel != nil {
		if !validAccess(*in.AccessLevel) {
			return link, domain.ErrConflict
		}
		link.AccessLevel = *in.AccessLevel
	}
	if in.ExpiresAt != nil {
		link.ExpiresAt = *in.ExpiresAt
	}
	if err = u.links.Update(ctx, link); err != nil {
		return link, err
	}
	_ = u.cache.Delete(ctx, key(link.ShortCode))
	return link, nil
}
func (u *LinkUseCase) Revoke(ctx context.Context, actor, id uuid.UUID) error {
	link, err := u.links.FindByIDLink(ctx, id)
	if err != nil {
		return err
	}
	if link.Submission.OwnerUserID != actor {
		return domain.ErrForbidden
	}
	if err = u.links.Revoke(ctx, id, u.now()); err != nil {
		return err
	}
	return u.cache.Delete(ctx, key(link.ShortCode))
}
func (u *LinkUseCase) Resolve(ctx context.Context, code string, viewer *uuid.UUID) (domain.ShareLink, error) {
	var link domain.ShareLink
	hit, err := u.cache.GetJSON(ctx, key(code), &link)
	if err != nil || !hit {
		link, err = u.links.FindByCode(ctx, code)
		if err != nil {
			return link, err
		}
		_ = u.cache.SetJSON(ctx, key(code), link)
	}
	if !link.IsAvailable(u.now(), viewer) {
		return link, domain.ErrNotFound
	}
	u.views.Record(link.ID, u.now())
	return link, nil
}
func (u *LinkUseCase) Analytics(ctx context.Context, actor, id uuid.UUID, from, to time.Time) (domain.ShareLink, []domain.DailyViews, error) {
	link, err := u.links.FindByIDLink(ctx, id)
	if err != nil {
		return link, nil, err
	}
	if link.Submission.OwnerUserID != actor {
		return link, nil, domain.ErrForbidden
	}
	rows, err := u.links.GetDailyViews(ctx, id, from, to)
	return link, rows, err
}
func key(code string) string { return "share-link:" + code }
func validAccess(a domain.AccessLevel) bool {
	return a == domain.AccessPublic || a == domain.AccessAuthenticated || a == domain.AccessRestricted
}
func shortCode() (string, error) {
	b := make([]byte, 9)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
