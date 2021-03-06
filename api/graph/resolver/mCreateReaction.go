package resolver

import (
	"context"
	"time"

	"github.com/romshark/dgraph_graphql_go/api/graph/auth"
	"github.com/romshark/dgraph_graphql_go/store"
	"github.com/romshark/dgraph_graphql_go/store/enum/emotion"
	strerr "github.com/romshark/dgraph_graphql_go/store/errors"
)

// CreateReaction resolves Mutation.createReaction
func (rsv *Resolver) CreateReaction(
	ctx context.Context,
	params struct {
		Author  string
		Subject string
		Emotion string
		Message string
	},
) *Reaction {
	if err := auth.Authorize(ctx, auth.IsOwner{
		Owner: store.ID(params.Author),
	}); err != nil {
		rsv.error(ctx, err)
		return nil
	}

	emot := emotion.Emotion(params.Emotion)

	// Validate input
	if err := rsv.validator.ReactionMessage(params.Message); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}
	if err := emotion.Validate(emot); err != nil {
		err = strerr.Wrap(strerr.ErrInvalidInput, err)
		rsv.error(ctx, err)
		return nil
	}

	creationTime := time.Now()

	// Create new reaction entity
	newReaction, err := rsv.str.CreateReaction(
		ctx,
		creationTime,
		store.ID(params.Author),
		store.ID(params.Subject),
		emot,
		params.Message,
	)
	if err != nil {
		rsv.error(ctx, err)
		return nil
	}

	return &Reaction{
		root:       rsv,
		uid:        newReaction.UID,
		id:         newReaction.ID,
		creation:   creationTime,
		authorUID:  newReaction.Author.UID,
		subjectUID: newReaction.Subject.NodeID(),
		emotion:    emot,
		message:    params.Message,
	}
}
