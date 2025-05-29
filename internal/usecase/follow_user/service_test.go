package follow_user

import (
	"context"
	"errors"
	"testing"
	"ualaTwitter/internal/domain/user"
	"ualaTwitter/internal/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestFollowUserService_Execute(t *testing.T) {
	ctx := context.Background()

	follower := &user.User{ID: "follower", Name: "F", Document: "100"}
	followee := &user.User{ID: "followee", Name: "EE", Document: "200"}

	tests := []struct {
		name      string
		input     Input
		users     map[string]*user.User
		followees map[string][]string
		followErr error
		expectErr string
	}{
		{
			name:      "successfully follows user",
			input:     Input{FollowerID: follower.ID, FolloweeID: followee.ID},
			users:     map[string]*user.User{follower.ID: follower, followee.ID: followee},
			followErr: nil,
		},
		{
			name:      "empty follower or followee",
			input:     Input{FollowerID: "", FolloweeID: ""},
			users:     map[string]*user.User{},
			expectErr: "cannot be empty",
		},
		{
			name:      "self-follow is forbidden",
			input:     Input{FollowerID: follower.ID, FolloweeID: follower.ID},
			users:     map[string]*user.User{follower.ID: follower},
			expectErr: "cannot follow yourself",
		},
		{
			name:      "follower not found",
			input:     Input{FollowerID: "ghost", FolloweeID: followee.ID},
			users:     map[string]*user.User{followee.ID: followee},
			expectErr: "follower not found",
		},
		{
			name:      "followee not found",
			input:     Input{FollowerID: follower.ID, FolloweeID: "ghost"},
			users:     map[string]*user.User{follower.ID: follower},
			expectErr: "followee not found",
		},
		{
			name:      "already following",
			input:     Input{FollowerID: follower.ID, FolloweeID: followee.ID},
			users:     map[string]*user.User{follower.ID: follower, followee.ID: followee},
			followErr: user.ErrAlreadyFollowing,
			expectErr: "already following",
		},
		{
			name:      "follow error returns internal server error",
			input:     Input{FollowerID: follower.ID, FolloweeID: followee.ID},
			users:     map[string]*user.User{follower.ID: follower, followee.ID: followee},
			followErr: errors.New("db failure"),
			expectErr: "could not follow user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repo := &mocks.FakeUserRepo{
				Users:     tc.users,
				Followees: make(map[string][]string),
				FollowErr: tc.followErr,
			}

			service := NewFollowUserService(repo)
			err := service.Execute(ctx, tc.input)

			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
