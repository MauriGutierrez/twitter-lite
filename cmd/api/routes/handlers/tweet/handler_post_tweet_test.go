package tweet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"ualaTwitter/internal/usecase/post_tweet"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePostTweetService struct {
	TweetID string
	Err     error
}

func (f *fakePostTweetService) Execute(_ context.Context, _ post_tweet.Input) (string, error) {
	return f.TweetID, f.Err
}

func TestPostTweetHandler(t *testing.T) {
	tests := []struct {
		name           string
		header         string
		body           map[string]string
		mockService    *fakePostTweetService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "successfully posts tweet",
			header: "usr_123",
			body: map[string]string{
				"content": "Hello Twitter!",
			},
			mockService: &fakePostTweetService{
				TweetID: "tweet_001",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"tweet_001"}`,
		},
		{
			name:           "missing X-User-ID",
			header:         "",
			body:           map[string]string{"content": "Hello"},
			mockService:    &fakePostTweetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "empty request body",
			header:         "usr_123",
			body:           nil,
			mockService:    &fakePostTweetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "empty tweet content",
			header: "usr_123",
			body: map[string]string{
				"content": "",
			},
			mockService:    &fakePostTweetService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "use case returns error",
			header: "usr_123",
			body: map[string]string{
				"content": "Something real",
			},
			mockService: &fakePostTweetService{
				Err: errors.New("internal failure"),
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

			req := httptest.NewRequest(http.MethodPost, "/tweets", bytes.NewReader(bodyBytes))
			if tc.header != "" {
				req.Header.Set("X-User-ID", tc.header)
			}

			rr := httptest.NewRecorder()

			handler := NewPostTweetHandler(tc.mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus == http.StatusOK {
				require.JSONEq(t, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
