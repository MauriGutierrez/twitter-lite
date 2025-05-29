package user_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"ualaTwitter/internal/domain/user"
	"unicode/utf8"
)

func TestUser_New(t *testing.T) {
	tests := []struct {
		name      string
		inputName string
		document  string
		wantError error
	}{
		{
			name:      "valid user with 7 digit document",
			inputName: "Alice",
			document:  "1234567",
			wantError: nil,
		},
		{
			name:      "valid user with 8 digit document",
			inputName: "Bob",
			document:  "87654321",
			wantError: nil,
		},
		{
			name:      "empty name",
			inputName: "",
			document:  "1234567",
			wantError: user.ErrInvalidName,
		},
		{
			name:      "empty document",
			inputName: "Charlie",
			document:  "",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "name exceeds max length (runes)",
			inputName: string(make([]rune, user.MaxNameLength+1)),
			document:  "1234567",
			wantError: user.ErrInvalidName,
		},
		{
			name:      "document too short",
			inputName: "Eve",
			document:  "123456",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "document too long",
			inputName: "Frank",
			document:  "123456789",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "document contains non-digit character",
			inputName: "Grace",
			document:  "1234a567",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "document contains special character",
			inputName: "Heidi",
			document:  "12345-67",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "document contains whitespace",
			inputName: "Ivan",
			document:  "1234 567",
			wantError: user.ErrInvalidDocument,
		},
		{
			name:      "name with multi-byte characters (valid)",
			inputName: "Ñandú ☀️",
			document:  "87654321",
			wantError: nil,
		},
		{
			name:      "name with only spaces (should fail)",
			inputName: "    ",
			document:  "1234567",
			wantError: user.ErrInvalidName,
		},
		{
			name:      "document all zeros (should be valid)",
			inputName: "Oscar",
			document:  "00000000",
			wantError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u, err := user.New(tc.inputName, tc.document)
			if tc.wantError != nil {
				assert.ErrorIs(t, err, tc.wantError)
				assert.Empty(t, u.ID)
				assert.Empty(t, u.Name)
				assert.Empty(t, u.Document)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.inputName, u.Name)
				assert.Equal(t, tc.document, u.Document)
				assert.Contains(t, u.ID, u.Document)
				assert.True(t, utf8.RuneCountInString(u.Name) <= user.MaxNameLength)
			}
		})
	}
}
