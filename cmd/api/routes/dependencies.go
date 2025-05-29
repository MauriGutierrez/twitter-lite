package routes

import (
	"net/http"
)

type Handlers struct {
	PostTweet   http.HandlerFunc
	FollowUser  http.HandlerFunc
	CreateUser  http.HandlerFunc
	GetTimeline http.HandlerFunc
	LikeTweet   http.HandlerFunc
	Health      http.HandlerFunc
}
