package like_tweet

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"ualaTwitter/internal/platform/logger"

	"ualaTwitter/internal/domain/like"
	"ualaTwitter/internal/domain/tweet"
	"ualaTwitter/internal/platform/errors/usecase"
)

type LikeTweetService struct {
	TweetRepo tweet.Repository
	LikeRepo  like.Repository
}

func NewLikeTweetService(tweetRepo tweet.Repository, likeRepo like.Repository) *LikeTweetService {
	return &LikeTweetService{
		TweetRepo: tweetRepo,
		LikeRepo:  likeRepo,
	}
}

func (s *LikeTweetService) Execute(ctx context.Context, input Input) error {
	if input.UserID == "" || input.TweetID == "" {
		return usecase.InvalidParam("user ID and tweet ID must not be empty", like.ErrInvalidInput)
	}

	alreadyLiked, err := s.LikeRepo.HasLiked(ctx, input.TweetID, input.UserID)
	if err != nil {
		return usecase.InternalServerError("failed to check like status", err)
	}
	if alreadyLiked {
		logger.Log.Warn("duplicate like attempt",
			zap.String("user_id", input.UserID),
			zap.String("tweet_id", input.TweetID),
		)
		return usecase.Forbidden("user has already liked this tweet", like.ErrAlreadyLiked)
	}

	if err := s.TweetRepo.IncrementLikes(ctx, input.TweetID); err != nil {
		switch {
		case errors.Is(err, tweet.ErrNotFound):
			return usecase.NotFound("tweet not found", err)
		default:
			return usecase.InternalServerError("failed to increment tweet like count", err)
		}
	}

	if err := s.LikeRepo.Like(ctx, input.TweetID, input.UserID); err != nil {
		return usecase.InternalServerError("failed to persist like", err)
	}

	return nil
}
