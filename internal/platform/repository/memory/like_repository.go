package memory

import (
	"context"
	"sync"
)

type InMemoryLikeRepository struct {
	mu    sync.RWMutex
	likes map[string]map[string]struct{}
}

func NewInMemoryLikeRepository() *InMemoryLikeRepository {
	return &InMemoryLikeRepository{
		likes: make(map[string]map[string]struct{}),
	}
}

func (r *InMemoryLikeRepository) Like(ctx context.Context, userID, tweetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.likes[userID] == nil {
		r.likes[userID] = make(map[string]struct{})
	}

	r.likes[userID][tweetID] = struct{}{}
	return nil
}

func (r *InMemoryLikeRepository) HasLiked(ctx context.Context, userID, tweetID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.likes[userID][tweetID]
	return ok, nil
}
