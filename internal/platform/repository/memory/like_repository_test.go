package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryLikeRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryLikeRepository()

	t.Run("like a tweet and check HasLiked returns true", func(t *testing.T) {
		userID, tweetID := "user1", "tweet1"
		assert.NoError(t, repo.Like(ctx, userID, tweetID))

		liked, err := repo.HasLiked(ctx, userID, tweetID)
		assert.NoError(t, err)
		assert.True(t, liked)
	})

	t.Run("HasLiked returns false if tweet was not liked", func(t *testing.T) {
		userID, tweetID := "user2", "tweet2"
		liked, err := repo.HasLiked(ctx, userID, tweetID)
		assert.NoError(t, err)
		assert.False(t, liked)
	})

	t.Run("Like is idempotent", func(t *testing.T) {
		userID, tweetID := "user3", "tweet3"
		assert.NoError(t, repo.Like(ctx, userID, tweetID))
		assert.NoError(t, repo.Like(ctx, userID, tweetID))

		liked, err := repo.HasLiked(ctx, userID, tweetID)
		assert.NoError(t, err)
		assert.True(t, liked)
	})

	t.Run("Different users can like same tweet", func(t *testing.T) {
		tweetID := "tweet4"
		assert.NoError(t, repo.Like(ctx, "userA", tweetID))
		assert.NoError(t, repo.Like(ctx, "userB", tweetID))

		likedA, _ := repo.HasLiked(ctx, "userA", tweetID)
		likedB, _ := repo.HasLiked(ctx, "userB", tweetID)
		assert.True(t, likedA)
		assert.True(t, likedB)
	})

	t.Run("Same user likes multiple tweets", func(t *testing.T) {
		userID := "userX"
		assert.NoError(t, repo.Like(ctx, userID, "tweetA"))
		assert.NoError(t, repo.Like(ctx, userID, "tweetB"))

		likedA, _ := repo.HasLiked(ctx, userID, "tweetA")
		likedB, _ := repo.HasLiked(ctx, userID, "tweetB")
		assert.True(t, likedA)
		assert.True(t, likedB)
	})

	t.Run("HasLiked returns false for unknown user", func(t *testing.T) {
		liked, err := repo.HasLiked(ctx, "ghost", "nope")
		assert.NoError(t, err)
		assert.False(t, liked)
	})
}
