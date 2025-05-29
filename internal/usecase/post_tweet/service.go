package post_tweet

import (
	"context"
	"errors"
	"time"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/platform/errors/usecase"

	"ualaTwitter/internal/domain/tweet"
)

type PostTweetService struct {
	TweetRepo tweet.Repository
	UserRepo  user.InMemoryRepository
}

func NewPostTweetService(tweetRepo tweet.Repository, memoryUserRepository user.InMemoryRepository) *PostTweetService {
	return &PostTweetService{
		TweetRepo: tweetRepo,
		UserRepo:  memoryUserRepository,
	}

}

func (s *PostTweetService) Execute(ctx context.Context, input Input) (string, error) {
	if _, err := s.UserRepo.GetByID(ctx, input.UserID); err != nil {
		return "", usecase.NotFound("", user.ErrUserNotFound)
	}

	newTweet, err := tweet.New(input.UserID, input.Content, time.Now())
	if err != nil {
		switch {
		case errors.Is(err, tweet.ErrEmptyTweet), errors.Is(err, tweet.ErrTooLong):
			return "", usecase.InvalidParam("invalid tweet content", err)
		default:
			return "", usecase.InternalServerError("failed to create tweet", err)
		}
	}

	if err := s.TweetRepo.Save(ctx, newTweet); err != nil {
		return "", usecase.InternalServerError("failed to persist tweet", err)
	}

	return newTweet.ID, nil
}
