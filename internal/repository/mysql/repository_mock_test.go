package mysql

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"explore-service/internal/domain/decision"
)

type mockDB struct {
	execFn     func(query string, args ...any) (sql.Result, error)
	queryFn    func(query string, args ...any) (rows, error)
	queryRowFn func(query string, args ...any) rowScanner
}

func (m mockDB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return m.execFn(query, args...)
}

func (m mockDB) QueryContext(ctx context.Context, query string, args ...any) (rows, error) {
	return m.queryFn(query, args...)
}

func (m mockDB) QueryRowContext(ctx context.Context, query string, args ...any) rowScanner {
	return m.queryRowFn(query, args...)
}

type mockResult struct {
	rowsAffected int64
}

func (r mockResult) LastInsertId() (int64, error) { return 0, nil }
func (r mockResult) RowsAffected() (int64, error) { return r.rowsAffected, nil }

type mockRows struct {
	data [][]any
	i    int
	err  error
}

func (r *mockRows) Next() bool {
	if r.i >= len(r.data) {
		return false
	}
	r.i++
	return true
}

func (r *mockRows) Scan(dest ...any) error {
	if r.i == 0 || r.i > len(r.data) {
		return errors.New("scan called out of bounds")
	}
	row := r.data[r.i-1]
	if len(row) != len(dest) {
		return errors.New("scan dest mismatch")
	}
	for idx, v := range row {
		switch d := dest[idx].(type) {
		case *string:
			*d = v.(string)
		case *uint64:
			*d = v.(uint64)
		case *int8:
			*d = v.(int8)
		default:
			return errors.New("unsupported scan type")
		}
	}
	return nil
}

func (r *mockRows) Close() error { return nil }
func (r *mockRows) Err() error   { return r.err }

type mockRow struct {
	values []any
	err    error
}

func (r mockRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan dest mismatch")
	}
	for i, v := range r.values {
		switch d := dest[i].(type) {
		case *string:
			*d = v.(string)
		case *uint64:
			*d = v.(uint64)
		case *int8:
			*d = v.(int8)
		default:
			return errors.New("unsupported scan type")
		}
	}
	return nil
}

func TestPutDecisionExecutes(t *testing.T) {
	called := false
	db := mockDB{
		execFn: func(query string, args ...any) (sql.Result, error) {
			called = true
			if len(args) != 4 {
				t.Fatalf("expected 4 args, got %d", len(args))
			}
			return mockResult{rowsAffected: 1}, nil
		},
		queryFn:    nil,
		queryRowFn: nil,
	}

	repo := DecisionRepository{db: db}
	err := repo.PutDecision(context.Background(), &decision.Decision{
		ActorId:        "a1",
		RecipientId:    "r1",
		LikedRecipient: true,
		UnixTimestamp:  10,
	})
	if err != nil {
		t.Fatalf("PutDecision err: %v", err)
	}
	if !called {
		t.Fatalf("expected exec called")
	}
}

func TestGetDecisionFoundAndNotFound(t *testing.T) {
	db := mockDB{
		execFn: nil,
		queryFn: func(string, ...any) (rows, error) {
			return nil, nil
		},
		queryRowFn: func(query string, args ...any) rowScanner {
			if args[0] == "a1" {
				return mockRow{values: []any{"a1", "r1", int8(1), uint64(5)}}
			}
			return mockRow{err: sql.ErrNoRows}
		},
	}
	repo := DecisionRepository{db: db}

	got, err := repo.GetDecision(context.Background(), "a1", "r1")
	if err != nil {
		t.Fatalf("GetDecision err: %v", err)
	}
	if got.ActorId != "a1" || got.RecipientId != "r1" || !got.LikedRecipient || got.UnixTimestamp != 5 {
		t.Fatalf("unexpected decision: %+v", got)
	}

	_, err = repo.GetDecision(context.Background(), "x", "y")
	if err == nil || !errors.Is(err, decision.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestListLikedYouPagination(t *testing.T) {
	db := mockDB{
		queryFn: func(query string, args ...any) (rows, error) {
			if args[0] != "r1" || args[len(args)-1] != 3 {
				t.Fatalf("unexpected args: %+v", args)
			}
			return &mockRows{
				data: [][]any{
					{"u1", uint64(3)},
					{"u2", uint64(2)},
					{"u3", uint64(1)},
				},
			}, nil
		},
	}
	repo := DecisionRepository{db: db}

	likers, next, err := repo.ListLikedYou(context.Background(), "r1", decision.Cursor{}, 2)
	if err != nil {
		t.Fatalf("ListLikedYou err: %v", err)
	}
	if len(likers) != 2 || likers[0].ActorId != "u1" || likers[1].ActorId != "u2" {
		t.Fatalf("unexpected likers: %+v", likers)
	}
	if next.IsZero() {
		t.Fatalf("expected next cursor")
	}
	if next.LastActorId != "u2" || next.LastUnixTimestamp != 2 {
		t.Fatalf("unexpected next cursor: %+v", next)
	}
}

func TestListNewLikedYou(t *testing.T) {
	db := mockDB{
		queryFn: func(query string, args ...any) (rows, error) {
			if args[0] != "r1" || args[len(args)-1] != 3 {
				t.Fatalf("unexpected args: %+v", args)
			}
			return &mockRows{
				data: [][]any{
					{"u2", uint64(2)},
				},
			}, nil
		},
	}
	repo := DecisionRepository{db: db}

	likers, next, err := repo.ListNewLikedYou(context.Background(), "r1", decision.Cursor{}, 2)
	if err != nil {
		t.Fatalf("ListNewLikedYou err: %v", err)
	}
	if len(likers) != 1 || likers[0].ActorId != "u2" {
		t.Fatalf("unexpected likers: %+v", likers)
	}
	if !next.IsZero() {
		t.Fatalf("expected empty next cursor")
	}
}

func TestCountLikedYou(t *testing.T) {
	db := mockDB{
		queryRowFn: func(query string, args ...any) rowScanner {
			if args[0] != "r1" {
				t.Fatalf("unexpected args: %+v", args)
			}
			return mockRow{values: []any{uint64(7)}}
		},
	}
	repo := DecisionRepository{db: db}

	count, err := repo.CountLikedYou(context.Background(), "r1")
	if err != nil {
		t.Fatalf("CountLikedYou err: %v", err)
	}
	if count != 7 {
		t.Fatalf("expected 7, got %d", count)
	}
}
