package decision

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"
)

type stubRepo struct {
	mu               sync.Mutex
	decisions        map[string]*Decision
	lastListLikedLim int
	lastListNewLim   int
	countCalls       int
}

func newStubRepo() *stubRepo {
	return &stubRepo{
		decisions: make(map[string]*Decision),
	}
}

func (r *stubRepo) key(actor, recipient string) string {
	return actor + "|" + recipient
}

func (r *stubRepo) PutDecision(_ context.Context, d *Decision) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *d
	r.decisions[r.key(d.ActorId, d.RecipientId)] = &cp
	return nil
}

func (r *stubRepo) GetDecision(_ context.Context, actorId, recipientId string) (*Decision, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if dec, ok := r.decisions[r.key(actorId, recipientId)]; ok {
		cp := *dec
		return &cp, nil
	}
	return nil, ErrNotFound
}

func (r *stubRepo) ListLikedYou(ctx context.Context, recipientId string, _ Cursor, limit int) ([]Liker, Cursor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastListLikedLim = limit

	likers := make([]Liker, 0)
	for _, d := range r.decisions {
		if d.RecipientId == recipientId && d.LikedRecipient {
			likers = append(likers, Liker{ActorId: d.ActorId, UnixTimestamp: d.UnixTimestamp})
		}
	}

	sort.Slice(likers, func(i, j int) bool {
		if likers[i].UnixTimestamp == likers[j].UnixTimestamp {
			return likers[i].ActorId < likers[j].ActorId
		}
		return likers[i].UnixTimestamp > likers[j].UnixTimestamp
	})

	if len(likers) > limit {
		likers = likers[:limit]
	}
	return likers, Cursor{}, nil
}

func (r *stubRepo) ListNewLikedYou(ctx context.Context, recipientId string, _ Cursor, limit int) ([]Liker, Cursor, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastListNewLim = limit

	likers := make([]Liker, 0)
	for _, d := range r.decisions {
		if d.RecipientId != recipientId || !d.LikedRecipient {
			continue
		}
		if other, ok := r.decisions[r.key(recipientId, d.ActorId)]; ok && other.LikedRecipient {
			continue
		}
		likers = append(likers, Liker{ActorId: d.ActorId, UnixTimestamp: d.UnixTimestamp})
	}

	sort.Slice(likers, func(i, j int) bool {
		if likers[i].UnixTimestamp == likers[j].UnixTimestamp {
			return likers[i].ActorId < likers[j].ActorId
		}
		return likers[i].UnixTimestamp > likers[j].UnixTimestamp
	})

	if len(likers) > limit {
		likers = likers[:limit]
	}
	return likers, Cursor{}, nil
}

func (r *stubRepo) CountLikedYou(ctx context.Context, recipientId string) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.countCalls++

	var count uint64
	for _, d := range r.decisions {
		if d.RecipientId == recipientId && d.LikedRecipient {
			count++
		}
	}
	return count, nil
}

type stubCache struct {
	mu          sync.Mutex
	values      map[string]uint64
	invalidated map[string]bool
}

func newStubCache() *stubCache {
	return &stubCache{
		values:      make(map[string]uint64),
		invalidated: make(map[string]bool),
	}
}

func (c *stubCache) GetCount(ctx context.Context, recipientId string) (uint64, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.values[recipientId]
	return val, ok, nil
}

func (c *stubCache) SetCount(ctx context.Context, recipientId string, count uint64, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[recipientId] = count
	return nil
}

func (c *stubCache) InvalidateCount(ctx context.Context, recipientId string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.values, recipientId)
	c.invalidated[recipientId] = true
	return nil
}

func TestPutDecisionMutualAndCacheInvalidation(t *testing.T) {
	ctx := context.Background()
	repo := newStubRepo()
	cache := newStubCache()
	svc := NewService(repo, cache)

	mutual, err := svc.PutDecision(ctx, "a1", "b1", true, 1)
	if err != nil {
		t.Fatalf("PutDecision first: %v", err)
	}
	if mutual {
		t.Fatalf("expected mutual=false on first like")
	}
	if !cache.invalidated["b1"] {
		t.Fatalf("expected cache invalidation for recipient b1")
	}

	mutual, err = svc.PutDecision(ctx, "b1", "a1", true, 2)
	if err != nil {
		t.Fatalf("PutDecision second: %v", err)
	}
	if !mutual {
		t.Fatalf("expected mutual=true after reciprocal like")
	}
	if !cache.invalidated["a1"] {
		t.Fatalf("expected cache invalidation for recipient a1")
	}
}

func TestCountLikedYouCacheHit(t *testing.T) {
	ctx := context.Background()
	repo := newStubRepo()
	cache := newStubCache()
	cache.values["r1"] = 9

	svc := NewService(repo, cache)
	count, err := svc.CountLikedYou(ctx, "r1")
	if err != nil {
		t.Fatalf("CountLikedYou: %v", err)
	}
	if count != 9 {
		t.Fatalf("expected cached count 9, got %d", count)
	}
	if repo.countCalls != 0 {
		t.Fatalf("expected repo count not called on cache hit")
	}
}

func TestCountLikedYouCacheMissSetsCache(t *testing.T) {
	ctx := context.Background()
	repo := newStubRepo()
	_ = repo.PutDecision(ctx, &Decision{ActorId: "x1", RecipientId: "r2", LikedRecipient: true, UnixTimestamp: 1})
	cache := newStubCache()

	svc := NewService(repo, cache)
	count, err := svc.CountLikedYou(ctx, "r2")
	if err != nil {
		t.Fatalf("CountLikedYou: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected count 1, got %d", count)
	}
	if repo.countCalls != 1 {
		t.Fatalf("expected repo count called once")
	}
	if cache.values["r2"] != 1 {
		t.Fatalf("expected cache set to 1, got %d", cache.values["r2"])
	}
}

func TestListLimitCapped(t *testing.T) {
	ctx := context.Background()
	repo := newStubRepo()
	cache := newStubCache()
	svc := NewService(repo, cache)

	_, _, err := svc.ListLikedYou(ctx, "r3", Cursor{}, 500)
	if err != nil {
		t.Fatalf("ListLikedYou: %v", err)
	}
	if repo.lastListLikedLim != maxPageSize {
		t.Fatalf("expected limit capped to %d, got %d", maxPageSize, repo.lastListLikedLim)
	}

	_, _, err = svc.ListNewLikedYou(ctx, "r3", Cursor{}, 500)
	if err != nil {
		t.Fatalf("ListNewLikedYou: %v", err)
	}
	if repo.lastListNewLim != maxPageSize {
		t.Fatalf("expected limit capped to %d, got %d", maxPageSize, repo.lastListNewLim)
	}
}
