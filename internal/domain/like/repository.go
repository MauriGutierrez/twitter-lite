package like

import "context"

type Repository interface {
	HasLiked(ctx context.Context, userID, tweetID string) (bool, error)
	Like(ctx context.Context, userID, tweetID string) error
}
