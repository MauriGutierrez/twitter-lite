package tweet

import (
	"github.com/google/uuid"
	"strings"
	"time"
	"unicode/utf8"
)

const MaxContentLength = 280

type Tweet struct {
	ID        string
	UserID    string
	Content   string
	CreatedAt time.Time
	Likes     int
}

func New(userID, content string, createdAt time.Time) (Tweet, error) {
	if userID == "" {
		return Tweet{}, ErrInvalidUser
	}

	trimmedTweet := strings.TrimSpace(content)

	if trimmedTweet == "" {
		return Tweet{}, ErrEmptyTweet
	}

	if utf8.RuneCountInString(trimmedTweet) > MaxContentLength {
		return Tweet{}, ErrTooLong
	}

	return Tweet{
		ID:        uuid.NewString(),
		UserID:    userID,
		Content:   trimmedTweet,
		CreatedAt: createdAt,
		Likes:     0,
	}, nil
}
