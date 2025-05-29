package create_user

import (
	"context"
	"go.uber.org/zap"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/platform/errors/usecase"
	"ualaTwitter/internal/platform/logger"
)

type CreateUserService struct {
	pgRepo     user.PostgresRepository
	memoryRepo user.InMemoryRepository
}

func NewCreateUserService(pgRepo user.PostgresRepository, memoryRepo user.InMemoryRepository) *CreateUserService {
	return &CreateUserService{
		pgRepo:     pgRepo,
		memoryRepo: memoryRepo}
}

func (s *CreateUserService) Execute(ctx context.Context, input Input) (Output, error) {
	newUser, err := user.New(input.Name, input.Document)
	if err != nil {
		return Output{}, usecase.InvalidParam("invalid user name", err)
	}

	if _, err := s.memoryRepo.GetByID(ctx, newUser.ID); err == nil {
		logger.Log.Info("user already exists in memory", zap.String("user_id", newUser.ID))
		return Output{}, usecase.Conflict("user already exists", user.ErrUserAlreadyExists)
	}

	existingUser, err := s.pgRepo.GetByID(ctx, newUser.ID)
	if err == nil {
		logger.Log.Info("user already exists in postgres", zap.String("user_id", existingUser.ID))

		if err := s.saveUserInMemory(ctx, existingUser); err != nil {
			logger.Log.Warn("failed to cache existing user in memory", zap.Error(err))
		}
		return Output{}, usecase.Conflict("user already exists", user.ErrUserAlreadyExists)
	}

	if err := s.pgRepo.Create(ctx, newUser); err != nil {
		return Output{}, usecase.InternalServerError("failed to persist user", err)
	}

	if err := s.saveUserInMemory(ctx, newUser); err != nil {
		logger.Log.Warn("failed to store user in memory",
			zap.String("user_id", newUser.ID),
			zap.Error(err),
		)
	}

	return Output{ID: newUser.ID}, nil
}

func (s *CreateUserService) saveUserInMemory(ctx context.Context, user user.User) error {
	return s.memoryRepo.Create(ctx, user)
}
