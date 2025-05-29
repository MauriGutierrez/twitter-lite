package follow_user

import (
	"context"
	"errors"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/platform/errors/usecase"
)

type FollowUserService struct {
	UserRepo user.InMemoryRepository
}

func NewFollowUserService(userRepo user.InMemoryRepository) *FollowUserService {
	return &FollowUserService{
		UserRepo: userRepo,
	}
}

func (s *FollowUserService) Execute(ctx context.Context, input Input) error {
	if input.FollowerID == "" || input.FolloweeID == "" {
		return usecase.InvalidParam("user ID and followee ID cannot be empty", user.ErrInvalidInput)
	}

	if input.FollowerID == input.FolloweeID {
		return usecase.Forbidden("cannot follow yourself", user.ErrSelfFollow)
	}

	if _, err := s.UserRepo.GetByID(ctx, input.FollowerID); err != nil {
		return usecase.NotFound("follower not found", user.ErrUserNotFound)
	}

	_, err := s.UserRepo.GetByID(ctx, input.FolloweeID)
	if err != nil {
		return usecase.NotFound("followee not found", user.ErrUserNotFound)
	}

	err = s.UserRepo.Follow(ctx, input.FollowerID, input.FolloweeID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAlreadyFollowing):
			return usecase.Forbidden("already following this user", err)
		case errors.Is(err, user.ErrUserNotFound):
			return usecase.NotFound("user not found", err)
		default:
			return usecase.InternalServerError("could not follow user", err)
		}
	}

	return nil
}
