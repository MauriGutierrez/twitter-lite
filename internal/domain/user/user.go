package user

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	userPrefix    = "usr_"
	MaxNameLength = 100
)

type User struct {
	ID       string
	Name     string
	Document string
}

func New(name, document string) (User, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" || utf8.RuneCountInString(trimmedName) > MaxNameLength {
		return User{}, ErrInvalidName
	}

	if len(document) < 7 || len(document) > 8 {
		return User{}, ErrInvalidDocument
	}

	for _, r := range document {
		if r < '0' || r > '9' {
			return User{}, ErrInvalidDocument
		}
	}

	return User{
		ID:       generateUserID(document),
		Name:     name,
		Document: document,
	}, nil
}

func generateUserID(document string) string {
	return fmt.Sprintf("%s%s", userPrefix, document)
}
