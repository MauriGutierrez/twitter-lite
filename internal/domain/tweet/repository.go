package tweet

import "context"

type Repository interface {
	Save(ctx context.Context, t Tweet) error
	GetByID(ctx context.Context, id string) (Tweet, error)
	FindTweetsAuthoredBy(ctx context.Context, userID string) ([]Tweet, error)
	IncrementLikes(ctx context.Context, tweetID string) error
}
