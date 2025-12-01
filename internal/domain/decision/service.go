package decision

import (
	"context"
	"time"
)

type LikedCountCache interface {
	GetCount(ctx context.Context, recipientID string) (count uint64, found bool, err error)
	SetCount(ctx context.Context, recipientID string, count uint64, ttl time.Duration) error
	InvalidateCount(ctx context.Context, recipientID string) error
}

type Service struct {
	repo          Repository
	cache         LikedCountCache
	likedCountTTL time.Duration
}

const (
	defaultPageSize = 20
	maxPageSize     = 100
)

func NewService(repo Repository, cache LikedCountCache) *Service {
	return &Service{
		repo:          repo,
		cache:         cache,
		likedCountTTL: 5 * time.Minute,
	}
}

func (s *Service) PutDecision(ctx context.Context, actorID, recipientID string, liked bool, ts uint64) (bool, error) {
	dec, err := NewDecision(actorID, recipientID, liked, ts)
	if err != nil {
		return false, err
	}

	if err := s.repo.PutDecision(ctx, dec); err != nil {
		return false, err
	}

	if s.cache != nil {
		_ = s.cache.InvalidateCount(ctx, recipientID)
	}

	opposite, err := s.repo.GetDecision(ctx, recipientID, actorID)
	if err != nil {
		if err == ErrNotFound {
			return false, nil
		}
		return false, err
	}

	mutual := liked && opposite.LikedRecipient
	return mutual, nil
}

func (s *Service) ListLikedYou(ctx context.Context, recipientID string, cursor Cursor, limit int) ([]Liker, Cursor, error) {
	if recipientID == "" {
		return nil, Cursor{}, ErrEmptyRecipientId
	}
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	likers, nextCursor, err := s.repo.ListLikedYou(ctx, recipientID, cursor, limit)
	if err != nil {
		return nil, Cursor{}, err
	}
	return likers, nextCursor, nil
}

func (s *Service) ListNewLikedYou(ctx context.Context, recipientID string, cursor Cursor, limit int) ([]Liker, Cursor, error) {
	if recipientID == "" {
		return nil, Cursor{}, ErrEmptyRecipientId
	}
	if limit <= 0 {
		limit = defaultPageSize
	}
	if limit > maxPageSize {
		limit = maxPageSize
	}

	likers, nextCursor, err := s.repo.ListNewLikedYou(ctx, recipientID, cursor, limit)
	if err != nil {
		return nil, Cursor{}, err
	}
	return likers, nextCursor, nil
}

func (s *Service) CountLikedYou(ctx context.Context, recipientID string) (uint64, error) {
	if recipientID == "" {
		return 0, ErrEmptyRecipientId
	}

	if s.cache != nil {
		if count, found, err := s.cache.GetCount(ctx, recipientID); err == nil && found {
			return count, nil
		}
	}

	count, err := s.repo.CountLikedYou(ctx, recipientID)
	if err != nil {
		return 0, err
	}

	if s.cache != nil {
		_ = s.cache.SetCount(ctx, recipientID, count, s.likedCountTTL)
	}

	return count, nil
}
