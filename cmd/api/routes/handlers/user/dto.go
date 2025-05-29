package user

type followUserRequest struct {
	FolloweeID string `json:"followee_id"`
}

type createUserRequest struct {
	Name     string `json:"name"`
	Document string `json:"document"`
}

type createUserResponse struct {
	ID string `json:"id"`
}
