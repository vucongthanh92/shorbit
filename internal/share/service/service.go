package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/vucongthanh92/shortbit/internal/platform/cache"
	"github.com/vucongthanh92/shortbit/internal/share/domain"
	"github.com/vucongthanh92/shortbit/internal/share/repository"
)

var ErrUnavailable = errors.New("share link is unavailable")

type Service struct {
	repo  *repository.Postgres
	cache *cache.Redis
	now   func() time.Time
}

func New(repo *repository.Postgres, cache *cache.Redis) *Service {
	return &Service{repo: repo, cache: cache, now: time.Now}
}
func cacheKey(code string) string { return "share-link:" + code }
func (s *Service) Resolve(ctx context.Context, code string) (domain.ShareLink, error) {
	var link domain.ShareLink
	if hit, err := s.cache.GetJSON(ctx, cacheKey(code), &link); err == nil && hit {
		if link.IsAccessible(s.now()) {
			return link, nil
		}
		return link, ErrUnavailable
	}
	link, err := s.repo.FindByCode(ctx, code)
	if err != nil {
		return link, err
	}
	if !link.IsAccessible(s.now()) {
		return link, ErrUnavailable
	}
	// Redis failure must not block a valid public share link.
	_ = s.cache.SetJSON(ctx, cacheKey(code), link)
	return link, nil
}
func GenerateShortCode() (string, error) {
	raw := make([]byte, 9)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("random short code: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
