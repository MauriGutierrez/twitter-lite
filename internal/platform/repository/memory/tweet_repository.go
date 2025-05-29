package memory

import (
	"context"
	"fmt"
	"sync"
	"ualaTwitter/internal/domain/user"

	"ualaTwitter/internal/domain/tweet"
)

type InMemoryTweetRepository struct {
	mu     sync.RWMutex
	byID   map[string]tweet.Tweet
	byUser map[string][]tweet.Tweet
}

func NewInMemoryTweetRepository() *InMemoryTweetRepository {
	return &InMemoryTweetRepository{
		byID:   make(map[string]tweet.Tweet),
		byUser: make(map[string][]tweet.Tweet),
	}
}

func (r *InMemoryTweetRepository) Save(ctx context.Context, t tweet.Tweet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[t.ID] = t
	r.byUser[t.UserID] = append(r.byUser[t.UserID], t)

	return nil
}

func (r *InMemoryTweetRepository) GetByID(ctx context.Context, id string) (tweet.Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, ok := r.byID[id]
	if !ok {
		return tweet.Tweet{}, fmt.Errorf("tweet not found")
	}
	return t, nil
}

func (r *InMemoryTweetRepository) FindTweetsAuthoredBy(ctx context.Context, userID string) ([]tweet.Tweet, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tweets, ok := r.byUser[userID]
	if !ok {
		return []tweet.Tweet{}, user.ErrUserNotFound
	}
	return tweets, nil
}

func (r *InMemoryTweetRepository) IncrementLikes(ctx context.Context, tweetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	t, ok := r.byID[tweetID]
	if !ok {
		return fmt.Errorf("tweet not found")
	}
	t.Likes++
	r.byID[tweetID] = t

	userTweets := r.byUser[t.UserID]
	for i, tw := range userTweets {
		if tw.ID == tweetID {
			userTweets[i] = t
			break
		}
	}

	return nil
}
