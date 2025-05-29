package mocks

import (
	"context"
	"ualaTwitter/internal/domain/tweet"
)

type FakeTweetRepo struct {
	SaveErr           error
	Saved             *tweet.Tweet
	LastLikedTweetID  string
	IncrementLikesErr error
	TweetsByUser      map[string][]tweet.Tweet
	TweetFetchErr     map[string]error
}

func (f *FakeTweetRepo) Save(_ context.Context, t tweet.Tweet) error {
	if f.SaveErr != nil {
		return f.SaveErr
	}
	f.Saved = &t
	return nil
}

func (f *FakeTweetRepo) FindTweetsAuthoredBy(_ context.Context, userID string) ([]tweet.Tweet, error) {
	if f.TweetFetchErr != nil {
		if err, exists := f.TweetFetchErr[userID]; exists {
			return nil, err
		}
	}
	if f.TweetsByUser != nil {
		if tweets, exists := f.TweetsByUser[userID]; exists {
			return tweets, nil
		}
	}
	return nil, nil
}

func (f *FakeTweetRepo) GetByID(_ context.Context, _ string) (tweet.Tweet, error) {
	return tweet.Tweet{}, nil
}

func (f *FakeTweetRepo) IncrementLikes(_ context.Context, tweetID string) error {
	f.LastLikedTweetID = tweetID
	return f.IncrementLikesErr
}
