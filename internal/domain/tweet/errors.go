package tweet

import "errors"

var (
	ErrTooLong     = errors.New("too long, maximum length (280 characters) exceeded")
	ErrEmptyTweet  = errors.New("empty")
	ErrNotFound    = errors.New("not found")
	ErrInvalidUser = errors.New("tweet user is not valid")
)
