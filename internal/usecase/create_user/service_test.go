package create_user

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"testing"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/platform/logger"
	"ualaTwitter/internal/test/mocks"

	"github.com/stretchr/testify/assert"
)

func init() {
	logger.Log = zap.NewNop()
}

func TestCreateUserService_Execute(t *testing.T) {
	ctx := context.Background()
	validName := "Alice"
	validDoc := "12345678"
	userID := "usr_" + validDoc
	validUser := user.User{ID: userID, Name: validName, Document: validDoc}

	tests := []struct {
		name              string
		input             Input
		memRepoUser       *user.User
		memRepoErr        error
		pgRepoUser        *user.User
		pgRepoCreateErr   error
		expectErrContains string
		expectUserID      string
	}{
		{
			name:         "successfully creates user",
			input:        Input{Name: validName, Document: validDoc},
			memRepoUser:  nil,
			pgRepoUser:   nil,
			expectUserID: userID,
		},
		{
			name:              "invalid user name",
			input:             Input{Name: "", Document: validDoc},
			expectErrContains: "invalid user name",
		},
		{
			name:              "user already exists in memory",
			input:             Input{Name: validName, Document: validDoc},
			memRepoUser:       &validUser,
			expectErrContains: "already exists",
		},
		{
			name:              "user exists in Postgres, not in memory (rehydrates)",
			input:             Input{Name: validName, Document: validDoc},
			memRepoUser:       nil,
			pgRepoUser:        &validUser,
			expectErrContains: "already exists",
		},
		{
			name:              "Postgres create error",
			input:             Input{Name: validName, Document: validDoc},
			memRepoUser:       nil,
			pgRepoUser:        nil,
			pgRepoCreateErr:   errors.New("db write failed"),
			expectErrContains: "failed to persist user",
		},
		{
			name:         "memory write after Postgres create fails",
			input:        Input{Name: validName, Document: validDoc},
			memRepoUser:  nil,
			pgRepoUser:   nil,
			memRepoErr:   errors.New("memory error"),
			expectUserID: userID,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fakeMemRepo := &mocks.FakeUserRepo{
				Users:     make(map[string]*user.User),
				CreateErr: tc.memRepoErr,
			}
			if tc.memRepoUser != nil {
				fakeMemRepo.Users[tc.memRepoUser.ID] = tc.memRepoUser
			}

			fakePgRepo := &mocks.FakeUserRepo{
				Users:     make(map[string]*user.User),
				CreateErr: tc.pgRepoCreateErr,
			}
			if tc.pgRepoUser != nil {
				fakePgRepo.Users[tc.pgRepoUser.ID] = tc.pgRepoUser
			}

			service := NewCreateUserService(fakePgRepo, fakeMemRepo)
			out, err := service.Execute(ctx, tc.input)

			if tc.expectErrContains != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectUserID, out.ID)
			}
		})
	}
}
