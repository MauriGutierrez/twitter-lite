package tweet_test

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"ualaTwitter/internal/domain/tweet"
)

func TestTweet_New(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		userID    string
		content   string
		wantError error
	}{
		{
			name:      "valid tweet",
			userID:    "user-123",
			content:   "hello, world!",
			wantError: nil,
		},
		{
			name:      "empty content",
			userID:    "user-123",
			content:   "   ",
			wantError: tweet.ErrEmptyTweet,
		}, {
			name:   "content too long",
			userID: "user-123",
			content: "A very very very very very very very very very very " +
				"very very very very very very very very very very very very" +
				"very very very very very very very very very very very very " +
				"very very very very very very very very very very very very " +
				"very very very very very very very very very very very very " +
				"very very very very very very very very very very very very " +
				"very very very very very very very very very very very very extensive tweet",
			wantError: tweet.ErrTooLong,
		},
		{
			name:      "content with spaces is trimmed",
			userID:    "user-123",
			content:   "   trimmed tweet    ",
			wantError: nil,
		},
		{
			name:      "tweet with emojis",
			userID:    "user-123",
			content:   "testing emojis 👍👍👍!",
			wantError: nil,
		},
		{
			name:      "content only whitespace/newlines",
			userID:    "user-123",
			content:   "\n\t  \n\t",
			wantError: tweet.ErrEmptyTweet,
		},
		{
			name:      "empty user id",
			userID:    "",
			content:   "valid tweet",
			wantError: tweet.ErrInvalidUser,
		},
		{
			name:      "tweet at max content length",
			userID:    "user-123",
			content:   strings.Repeat("a", tweet.MaxContentLength),
			wantError: nil,
		},
		{
			name:      "tweet at max length with emojis",
			userID:    "user-123",
			content:   strings.Repeat("👍", tweet.MaxContentLength),
			wantError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tw, err := tweet.New(tc.userID, tc.content, now)

			if tc.wantError != nil {
				assert.ErrorIs(t, err, tc.wantError)
				assert.Empty(t, tw.ID)
				assert.Empty(t, tw.Content)
				assert.Zero(t, tw.Likes)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tw.ID)
				assert.Equal(t, tc.userID, tw.UserID)
				assert.Equal(t, now, tw.CreatedAt)
				assert.Greater(t, utf8.RuneCountInString(tw.Content), 0)
				assert.LessOrEqual(t, utf8.RuneCountInString(tw.Content), tweet.MaxContentLength)
				if strings.HasPrefix(tc.name, "tweet at max") {
					assert.Equal(t, tweet.MaxContentLength, utf8.RuneCountInString(tw.Content))
				}
				assert.Zero(t, tw.Likes)
			}
		})
	}
}
