package decision

import (
	"context"
)

type Repository interface {
	PutDecision(ctx context.Context, d *Decision) error
	GetDecision(ctx context.Context, actorId, recipientId string) (*Decision, error)
	ListLikedYou(ctx context.Context, recipientId string, cursor Cursor, limit int) ([]Liker, Cursor, error)
	ListNewLikedYou(ctx context.Context, recipientId string, cursor Cursor, limit int) ([]Liker, Cursor, error)
	CountLikedYou(ctx context.Context, recipientId string) (uint64, error)
}
