package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"ualaTwitter/internal/domain/user"
)

func makeMockUser(id, name, document string) user.User {
	return user.User{
		ID:       id,
		Name:     name,
		Document: document,
	}
}

func TestInMemoryUserRepository(t *testing.T) {
	ctx := context.Background()
	repo := NewInMemoryUserRepository()

	t.Run("Create and GetByID stores and retrieves user", func(t *testing.T) {
		u := makeMockUser("usr1", "Test User", "12345678")
		assert.NoError(t, repo.Create(ctx, u))

		got, err := repo.GetByID(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.Name, got.Name)
		assert.Equal(t, u.Document, got.Document)
	})

	t.Run("Create returns error for duplicate user", func(t *testing.T) {
		u := makeMockUser("usr2", "User2", "87654321")
		assert.NoError(t, repo.Create(ctx, u))
		err := repo.Create(ctx, u)
		assert.ErrorIs(t, err, user.ErrUserAlreadyExists)
	})

	t.Run("GetByID returns error for unknown user", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "ghost")
		assert.ErrorIs(t, err, user.ErrUserNotFound)
	})

	t.Run("Follow adds followee, GetUsersFollowedBy returns them", func(t *testing.T) {
		repo = NewInMemoryUserRepository()
		follower := makeMockUser("follower", "Follower", "111")
		followee1 := makeMockUser("followee1", "Followee 1", "222")
		followee2 := makeMockUser("followee2", "Followee 2", "333")

		repo.Create(ctx, follower)
		repo.Create(ctx, followee1)
		repo.Create(ctx, followee2)

		assert.NoError(t, repo.Follow(ctx, follower.ID, followee1.ID))
		assert.NoError(t, repo.Follow(ctx, follower.ID, followee2.ID))

		followees, err := repo.GetUsersFollowedBy(ctx, follower.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{followee1.ID, followee2.ID}, followees)
	})

	t.Run("Follow is idempotent (no error on repeat)", func(t *testing.T) {
		repo = NewInMemoryUserRepository()
		follower := makeMockUser("follower2", "Follower 2", "444")
		followee := makeMockUser("followee3", "Followee 3", "555")

		repo.Create(ctx, follower)
		repo.Create(ctx, followee)

		assert.NoError(t, repo.Follow(ctx, follower.ID, followee.ID))
		assert.NoError(t, repo.Follow(ctx, follower.ID, followee.ID))

		followees, _ := repo.GetUsersFollowedBy(ctx, follower.ID)
		assert.Equal(t, []string{followee.ID}, followees)
	})

	t.Run("GetUsersFollowedBy returns empty slice for user with no follows", func(t *testing.T) {
		repo = NewInMemoryUserRepository()
		u := makeMockUser("solo", "Solo User", "666")
		repo.Create(ctx, u)
		followees, err := repo.GetUsersFollowedBy(ctx, u.ID)
		assert.NoError(t, err)
		assert.Len(t, followees, 0)
	})
}
