package decision

import (
	"context"
)

type Repository interface {
	PutDecision(ctx context.Context, d *Decision) error
	GetDecision(ctx context.Context, actorID, recipientID string) (*Decision, error)
	ListLikedYou(ctx context.Context, recipientID string, cursor Cursor, limit int) ([]Liker, Cursor, error)
	ListNewLikedYou(ctx context.Context, recipientID string, cursor Cursor, limit int) ([]Liker, Cursor, error)
	CountLikedYou(ctx context.Context, recipientID string) (uint64, error)
}
