package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrAlreadyFollowing  = errors.New("user already followed")
	ErrInvalidInput      = errors.New("user_id or followee_id is empty")
	ErrInvalidName       = errors.New("name is empty or too long")
	ErrInvalidDocument   = errors.New("invalid document")
	ErrSelfFollow        = errors.New("cannot follow yourself")
	ErrUserAlreadyExists = errors.New("user already exists")
)
