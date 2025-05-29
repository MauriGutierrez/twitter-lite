package postgres

import (
	"context"
	"testing"
	"ualaTwitter/cmd/api/config"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"ualaTwitter/internal/domain/user"
)

func setupTestDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), config.Load().PostgresDSN)
	if err != nil {
		panic("Failed to connect to test DB: " + err.Error())
	}
	_, _ = conn.Exec(context.Background(), "DELETE FROM users")
	return conn
}

func TestPostgresUserRepository(t *testing.T) {
	ctx := context.Background()
	conn := setupTestDB()
	repo := NewPostgresUserRepository(conn)

	t.Cleanup(func() {
		_, _ = conn.Exec(ctx, "DELETE FROM users")
	})

	t.Run("Create and GetByID happy path", func(t *testing.T) {
		u := user.User{ID: "usr_test1", Name: "Test User", Document: "12345678"}
		assert.NoError(t, repo.Create(ctx, u))

		got, err := repo.GetByID(ctx, u.ID)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, got.ID)
		assert.Equal(t, u.Name, got.Name)
	})

	t.Run("GetByID returns error for missing user", func(t *testing.T) {
		_, err := repo.GetByID(ctx, "no_such_id")
		assert.Error(t, err)
	})

	t.Run("Create returns error for duplicate ID", func(t *testing.T) {
		u := user.User{ID: "usr_test2", Name: "Test Duplicate", Document: "222"}
		assert.NoError(t, repo.Create(ctx, u))
		err := repo.Create(ctx, u)
		assert.Error(t, err)
	})
}
