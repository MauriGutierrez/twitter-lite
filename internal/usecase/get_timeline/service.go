package get_timeline

import (
	"context"
	"errors"
	"sort"

	"golang.org/x/sync/errgroup"
	"ualaTwitter/internal/domain/tweet"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/platform/errors/usecase"
)

const (
	defaultLimit = 10
	maxLimit     = 100
	maxOffset    = 1000
)

type GetTimelineService struct {
	TweetRepo tweet.Repository
	UserRepo  user.InMemoryRepository
}

func NewGetTimelineService(tweetRepo tweet.Repository, userRepo user.InMemoryRepository) *GetTimelineService {
	return &GetTimelineService{
		TweetRepo: tweetRepo,
		UserRepo:  userRepo,
	}
}

func (s *GetTimelineService) Execute(ctx context.Context, input Input) ([]TweetTimeline, error) {
	if _, err := s.UserRepo.GetByID(ctx, input.UserID); err != nil {
		return nil, usecase.NotFound(user.ErrUserNotFound.Error())
	}

	followees, err := s.getFollowees(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	tweets, err := s.aggregateTweets(ctx, followees)
	if err != nil {
		return nil, err
	}

	sorted := s.sortTweets(tweets)
	paginated := s.paginateTweets(sorted, input.Offset, input.Limit)
	return s.mapToTimelineResponse(paginated), nil
}

func (s *GetTimelineService) getFollowees(ctx context.Context, userID string) ([]string, error) {
	followees, err := s.UserRepo.GetUsersFollowedBy(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			return nil, usecase.NotFound("user not found", err)
		default:
			return nil, usecase.InternalServerError("failed to fetch followees", err)
		}
	}
	return followees, nil
}

func (s *GetTimelineService) aggregateTweets(ctx context.Context, followees []string) ([]tweet.Tweet, error) {
	var (
		timelineTweetsMutex = make(chan []tweet.Tweet, len(followees))
	)

	g, ctx := errgroup.WithContext(ctx)

	for _, fid := range followees {
		fid := fid
		g.Go(func() error {
			tweets, err := s.TweetRepo.FindTweetsAuthoredBy(ctx, fid)
			if err != nil {
				if errors.Is(err, user.ErrUserNotFound) {
					return nil
				}
				return err
			}
			timelineTweetsMutex <- tweets
			return nil
		})
	}

	err := g.Wait()
	close(timelineTweetsMutex)

	if err != nil {
		return nil, usecase.InternalServerError("failed to fetch tweets for followed users", err)
	}

	var timeline []tweet.Tweet
	for group := range timelineTweetsMutex {
		timeline = append(timeline, group...)
	}

	return timeline, nil
}

func (s *GetTimelineService) sortTweets(tweets []tweet.Tweet) []tweet.Tweet {
	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.After(tweets[j].CreatedAt)
	})
	return tweets
}

func (s *GetTimelineService) paginateTweets(tweets []tweet.Tweet, offset, limit int) []tweet.Tweet {
	offset, limit = normalizePaginationParams(offset, limit)

	if offset >= len(tweets) {
		return []tweet.Tweet{}
	}

	end := offset + limit
	if end > len(tweets) {
		end = len(tweets)
	}

	return tweets[offset:end]
}

func (s *GetTimelineService) mapToTimelineResponse(tweets []tweet.Tweet) []TweetTimeline {
	result := make([]TweetTimeline, len(tweets))
	for i, t := range tweets {
		result[i] = TweetTimeline{
			ID:        t.ID,
			UserID:    t.UserID,
			Content:   t.Content,
			Likes:     t.Likes,
			CreatedAt: t.CreatedAt,
		}
	}
	return result
}

func normalizePaginationParams(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if offset > maxOffset {
		offset = maxOffset
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return offset, limit
}
