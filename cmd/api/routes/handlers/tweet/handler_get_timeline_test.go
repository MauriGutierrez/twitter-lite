package tweet

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"ualaTwitter/internal/usecase/get_timeline"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGetTimelineService struct {
	Output []get_timeline.TweetTimeline
	Err    error
}

func (f *fakeGetTimelineService) Execute(_ context.Context, _ get_timeline.Input) ([]get_timeline.TweetTimeline, error) {
	return f.Output, f.Err
}

func TestGetTimelineHandler(t *testing.T) {
	tweetTime := time.Now()
	mockTweets := []get_timeline.TweetTimeline{
		{
			ID:        "tweet_1",
			UserID:    "usr_456",
			Content:   "Ualá tweeting",
			CreatedAt: tweetTime,
		},
		{
			ID:        "tweet_2",
			UserID:    "usr_789",
			Content:   "Second tweet",
			CreatedAt: tweetTime,
		},
	}

	tests := []struct {
		name           string
		headerUserID   string
		queryParams    string
		mockService    *fakeGetTimelineService
		expectedStatus int
		expectJSON     bool
		expectedBody   []map[string]string
	}{
		{
			name:         "returns timeline successfully",
			headerUserID: "usr_123",
			queryParams:  "?limit=2&offset=0",
			mockService: &fakeGetTimelineService{
				Output: mockTweets,
			},
			expectedStatus: http.StatusOK,
			expectJSON:     true,
			expectedBody: []map[string]string{
				{"id": "tweet_1", "user_id": "usr_456", "content": "Ualá tweeting", "created_at": tweetTime.Format(time.RFC3339)},
				{"id": "tweet_2", "user_id": "usr_789", "content": "Second tweet", "created_at": tweetTime.Format(time.RFC3339)},
			},
		},
		{
			name:           "missing X-User-ID header",
			headerUserID:   "",
			queryParams:    "",
			mockService:    &fakeGetTimelineService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "use case returns error",
			headerUserID: "usr_123",
			queryParams:  "",
			mockService: &fakeGetTimelineService{
				Err: errors.New("something failed"),
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:         "uses default limit and offset when query params are missing",
			headerUserID: "usr_123",
			queryParams:  "",
			mockService: &fakeGetTimelineService{
				Output: mockTweets,
			},
			expectedStatus: http.StatusOK,
			expectJSON:     true,
			expectedBody: []map[string]string{
				{"id": "tweet_1", "user_id": "usr_456", "content": "Ualá tweeting", "created_at": tweetTime.Format(time.RFC3339)},
				{"id": "tweet_2", "user_id": "usr_789", "content": "Second tweet", "created_at": tweetTime.Format(time.RFC3339)},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/timeline"+tc.queryParams, nil)
			if tc.headerUserID != "" {
				req.Header.Set("X-User-ID", tc.headerUserID)
			}

			rr := httptest.NewRecorder()

			handler := NewGetTimelineHandler(tc.mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectJSON {
				var parsed []map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &parsed)
				require.NoError(t, err)
				require.Equal(t, tc.expectedBody, parsed)
			}
		})
	}
}
