package grpc

import (
	"context"
	"time"

	"explore-service/internal/domain/decision"
	explorepb "explore-service/pkg/proto/explore/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ExploreHandler struct {
	explorepb.UnimplementedExploreServiceServer
	svc *decision.Service
}

func NewExploreHandler(svc *decision.Service) *ExploreHandler {
	return &ExploreHandler{svc: svc}
}

func (h *ExploreHandler) PutDecision(
	ctx context.Context,
	req *explorepb.PutDecisionRequest,
) (*explorepb.PutDecisionResponse, error) {
	actorID := req.GetActorUserId()
	recipientID := req.GetRecipientUserId()
	liked := req.GetLikedRecipient()

	if actorID == "" || recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "actor_user_id and recipient_user_id must not be empty")
	}

	ts := uint64(time.Now().Unix())

	mutual, err := h.svc.PutDecision(ctx, actorID, recipientID, liked, ts)
	if err != nil {
		switch err {
		case decision.ErrEmptyActorID, decision.ErrEmptyRecipientID, decision.ErrSameUser:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "failed to put decision: %v", err)
		}
	}

	return &explorepb.PutDecisionResponse{
		MutualLikes: mutual,
	}, nil
}

func (h *ExploreHandler) ListLikedYou(
	ctx context.Context,
	req *explorepb.ListLikedYouRequest,
) (*explorepb.ListLikedYouResponse, error) {
	recipientID := req.GetRecipientUserId()
	if recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient_user_id must not be empty")
	}

	cursor, err := decodeTokenToCursor(req.GetPaginationToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid pagination_token")
	}

	pageSize := int(req.GetPageSize())
	likers, nextCursor, err := h.svc.ListLikedYou(ctx, recipientID, cursor, pageSize)
	if err != nil {
		if err == decision.ErrEmptyRecipientID {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to list liked you: %v", err)
	}

	resp := &explorepb.ListLikedYouResponse{
		Likers: make([]*explorepb.ListLikedYouResponse_Liker, 0, len(likers)),
	}

	for _, l := range likers {
		resp.Likers = append(resp.Likers, &explorepb.ListLikedYouResponse_Liker{
			ActorId:       l.ActorID,
			UnixTimestamp: l.UnixTimestamp,
		})
	}

	if !nextCursor.IsZero() {
		if token, err := encodeCursorToToken(nextCursor); err == nil {
			resp.NextPaginationToken = &token
		}
	}

	return resp, nil
}

func (h *ExploreHandler) ListNewLikedYou(
	ctx context.Context,
	req *explorepb.ListLikedYouRequest,
) (*explorepb.ListLikedYouResponse, error) {
	recipientID := req.GetRecipientUserId()
	if recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient_user_id must not be empty")
	}

	cursor, err := decodeTokenToCursor(req.GetPaginationToken())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid pagination_token")
	}

	pageSize := int(req.GetPageSize())
	likers, nextCursor, err := h.svc.ListNewLikedYou(ctx, recipientID, cursor, pageSize)
	if err != nil {
		if err == decision.ErrEmptyRecipientID {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to list new liked you: %v", err)
	}

	resp := &explorepb.ListLikedYouResponse{
		Likers: make([]*explorepb.ListLikedYouResponse_Liker, 0, len(likers)),
	}

	for _, l := range likers {
		resp.Likers = append(resp.Likers, &explorepb.ListLikedYouResponse_Liker{
			ActorId:       l.ActorID,
			UnixTimestamp: l.UnixTimestamp,
		})
	}

	if !nextCursor.IsZero() {
		if token, err := encodeCursorToToken(nextCursor); err == nil {
			resp.NextPaginationToken = &token
		}
	}

	return resp, nil
}

func (h *ExploreHandler) CountLikedYou(
	ctx context.Context,
	req *explorepb.CountLikedYouRequest,
) (*explorepb.CountLikedYouResponse, error) {
	recipientID := req.GetRecipientUserId()
	if recipientID == "" {
		return nil, status.Error(codes.InvalidArgument, "recipient_user_id must not be empty")
	}

	count, err := h.svc.CountLikedYou(ctx, recipientID)
	if err != nil {
		if err == decision.ErrEmptyRecipientID {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to count liked you: %v", err)
	}

	return &explorepb.CountLikedYouResponse{
		Count: count,
	}, nil
}
