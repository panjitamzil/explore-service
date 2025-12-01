package memory

import (
	"context"
	"sort"
	"sync"

	"explore-service/internal/domain/decision"
)

type DecisionRepository struct {
	mu        sync.RWMutex
	decisions map[string]*decision.Decision
}

func NewDecisionRepository() *DecisionRepository {
	return &DecisionRepository{
		decisions: make(map[string]*decision.Decision),
	}
}

func (r *DecisionRepository) key(actor, recipient string) string {
	return actor + "|" + recipient
}

func (r *DecisionRepository) PutDecision(_ context.Context, d *decision.Decision) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cp := *d
	r.decisions[r.key(d.ActorId, d.RecipientId)] = &cp
	return nil
}

func (r *DecisionRepository) GetDecision(ctx context.Context, actorId, recipientId string) (*decision.Decision, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if dec, ok := r.decisions[r.key(actorId, recipientId)]; ok {
		cp := *dec
		return &cp, nil
	}
	return nil, decision.ErrNotFound
}

func (r *DecisionRepository) ListLikedYou(ctx context.Context, recipientId string, cursor decision.Cursor, limit int) ([]decision.Liker, decision.Cursor, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	likers := make([]decision.Liker, 0, len(r.decisions))
	for _, d := range r.decisions {
		if d.RecipientId == recipientId && d.LikedRecipient {
			likers = append(likers, decision.Liker{
				ActorId:       d.ActorId,
				UnixTimestamp: d.UnixTimestamp,
			})
		}
	}

	likers, next := paginate(likers, cursor, limit)
	return likers, next, nil
}

func (r *DecisionRepository) ListNewLikedYou(ctx context.Context, recipientId string, cursor decision.Cursor, limit int) ([]decision.Liker, decision.Cursor, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	likers := make([]decision.Liker, 0, len(r.decisions))
	for _, d := range r.decisions {
		if d.RecipientId != recipientId || !d.LikedRecipient {
			continue
		}
		if other, ok := r.decisions[r.key(recipientId, d.ActorId)]; ok && other.LikedRecipient {
			continue
		}
		likers = append(likers, decision.Liker{
			ActorId:       d.ActorId,
			UnixTimestamp: d.UnixTimestamp,
		})
	}

	likers, next := paginate(likers, cursor, limit)
	return likers, next, nil
}

func (r *DecisionRepository) CountLikedYou(ctx context.Context, recipientId string) (uint64, error) {
	_ = ctx

	r.mu.RLock()
	defer r.mu.RUnlock()

	var count uint64
	for _, d := range r.decisions {
		if d.RecipientId == recipientId && d.LikedRecipient {
			count++
		}
	}
	return count, nil
}

func paginate(likers []decision.Liker, cursor decision.Cursor, limit int) ([]decision.Liker, decision.Cursor) {
	sort.Slice(likers, func(i, j int) bool {
		if likers[i].UnixTimestamp == likers[j].UnixTimestamp {
			return likers[i].ActorId < likers[j].ActorId
		}
		return likers[i].UnixTimestamp > likers[j].UnixTimestamp
	})

	filtered := likers
	if !cursor.IsZero() {
		filtered = make([]decision.Liker, 0, len(likers))
		for _, l := range likers {
			if l.UnixTimestamp < cursor.LastUnixTimestamp {
				filtered = append(filtered, l)
				continue
			}
			if l.UnixTimestamp == cursor.LastUnixTimestamp && l.ActorId > cursor.LastActorId {
				filtered = append(filtered, l)
			}
		}
	}

	var next decision.Cursor
	if len(filtered) > limit {
		next = decision.Cursor{
			LastUnixTimestamp: filtered[limit-1].UnixTimestamp,
			LastActorId:       filtered[limit-1].ActorId,
		}
		filtered = filtered[:limit]
	}

	return filtered, next
}
