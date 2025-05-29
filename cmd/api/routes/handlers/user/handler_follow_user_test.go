package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"ualaTwitter/internal/usecase/follow_user"
)

type fakeFollowUserService struct {
	Err error
}

func (f *fakeFollowUserService) Execute(_ context.Context, _ follow_user.Input) error {
	return f.Err
}

func TestFollowUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		header         string
		body           map[string]string
		mockService    *fakeFollowUserService
		expectedStatus int
	}{
		{
			name:   "successfully follows a user",
			header: "usr_123",
			body: map[string]string{
				"followee_id": "usr_456",
			},
			mockService:    &fakeFollowUserService{},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing X-User-ID header",
			header:         "",
			body:           map[string]string{"followee_id": "usr_456"},
			mockService:    &fakeFollowUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty body",
			header:         "usr_123",
			body:           nil,
			mockService:    &fakeFollowUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "missing followee_id",
			header: "usr_123",
			body:   map[string]string{"followee_id": ""},
			mockService: &fakeFollowUserService{
				Err: nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "use case returns error",
			header: "usr_123",
			body:   map[string]string{"followee_id": "usr_456"},
			mockService: &fakeFollowUserService{
				Err: errors.New("already following user"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var bodyBytes []byte
			if tc.body != nil {
				bodyBytes, _ = json.Marshal(tc.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/follow", bytes.NewReader(bodyBytes))
			if tc.header != "" {
				req.Header.Set("X-User-ID", tc.header)
			}

			rr := httptest.NewRecorder()

			handler := NewFollowUserHandler(tc.mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
