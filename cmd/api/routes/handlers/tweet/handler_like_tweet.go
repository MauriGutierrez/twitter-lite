package tweet

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"ualaTwitter/internal/platform/httphelper"
	"ualaTwitter/internal/usecase/like_tweet"
)

var (
	ErrMissingTweetID = errors.New("missing tweet ID in path")
)

type likeTweetService interface {
	Execute(ctx context.Context, input like_tweet.Input) error
}

type LikeTweetHandler struct {
	service likeTweetService
}

func NewLikeTweetHandler(service likeTweetService) *LikeTweetHandler {
	return &LikeTweetHandler{
		service: service,
	}
}

func (h *LikeTweetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

func (h *LikeTweetHandler) parseRequest(r *http.Request) (*like_tweet.Input, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return nil, ErrMissingUserID
	}

	vars := mux.Vars(r)
	tweetID := vars["id"]
	if tweetID == "" {
		return nil, ErrMissingTweetID
	}

	return &like_tweet.Input{
		UserID:  userID,
		TweetID: tweetID,
	}, nil
}
