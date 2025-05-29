package mocks

import (
	"context"
	"ualaTwitter/internal/domain/user"
)

type FakeUserRepo struct {
	Users     map[string]*user.User
	Followees map[string][]string
	FollowErr error
	CreateErr error
}

func (f *FakeUserRepo) GetByID(_ context.Context, id string) (user.User, error) {
	u, ok := f.Users[id]
	if !ok {
		return user.User{}, user.ErrUserNotFound
	}
	return *u, nil
}

func (f *FakeUserRepo) Create(_ context.Context, u user.User) error {
	if f.CreateErr != nil {
		return f.CreateErr
	}
	f.Users[u.ID] = &u
	return nil
}

func (f *FakeUserRepo) Follow(_ context.Context, followerID, followeeID string) error {
	if f.FollowErr != nil {
		return f.FollowErr
	}
	f.Followees[followerID] = append(f.Followees[followerID], followeeID)
	return nil
}

func (f *FakeUserRepo) GetUsersFollowedBy(_ context.Context, userID string) ([]string, error) {
	followees, ok := f.Followees[userID]
	if !ok {
		return nil, nil
	}
	return followees, nil
}
