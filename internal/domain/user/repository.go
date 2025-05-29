package user

import "context"

type InMemoryRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Follow(ctx context.Context, followerID, followeeID string) error
	GetUsersFollowedBy(ctx context.Context, userID string) ([]string, error)
}

type PostgresRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
}
