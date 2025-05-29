package memory

import (
	"context"
	"sync"
	"ualaTwitter/internal/domain/user"
)

type InMemoryUserRepository struct {
	mu      sync.RWMutex
	follows map[string]map[string]struct{}
	users   map[string]user.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		follows: make(map[string]map[string]struct{}),
		users:   make(map[string]user.User),
	}
}

func (r *InMemoryUserRepository) Create(ctx context.Context, u user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[u.ID]; exists {
		return user.ErrUserAlreadyExists
	}
	r.users[u.ID] = u
	return nil
}

func (r *InMemoryUserRepository) GetByID(ctx context.Context, id string) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.users[id]
	if !exists {
		return user.User{}, user.ErrUserNotFound
	}
	return u, nil
}

func (r *InMemoryUserRepository) Follow(ctx context.Context, followerID, followeeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.follows[followerID] == nil {
		r.follows[followerID] = make(map[string]struct{})
	}

	r.follows[followerID][followeeID] = struct{}{}
	return nil
}

func (r *InMemoryUserRepository) GetUsersFollowedBy(ctx context.Context, userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	followees := make([]string, 0)

	for fid := range r.follows[userID] {
		followees = append(followees, fid)
	}

	return followees, nil
}
