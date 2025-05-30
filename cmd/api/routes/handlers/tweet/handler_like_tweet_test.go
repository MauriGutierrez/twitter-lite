package tweet

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"ualaTwitter/internal/usecase/like_tweet"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type fakeLikeTweetService struct {
	Err error
}

func (f *fakeLikeTweetService) Execute(_ context.Context, _ like_tweet.Input) error {
	return f.Err
}

func TestLikeTweetHandler(t *testing.T) {
	tests := []struct {
		name           string
		headerUserID   string
		tweetIDPathVar string
		mockService    *fakeLikeTweetService
		expectedStatus int
	}{
		{
			name:           "successfully likes a tweet",
			headerUserID:   "usr_123",
			tweetIDPathVar: "tweet_abc",
			mockService:    &fakeLikeTweetService{},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "missing X-User-ID header",
			headerUserID:   "",
			tweetIDPathVar: "tweet_abc",
			mockService:    &fakeLikeTweetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing tweet ID path param",
			headerUserID:   "usr_123",
			tweetIDPathVar: "",
			mockService:    &fakeLikeTweetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "use case error",
			headerUserID:   "usr_123",
			tweetIDPathVar: "tweet_abc",
			mockService:    &fakeLikeTweetService{Err: errors.New("already liked")},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/tweets/{id}/like", nil)
			if tc.headerUserID != "" {
				req.Header.Set("X-User-ID", tc.headerUserID)
			}

			req = mux.SetURLVars(req, map[string]string{"id": tc.tweetIDPathVar})
			rr := httptest.NewRecorder()

			handler := NewLikeTweetHandler(tc.mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)
		})
	}
}
