package get_timeline

import (
	"context"
	"errors"
	"testing"
	"time"

	"ualaTwitter/internal/domain/tweet"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetTimelineService_Execute(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tweet1 := tweet.Tweet{ID: "t1", UserID: "usr_1", Content: "A", CreatedAt: now.Add(-2 * time.Minute)}
	tweet2 := tweet.Tweet{ID: "t2", UserID: "usr_1", Content: "B", CreatedAt: now.Add(-1 * time.Minute)}
	tweet3 := tweet.Tweet{ID: "t3", UserID: "usr_2", Content: "C", CreatedAt: now.Add(-3 * time.Minute)}
	tweet4 := tweet.Tweet{ID: "t4", UserID: "usr_3", Content: "", CreatedAt: now.Add(-3 * time.Minute)}

	tests := []struct {
		name          string
		userExists    bool
		followees     []string
		tweetsByUser  map[string][]tweet.Tweet
		tweetErrors   map[string]error
		offset        int
		limit         int
		expectedCount int
		expectError   string
	}{
		{
			name:       "happy path, two followees, sorted & paginated",
			userExists: true,
			followees:  []string{"usr_1", "usr_2"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1, tweet2},
				"usr_2": {tweet3},
			},
			offset:        0,
			limit:         10,
			expectedCount: 3,
		},
		{
			name:        "user not found",
			userExists:  false,
			expectError: "not found",
		},
		{
			name:          "no followees",
			userExists:    true,
			followees:     []string{},
			expectedCount: 0,
		},
		{
			name:       "tweet fetch error for one followee",
			userExists: true,
			followees:  []string{"usr_1", "usr_2"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1},
			},
			tweetErrors: map[string]error{
				"usr_2": errors.New("db failure"),
			},
			expectError: "failed to fetch tweets",
		},
		{
			name:       "tweet fetch not found, skip that followee",
			userExists: true,
			followees:  []string{"usr_1", "usr_2"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1},
			},
			tweetErrors: map[string]error{
				"usr_2": user.ErrUserNotFound,
			},
			expectedCount: 1,
		},
		{
			name:       "pagination offset exceeds tweets",
			userExists: true,
			followees:  []string{"usr_1"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1, tweet2},
			},
			offset:        10,
			limit:         10,
			expectedCount: 0,
		},
		{
			name:       "pagination trims to limit",
			userExists: true,
			followees:  []string{"usr_1"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1, tweet2},
			},
			offset:        0,
			limit:         1,
			expectedCount: 1,
		},
		{
			name:       "negative offset and limit normalized",
			userExists: true,
			followees:  []string{"usr_1"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1, tweet2},
			},
			offset:        -5,
			limit:         -10,
			expectedCount: 2,
		},
		{
			name:       "timeline contains empty-content tweet",
			userExists: true,
			followees:  []string{"usr_1", "usr_3"},
			tweetsByUser: map[string][]tweet.Tweet{
				"usr_1": {tweet1},
				"usr_3": {tweet4},
			},
			offset:        0,
			limit:         10,
			expectedCount: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := &mocks.FakeUserRepo{}
			if tc.userExists {
				userRepo.Users = map[string]*user.User{"test_user": {ID: "test_user", Name: "User", Document: "123"}}
			} else {
				userRepo.Users = map[string]*user.User{}
			}
			userRepo.Followees = map[string][]string{"test_user": tc.followees}

			tweetRepo := &mocks.FakeTweetRepo{}
			tweetRepo.TweetsByUser = tc.tweetsByUser
			tweetRepo.TweetFetchErr = tc.tweetErrors

			service := NewGetTimelineService(tweetRepo, userRepo)

			input := Input{
				UserID: "test_user",
				Limit:  tc.limit,
				Offset: tc.offset,
			}
			result, err := service.Execute(ctx, input)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCount, len(result))
				assert.Equal(t, tc.expectedCount, len(result))
				if len(result) > 1 {
					for i := 1; i < len(result); i++ {
						assert.True(t,
							result[i-1].CreatedAt.After(result[i].CreatedAt) || result[i-1].CreatedAt.Equal(result[i].CreatedAt),
							"timeline is not sorted: tweet %d (at %v) should be >= tweet %d (at %v)",
							i-1, result[i-1].CreatedAt, i, result[i].CreatedAt,
						)
					}
				}
			}
		})
	}
}
