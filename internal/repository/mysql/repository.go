package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"explore-service/internal/domain/decision"
)

type DecisionRepository struct {
	db dbExecutor
}

func NewDecisionRepository(db *sql.DB) *DecisionRepository {
	return &DecisionRepository{db: sqlDB{db: db}}
}

type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) rowScanner
}

type rows interface {
	Next() bool
	Scan(dest ...any) error
	Close() error
	Err() error
}

type rowScanner interface {
	Scan(dest ...any) error
}

type sqlDB struct {
	db *sql.DB
}

func (s sqlDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s sqlDB) QueryContext(ctx context.Context, query string, args ...any) (rows, error) {
	r, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s sqlDB) QueryRowContext(ctx context.Context, query string, args ...any) rowScanner {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (r *DecisionRepository) PutDecision(ctx context.Context, d *decision.Decision) error {
	query := `
		INSERT INTO decisions (
			actor_user_id,
			recipient_user_id,
			liked_recipient,
			decision_unix_ts
		) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			liked_recipient = VALUES(liked_recipient),
			decision_unix_ts = VALUES(decision_unix_ts)
	`

	_, err := r.db.ExecContext(ctx, query, d.ActorId, d.RecipientId, boolToTinyInt(d.LikedRecipient), d.UnixTimestamp)
	if err != nil {
		return fmt.Errorf("put decision: %w", err)
	}

	return nil
}

func (r *DecisionRepository) GetDecision(ctx context.Context, actorId, recipientId string) (*decision.Decision, error) {
	var (
		liked int8
		ts    uint64
	)

	query := `
		SELECT actor_user_id, recipient_user_id, liked_recipient, decision_unix_ts
		FROM decisions
		WHERE actor_user_id = ? AND recipient_user_id = ?
	`

	row := r.db.QueryRowContext(ctx, query, actorId, recipientId)
	if err := row.Scan(&actorId, &recipientId, &liked, &ts); err != nil {
		if err == sql.ErrNoRows {
			return nil, decision.ErrNotFound
		}
		return nil, fmt.Errorf("get decision: %w", err)
	}

	return &decision.Decision{
		ActorId:        actorId,
		RecipientId:    recipientId,
		LikedRecipient: tinyIntToBool(liked),
		UnixTimestamp:  ts,
	}, nil
}

func (r *DecisionRepository) ListLikedYou(ctx context.Context, recipientId string, cursor decision.Cursor, limit int) ([]decision.Liker, decision.Cursor, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT actor_user_id, decision_unix_ts
		FROM decisions
		WHERE recipient_user_id = ?
		  AND liked_recipient = 1
	`
	args := []any{recipientId}

	if !cursor.IsZero() {
		query += `
		  AND (
			decision_unix_ts < ?
			OR (decision_unix_ts = ? AND actor_user_id > ?)
		  )
		`
		args = append(args, cursor.LastUnixTimestamp, cursor.LastUnixTimestamp, cursor.LastActorId)
	}

	query += `
		ORDER BY decision_unix_ts DESC, actor_user_id
		LIMIT ?
	`
	args = append(args, limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, decision.Cursor{}, fmt.Errorf("list liked you: %w", err)
	}
	defer rows.Close()

	likers := make([]decision.Liker, 0, limit+1)
	for rows.Next() {
		var actorId string
		var ts uint64
		if err := rows.Scan(&actorId, &ts); err != nil {
			return nil, decision.Cursor{}, fmt.Errorf("list liked you scan: %w", err)
		}
		likers = append(likers, decision.Liker{ActorId: actorId, UnixTimestamp: ts})
	}

	if err := rows.Err(); err != nil {
		return nil, decision.Cursor{}, fmt.Errorf("list liked you rows: %w", err)
	}

	var nextCursor decision.Cursor
	if len(likers) > limit {
		last := likers[limit-1]
		likers = likers[:limit]
		nextCursor = decision.Cursor{
			LastUnixTimestamp: last.UnixTimestamp,
			LastActorId:       last.ActorId,
		}
	}

	return likers, nextCursor, nil
}

func (r *DecisionRepository) ListNewLikedYou(ctx context.Context, recipientId string, cursor decision.Cursor, limit int) ([]decision.Liker, decision.Cursor, error) {
	if limit <= 0 {
		limit = 20
	}

	query := `
		SELECT d.actor_user_id, d.decision_unix_ts
		FROM decisions d
		WHERE d.recipient_user_id = ?
		  AND d.liked_recipient = 1
		  AND NOT EXISTS (
			SELECT 1
			FROM decisions r2
			WHERE r2.actor_user_id = d.recipient_user_id
			  AND r2.recipient_user_id = d.actor_user_id
			  AND r2.liked_recipient = 1
		  )
	`
	args := []any{recipientId}

	if !cursor.IsZero() {
		query += `
		  AND (
			d.decision_unix_ts < ?
			OR (d.decision_unix_ts = ? AND d.actor_user_id > ?)
		  )
		`
		args = append(args, cursor.LastUnixTimestamp, cursor.LastUnixTimestamp, cursor.LastActorId)
	}

	query += `
		ORDER BY d.decision_unix_ts DESC, d.actor_user_id
		LIMIT ?
	`
	args = append(args, limit+1)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, decision.Cursor{}, fmt.Errorf("list new liked you: %w", err)
	}

	defer rows.Close()

	likers := make([]decision.Liker, 0, limit+1)
	for rows.Next() {
		var actorId string
		var ts uint64
		if err := rows.Scan(&actorId, &ts); err != nil {
			return nil, decision.Cursor{}, fmt.Errorf("list new liked you scan: %w", err)
		}
		likers = append(likers, decision.Liker{ActorId: actorId, UnixTimestamp: ts})
	}

	if err := rows.Err(); err != nil {
		return nil, decision.Cursor{}, fmt.Errorf("list new liked you rows: %w", err)
	}

	var nextCursor decision.Cursor
	if len(likers) > limit {
		last := likers[limit-1]
		likers = likers[:limit]
		nextCursor = decision.Cursor{
			LastUnixTimestamp: last.UnixTimestamp,
			LastActorId:       last.ActorId,
		}
	}

	return likers, nextCursor, nil
}

func (r *DecisionRepository) CountLikedYou(ctx context.Context, recipientId string) (uint64, error) {
	var count uint64

	query := `
		SELECT COUNT(*)
		FROM decisions
		WHERE recipient_user_id = ?
		  AND liked_recipient = 1
	`

	if err := r.db.QueryRowContext(ctx, query, recipientId).Scan(&count); err != nil {
		return 0, fmt.Errorf("count liked you: %w", err)
	}
	return count, nil
}

func boolToTinyInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func tinyIntToBool(v int8) bool {
	return v != 0
}
