package like

import "errors"

var (
	ErrAlreadyLiked = errors.New("user has already liked this tweet")
	ErrInvalidInput = errors.New("invalid userID or tweetID")
)
