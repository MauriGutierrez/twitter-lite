package post_tweet

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"ualaTwitter/internal/domain/tweet"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/test/mocks"
)

func TestPostTweetService_Execute(t *testing.T) {
	validUser := &user.User{ID: "usr_123", Name: "Test User", Document: "12345678"}
	ctx := context.Background()

	t.Run("successfully posts a tweet", func(t *testing.T) {
		userRepo := &mocks.FakeUserRepo{Users: map[string]*user.User{validUser.ID: validUser}}
		tweetRepo := &mocks.FakeTweetRepo{}

		service := NewPostTweetService(tweetRepo, userRepo)

		input := Input{
			UserID:  validUser.ID,
			Content: "hello world",
		}

		tweetID, err := service.Execute(ctx, input)
		assert.NoError(t, err)
		assert.NotEmpty(t, tweetID)
		assert.NotNil(t, tweetRepo.Saved)
		assert.Equal(t, input.Content, tweetRepo.Saved.Content)
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := &mocks.FakeUserRepo{Users: map[string]*user.User{validUser.ID: validUser}}
		tweetRepo := &mocks.FakeTweetRepo{}

		service := NewPostTweetService(tweetRepo, userRepo)

		_, err := service.Execute(ctx, Input{
			UserID:  "nonexistent",
			Content: "hi",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not_found")
	})

	t.Run("empty tweet content", func(t *testing.T) {
		userRepo := &mocks.FakeUserRepo{Users: map[string]*user.User{validUser.ID: validUser}}
		tweetRepo := &mocks.FakeTweetRepo{}

		service := NewPostTweetService(tweetRepo, userRepo)

		_, err := service.Execute(ctx, Input{
			UserID:  validUser.ID,
			Content: "",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid tweet content")
	})

	t.Run("too long tweet content", func(t *testing.T) {
		userRepo := &mocks.FakeUserRepo{Users: map[string]*user.User{validUser.ID: validUser}}
		tweetRepo := &mocks.FakeTweetRepo{}

		service := NewPostTweetService(tweetRepo, userRepo)

		tooLong := make([]byte, tweet.MaxContentLength+1)
		for i := range tooLong {
			tooLong[i] = 'a'
		}

		_, err := service.Execute(ctx, Input{
			UserID:  validUser.ID,
			Content: string(tooLong),
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid tweet content")
	})

	t.Run("tweet repo fails to save", func(t *testing.T) {
		userRepo := &mocks.FakeUserRepo{Users: map[string]*user.User{validUser.ID: validUser}}
		tweetRepo := &mocks.FakeTweetRepo{SaveErr: errors.New("db failure")}

		service := NewPostTweetService(tweetRepo, userRepo)

		_, err := service.Execute(ctx, Input{
			UserID:  validUser.ID,
			Content: "valid tweet",
		})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to persist tweet")
	})
}
