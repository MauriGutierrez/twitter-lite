package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"ualaTwitter/internal/usecase/create_user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCreateUserService struct {
	Output create_user.Output
	Err    error
}

func (f *fakeCreateUserService) Execute(_ context.Context, input create_user.Input) (create_user.Output, error) {
	if f.Err != nil {
		return create_user.Output{}, f.Err
	}
	return f.Output, nil
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           map[string]string
		mockService    *fakeCreateUserService
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successfully creates user",
			body: map[string]string{
				"name":     "Mauricio",
				"document": "12345678",
			},
			mockService: &fakeCreateUserService{
				Output: create_user.Output{ID: "usr_12345678"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"usr_12345678"}`,
		},
		{
			name:           "empty request body",
			body:           nil,
			mockService:    &fakeCreateUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid DNI format",
			body: map[string]string{
				"name":     "Mauricio",
				"document": "abc123",
			},
			mockService:    &fakeCreateUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "empty user name",
			body: map[string]string{
				"name":     "",
				"document": "12345678",
			},
			mockService:    &fakeCreateUserService{},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "use case returns error",
			body: map[string]string{
				"name":     "Mauricio",
				"document": "12345678",
			},
			mockService: &fakeCreateUserService{
				Err: errors.New("something went wrong"),
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

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
			rr := httptest.NewRecorder()

			handler := NewCreateUserHandler(tc.mockService)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code)

			if tc.expectedStatus == http.StatusOK {
				require.JSONEq(t, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
