package decision

import "errors"

type Decision struct {
	ActorID        string
	RecipientID    string
	LikedRecipient bool
	UnixTimestamp  uint64
}

type Liker struct {
	ActorID       string
	UnixTimestamp uint64
}

type Cursor struct {
	LastUnixTimestamp uint64
	LastActorID       string
}

func (c Cursor) IsZero() bool {
	return c.LastUnixTimestamp == 0 && c.LastActorID == ""
}

var (
	ErrEmptyActorID     = errors.New("actor id must not be empty")
	ErrEmptyRecipientID = errors.New("recipient id must not be empty")
	ErrSameUser         = errors.New("actor and recipient must not be the same user")
	ErrNotFound         = errors.New("decision not found")
)

func NewDecision(actorID, recipientID string, liked bool, ts uint64) (*Decision, error) {
	if actorID == "" {
		return nil, ErrEmptyActorID
	}

	if recipientID == "" {
		return nil, ErrEmptyRecipientID
	}

	if actorID == recipientID {
		return nil, ErrSameUser
	}

	return &Decision{
		ActorID:        actorID,
		RecipientID:    recipientID,
		LikedRecipient: liked,
		UnixTimestamp:  ts,
	}, nil
}
