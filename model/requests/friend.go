package requests

type AddFriend struct {
	ID   int    `json:"id"`
	user string `json:"user"`
}

type RemoveFriend struct {
	ID   int    `json:"id"`
	user string `json:"user"`
}

type SendFriendRequest struct {
	ID   int    `json:"id"`
	user string `json:"user"`
}

type AcceptFriendRequest struct {
	ID int `json:"id"`
}

type DeclineFriendRequest struct {
	ID int `json:"id"`
}
