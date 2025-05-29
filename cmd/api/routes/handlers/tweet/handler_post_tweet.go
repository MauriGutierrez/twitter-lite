package tweet

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ualaTwitter/internal/platform/httphelper"
	"ualaTwitter/internal/usecase/post_tweet"
)

const maxPostTweetBodySize = 4 * 1024

var (
	ErrMissingUserID = errors.New("missing X-User-ID header")
	ErrInvalidBody   = errors.New("invalid request body")
	ErrEmptyTweet    = errors.New("tweet content cannot be empty")
)

type postTweetService interface {
	Execute(ctx context.Context, input post_tweet.Input) (string, error)
}

type PostTweetHandler struct {
	service postTweetService
}

func NewPostTweetHandler(service postTweetService) *PostTweetHandler {
	return &PostTweetHandler{
		service: service,
	}
}

func (h *PostTweetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input, err := h.parseRequest(r)
	if err != nil {
		httphelper.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	tweetID, err := h.service.Execute(ctx, *input)
	if err != nil {
		httphelper.RenderError(w, httphelper.StatusFromError(err), err.Error())
		return
	}

	h.renderResponse(w, tweetID)
}

func (h *PostTweetHandler) parseRequest(r *http.Request) (*post_tweet.Input, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return nil, ErrMissingUserID
	}

	r.Body = http.MaxBytesReader(nil, r.Body, maxPostTweetBodySize)

	var req postTweetRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		return nil, ErrInvalidBody
	}

	if req.Content == "" {
		return nil, ErrEmptyTweet
	}

	return &post_tweet.Input{
		UserID:  userID,
		Content: req.Content,
	}, nil
}

func (h *PostTweetHandler) renderResponse(w http.ResponseWriter, tweetID string) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(postTweetResponse{ID: tweetID}); err != nil {
		log.Printf("failed to encode post tweet response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
