package tweet

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"ualaTwitter/internal/platform/httphelper"
	"ualaTwitter/internal/usecase/get_timeline"
)

const (
	defaultLimitValue  = 50
	defaultOffsetValue = 0
)

type getTimelineService interface {
	Execute(ctx context.Context, input get_timeline.Input) ([]get_timeline.TweetTimeline, error)
}

type GetTimelineHandler struct {
	service getTimelineService
}

func NewGetTimelineHandler(service getTimelineService) *GetTimelineHandler {
	return &GetTimelineHandler{
		service: service,
	}
}

func (h *GetTimelineHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	input, err := h.parseRequest(r)
	if err != nil {
		httphelper.RenderError(w, http.StatusBadRequest, err.Error())
		return
	}

	tweets, err := h.service.Execute(ctx, *input)
	if err != nil {
		httphelper.RenderError(w, httphelper.StatusFromError(err), err.Error())
		return
	}

	h.renderResponse(w, tweets)
}

func (h *GetTimelineHandler) parseRequest(r *http.Request) (*get_timeline.Input, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return nil, ErrMissingUserID
	}

	limit := parseQueryInt(r, "limit", defaultLimitValue)
	offset := parseQueryInt(r, "offset", defaultOffsetValue)

	return &get_timeline.Input{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (h *GetTimelineHandler) renderResponse(w http.ResponseWriter, tweets []get_timeline.TweetTimeline) {
	w.Header().Set("Content-Type", "application/json")

	response := make([]tweetTimelineResponse, len(tweets))
	for i, t := range tweets {
		response[i] = tweetTimelineResponse{
			ID:        t.ID,
			UserID:    t.UserID,
			Content:   t.Content,
			Likes:     t.Likes,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode timeline response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func parseQueryInt(r *http.Request, key string, defaultVal int) int {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil || val < 0 {
		return defaultVal
	}

	return val
}
