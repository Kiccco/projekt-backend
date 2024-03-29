package friends

import "github.com/gofiber/contrib/websocket"

type LogoutEvent struct {
	UserID int
}

type LoginEvent struct {
	UserID     int
	Connection *websocket.Conn
}

type FriendListEvent struct {
	UserID int
	Users  []int
}

type FriendRequestEvent struct {
	UserID int
	Friend int
}

type FriendAcceptEvent struct {
	UserID    int
	NewFriend int
}

type FriendDeclineEvent struct {
	UserID int
	Friend int
}
