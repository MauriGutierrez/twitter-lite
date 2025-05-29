package mocks

import (
	"context"
)

type FakeLikeRepo struct {
	AlreadyLiked bool
	HasLikedErr  error
	LikeErr      error
	LikedTweetID string
	LikedUserID  string
}

func (f *FakeLikeRepo) HasLiked(_ context.Context, tweetID, userID string) (bool, error) {
	return f.AlreadyLiked, f.HasLikedErr
}

func (f *FakeLikeRepo) Like(_ context.Context, tweetID, userID string) error {
	f.LikedTweetID = tweetID
	f.LikedUserID = userID
	return f.LikeErr
}
