package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterRoutes(r *mux.Router, h Handlers) {
	r.HandleFunc("/tweets", h.PostTweet).Methods(http.MethodPost)
	r.HandleFunc("/timeline", h.GetTimeline).Methods(http.MethodGet)
	r.HandleFunc("/tweets/{id}/like", h.LikeTweet).Methods(http.MethodPost)
	r.HandleFunc("/follow", h.FollowUser).Methods(http.MethodPost)
	r.HandleFunc("/users", h.CreateUser).Methods(http.MethodPost)

	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)
}
