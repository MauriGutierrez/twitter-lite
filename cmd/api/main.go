package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"ualaTwitter/cmd/api/config"
	"ualaTwitter/cmd/api/routes/handlers/health"
	"ualaTwitter/internal/platform/logger"
	"ualaTwitter/internal/platform/repository/memory"
	"ualaTwitter/internal/platform/repository/postgres"
	"ualaTwitter/internal/usecase/create_user"

	"ualaTwitter/cmd/api/routes"
	"ualaTwitter/cmd/api/routes/handlers/tweet"
	"ualaTwitter/cmd/api/routes/handlers/user"

	"ualaTwitter/internal/usecase/follow_user"
	"ualaTwitter/internal/usecase/get_timeline"
	"ualaTwitter/internal/usecase/like_tweet"
	"ualaTwitter/internal/usecase/post_tweet"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()
	conn := initializePsx(ctx, cfg.PostgresDSN)

	// === Repositories ===
	memoryUserRepository := memory.NewInMemoryUserRepository()
	psxUserRepository := postgres.NewPostgresUserRepository(conn)
	tweetRepo := memory.NewInMemoryTweetRepository()
	likeRepo := memory.NewInMemoryLikeRepository()

	// === Usecases ===
	postTweetService := post_tweet.NewPostTweetService(tweetRepo, memoryUserRepository)
	followUserService := follow_user.NewFollowUserService(memoryUserRepository)
	getTimelineService := get_timeline.NewGetTimelineService(tweetRepo, memoryUserRepository)
	createUserService := create_user.NewCreateUserService(psxUserRepository, memoryUserRepository)
	likeTweetService := like_tweet.NewLikeTweetService(tweetRepo, likeRepo)

	// === Handlers ===
	postTweetHandler := tweet.NewPostTweetHandler(postTweetService)
	followUserHandler := user.NewFollowUserHandler(followUserService)
	getTimelineHandler := tweet.NewGetTimelineHandler(getTimelineService)
	createUserHandler := user.NewCreateUserHandler(createUserService)
	likeTweetHandler := tweet.NewLikeTweetHandler(likeTweetService)

	healthHandler := health.NewHealthHandler(cfg.Env, cfg.AppName, cfg.Version)

	logger.Init()

	// === Route Bindings ===
	handlers := routes.Handlers{
		PostTweet:   postTweetHandler.ServeHTTP,
		FollowUser:  followUserHandler.ServeHTTP,
		GetTimeline: getTimelineHandler.ServeHTTP,
		CreateUser:  createUserHandler.ServeHTTP,
		LikeTweet:   likeTweetHandler.ServeHTTP,
		Health:      healthHandler.ServeHTTP,
	}

	r := mux.NewRouter()
	routes.RegisterRoutes(r, handlers)
	log.Printf("Server started on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}

func initializePsx(ctx context.Context, dsn string) *pgx.Conn {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	return conn
}
