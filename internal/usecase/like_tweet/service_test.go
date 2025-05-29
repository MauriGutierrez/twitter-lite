package like_tweet

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"
	"ualaTwitter/internal/domain/tweet"
	"ualaTwitter/internal/platform/logger"
	"ualaTwitter/internal/test/mocks"

	"github.com/stretchr/testify/assert"
)

func init() {
	logger.Log = zap.NewNop()
}

func TestLikeTweetService_Execute(t *testing.T) {
	ctx := context.Background()
	validUserID := "usr_123"
	validTweetID := "tweet_456"

	t.Run("successfully likes a tweet", func(t *testing.T) {
		likeRepo := &mocks.FakeLikeRepo{}
		tweetRepo := &mocks.FakeTweetRepo{}

		service := NewLikeTweetService(tweetRepo, likeRepo)

		err := service.Execute(ctx, Input{
			UserID:  validUserID,
			TweetID: validTweetID,
		})

		assert.NoError(t, err)
		assert.Equal(t, validTweetID, likeRepo.LikedTweetID)
		assert.Equal(t, validUserID, likeRepo.LikedUserID)
		assert.Equal(t, validTweetID, tweetRepo.LastLikedTweetID)
	})

	t.Run("missing user or tweet ID", func(t *testing.T) {
		service := NewLikeTweetService(&mocks.FakeTweetRepo{}, &mocks.FakeLikeRepo{})

		err := service.Execute(ctx, Input{
			UserID:  "",
			TweetID: "123",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID and tweet ID")
	})

	t.Run("already liked tweet", func(t *testing.T) {
		likeRepo := &mocks.FakeLikeRepo{AlreadyLiked: true}
		service := NewLikeTweetService(&mocks.FakeTweetRepo{}, likeRepo)

		err := service.Execute(ctx, Input{
			UserID:  validUserID,
			TweetID: validTweetID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already liked")
	})

	t.Run("error checking like status", func(t *testing.T) {
		likeRepo := &mocks.FakeLikeRepo{HasLikedErr: errors.New("db error")}
		service := NewLikeTweetService(&mocks.FakeTweetRepo{}, likeRepo)

		err := service.Execute(ctx, Input{
			UserID:  validUserID,
			TweetID: validTweetID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "check like status")
	})

	t.Run("tweet not found on increment", func(t *testing.T) {
		tweetRepo := &mocks.FakeTweetRepo{IncrementLikesErr: tweet.ErrNotFound}
		service := NewLikeTweetService(tweetRepo, &mocks.FakeLikeRepo{})

		err := service.Execute(ctx, Input{
			UserID:  validUserID,
			TweetID: validTweetID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("error persisting like", func(t *testing.T) {
		likeRepo := &mocks.FakeLikeRepo{LikeErr: errors.New("db error")}
		tweetRepo := &mocks.FakeTweetRepo{}
		service := NewLikeTweetService(tweetRepo, likeRepo)

		err := service.Execute(ctx, Input{
			UserID:  validUserID,
			TweetID: validTweetID,
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "persist like")
	})
}
