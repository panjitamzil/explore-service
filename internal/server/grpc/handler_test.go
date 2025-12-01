package grpc

import (
	"context"
	"testing"

	"explore-service/internal/domain/decision"
	memrepo "explore-service/internal/repository/memory"
	explorepb "explore-service/pkg/proto/explore/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func newTestHandler() (*ExploreHandler, *memrepo.DecisionRepository) {
	repo := memrepo.NewDecisionRepository()
	svc := decision.NewService(repo, nil)
	return NewExploreHandler(svc), repo
}

func TestPutDecisionInvalidArgument(t *testing.T) {
	h, _ := newTestHandler()
	_, err := h.PutDecision(context.Background(), &explorepb.PutDecisionRequest{
		ActorUserId:     "",
		RecipientUserId: "r1",
		LikedRecipient:  true,
	})
	if err == nil || status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected invalid argument, got %v", err)
	}
}

func TestPutDecisionMutual(t *testing.T) {
	h, repo := newTestHandler()
	ctx := context.Background()

	_ = repo.PutDecision(ctx, &decision.Decision{
		ActorId:        "b",
		RecipientId:    "a",
		LikedRecipient: true,
		UnixTimestamp:  1,
	})

	resp, err := h.PutDecision(ctx, &explorepb.PutDecisionRequest{
		ActorUserId:     "a",
		RecipientUserId: "b",
		LikedRecipient:  true,
	})
	if err != nil {
		t.Fatalf("PutDecision error: %v", err)
	}
	if !resp.MutualLikes {
		t.Fatalf("expected mutual likes true")
	}
}

func TestListLikedYouInvalidToken(t *testing.T) {
	h, _ := newTestHandler()
	_, err := h.ListLikedYou(context.Background(), &explorepb.ListLikedYouRequest{
		RecipientUserId: "r1",
		PaginationToken: strPtr("not-base64"),
	})
	if err == nil || status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected invalid argument for bad token, got %v", err)
	}
}

func TestListLikedYouPaginationFlow(t *testing.T) {
	h, repo := newTestHandler()
	ctx := context.Background()

	_ = repo.PutDecision(ctx, &decision.Decision{ActorId: "u1", RecipientId: "r1", LikedRecipient: true, UnixTimestamp: 3})
	_ = repo.PutDecision(ctx, &decision.Decision{ActorId: "u2", RecipientId: "r1", LikedRecipient: true, UnixTimestamp: 2})

	first, err := h.ListLikedYou(ctx, &explorepb.ListLikedYouRequest{
		RecipientUserId: "r1",
		PageSize:        uint32Ptr(1),
	})
	if err != nil {
		t.Fatalf("ListLikedYou first page error: %v", err)
	}
	if len(first.Likers) != 1 || first.Likers[0].ActorId != "u1" {
		t.Fatalf("unexpected first page likers: %+v", first.Likers)
	}
	if first.NextPaginationToken == nil || *first.NextPaginationToken == "" {
		t.Fatalf("expected next pagination token")
	}

	second, err := h.ListLikedYou(ctx, &explorepb.ListLikedYouRequest{
		RecipientUserId: "r1",
		PaginationToken: first.NextPaginationToken,
		PageSize:        uint32Ptr(1),
	})
	if err != nil {
		t.Fatalf("ListLikedYou second page error: %v", err)
	}
	if len(second.Likers) != 1 || second.Likers[0].ActorId != "u2" {
		t.Fatalf("unexpected second page likers: %+v", second.Likers)
	}
}

func TestCountLikedYouValidationAndSuccess(t *testing.T) {
	h, repo := newTestHandler()
	ctx := context.Background()

	_, err := h.CountLikedYou(ctx, &explorepb.CountLikedYouRequest{})
	if err == nil || status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected invalid argument on empty recipient, got %v", err)
	}

	_ = repo.PutDecision(ctx, &decision.Decision{ActorId: "u1", RecipientId: "r9", LikedRecipient: true, UnixTimestamp: 1})
	_ = repo.PutDecision(ctx, &decision.Decision{ActorId: "u2", RecipientId: "r9", LikedRecipient: true, UnixTimestamp: 2})

	resp, err := h.CountLikedYou(ctx, &explorepb.CountLikedYouRequest{RecipientUserId: "r9"})
	if err != nil {
		t.Fatalf("CountLikedYou error: %v", err)
	}
	if resp.Count != 2 {
		t.Fatalf("expected count 2, got %d", resp.Count)
	}
}

func strPtr(s string) *string {
	return &s
}

func uint32Ptr(v uint32) *uint32 {
	return &v
}
