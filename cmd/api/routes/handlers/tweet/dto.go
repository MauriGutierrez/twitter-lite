package tweet

type postTweetRequest struct {
	Content string `json:"content"`
}

type postTweetResponse struct {
	ID string `json:"id"`
}

type tweetTimelineResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	Likes     int    `json:"likes"`
}
