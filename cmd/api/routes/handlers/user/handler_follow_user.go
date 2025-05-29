package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"ualaTwitter/internal/platform/httphelper"
	"ualaTwitter/internal/usecase/follow_user"
)

const maxFollowBodySize = 4 * 1024 // 4 KB

var (
	ErrMissingUserID   = errors.New("missing X-User-ID header")
	ErrInvalidBody     = errors.New("invalid request body")
	ErrEmptyFolloweeID = errors.New("followee_id cannot be empty")
)

type followUserService interface {
	Execute(ctx context.Context, input follow_user.Input) error
}

type FollowUserHandler struct {
	service followUserService
}

func NewFollowUserHandler(service followUserService) *FollowUserHandler {
	return &FollowUserHandler{
		service: service,
	}
}

func (h *FollowUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, maxFollowBodySize)

	input, err := h.parseRequest(r)
	if err != nil {
		httphelper.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Execute(ctx, *input); err != nil {
		httphelper.RenderError(w, httphelper.StatusFromError(err), err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FollowUserHandler) parseRequest(r *http.Request) (*follow_user.Input, error) {
	followerID := r.Header.Get("X-User-ID")

	if followerID == "" {
		return nil, ErrMissingUserID
	}

	var req followUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, ErrInvalidBody
	}

	if req.FolloweeID == "" {
		return nil, ErrEmptyFolloweeID
	}

	return &follow_user.Input{
		FollowerID: followerID,
		FolloweeID: req.FolloweeID,
	}, nil
}
