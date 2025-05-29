package get_timeline

import "time"

type TweetTimeline struct {
	ID        string
	UserID    string
	Content   string
	Likes     int
	CreatedAt time.Time
}
