package grpc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"explore-service/internal/domain/decision"
)

type cursorToken struct {
	Ts      uint64 `json:"ts"`
	ActorID string `json:"actor_id"`
}

func encodeCursorToToken(c decision.Cursor) (string, error) {
	if c.IsZero() {
		return "", nil
	}

	tok := cursorToken{
		Ts:      c.LastUnixTimestamp,
		ActorID: c.LastActorID,
	}

	data, err := json.Marshal(tok)
	if err != nil {
		return "", fmt.Errorf("failed to marshal cursor token: %w", err)
	}

	return base64.URLEncoding.EncodeToString(data), nil
}

func decodeTokenToCursor(token string) (decision.Cursor, error) {
	if token == "" {
		return decision.Cursor{}, nil
	}

	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return decision.Cursor{}, fmt.Errorf("failed to decode pagination token: %w", err)
	}

	var t cursorToken
	if err := json.Unmarshal(raw, &t); err != nil {
		return decision.Cursor{}, fmt.Errorf("failed to unmarshal pagination token: %w", err)
	}

	return decision.Cursor{
		LastUnixTimestamp: t.Ts,
		LastActorID:       t.ActorID,
	}, nil
}
