package memory

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"ualaTwitter/internal/domain/tweet"
)

func makeMockTweet(userID, content string, createdAt time.Time, likes int) tweet.Tweet {
	return tweet.Tweet{
		ID:        userID + "-" + content,
		UserID:    userID,
		Content:   content,
		CreatedAt: createdAt,
		Likes:     likes,
	}
}

func TestInMemoryTweetRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryTweetRepository()
	userID := "usr1"

	t.Run("Save and GetByID works", func(t *testing.T) {
		tw := makeMockTweet(userID, "Hello world", time.Now(), 0)
		assert.NoError(t, repo.Save(ctx, tw))
		got, err := repo.GetByID(ctx, tw.ID)
		assert.NoError(t, err)
		assert.Equal(t, tw.Content, got.Content)
		assert.Equal(t, tw.UserID, got.UserID)
	})

	t.Run("GetByID returns error for unknown tweet", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "no-such-id")
		assert.Error(t, err)
	})

	t.Run("FindTweetsAuthoredBy returns all tweets for user", func(t *testing.T) {
		repo := NewInMemoryTweetRepository()
		tw1 := makeMockTweet(userID, "First", time.Now(), 0)
		tw2 := makeMockTweet(userID, "Second", time.Now().Add(1*time.Minute), 0)
		assert.NoError(t, repo.Save(ctx, tw1))
		assert.NoError(t, repo.Save(ctx, tw2))

		tweets, err := repo.FindTweetsAuthoredBy(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, tweets, 2)
		assert.Contains(t, []string{tweets[0].Content, tweets[1].Content}, "First")
		assert.Contains(t, []string{tweets[0].Content, tweets[1].Content}, "Second")
	})

	t.Run("FindTweetsAuthoredBy returns error for unknown user", func(t *testing.T) {
		_, err := repo.FindTweetsAuthoredBy(ctx, "ghost")
		assert.Error(t, err)
	})

	t.Run("IncrementLikes increments like count in all views", func(t *testing.T) {
		repo := NewInMemoryTweetRepository()
		tw := makeMockTweet(userID, "Like me", time.Now(), 0)
		assert.NoError(t, repo.Save(ctx, tw))

		assert.NoError(t, repo.IncrementLikes(ctx, tw.ID))
		got, _ := repo.GetByID(ctx, tw.ID)
		assert.Equal(t, 1, got.Likes)

		userTweets, _ := repo.FindTweetsAuthoredBy(ctx, userID)
		found := false
		for _, ut := range userTweets {
			if ut.ID == tw.ID {
				assert.Equal(t, 1, ut.Likes)
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("IncrementLikes returns error for unknown tweet", func(t *testing.T) {
		err := repo.IncrementLikes(ctx, "no-such-tweet")
		assert.Error(t, err)
	})
}
