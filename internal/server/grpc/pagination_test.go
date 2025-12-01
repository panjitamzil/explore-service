package grpc

import (
	"testing"

	"explore-service/internal/domain/decision"
)

func TestEncodeDecodeCursorRoundTrip(t *testing.T) {
	orig := decision.Cursor{LastUnixTimestamp: 123, LastActorId: "user-1"}
	token, err := encodeCursorToToken(orig)
	if err != nil {
		t.Fatalf("encode cursor: %v", err)
	}
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	cur, err := decodeTokenToCursor(token)
	if err != nil {
		t.Fatalf("decode cursor: %v", err)
	}

	if cur.LastUnixTimestamp != orig.LastUnixTimestamp || cur.LastActorId != orig.LastActorId {
		t.Fatalf("round trip mismatch: got %+v want %+v", cur, orig)
	}
}

func TestDecodeEmptyToken(t *testing.T) {
	cur, err := decodeTokenToCursor("")
	if err != nil {
		t.Fatalf("decode empty token error: %v", err)
	}
	if !cur.IsZero() {
		t.Fatalf("expected zero cursor for empty token")
	}
}

func TestDecodeInvalidToken(t *testing.T) {
	if _, err := decodeTokenToCursor("not-base64"); err == nil {
		t.Fatalf("expected error for invalid token")
	}
}
