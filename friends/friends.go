package friends

import (
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	activeFriends   map[int][]int
	userConnections map[int]*websocket.Conn
	FriendChannel   chan interface{}
)

func Init() {
	activeFriends = make(map[int][]int)
	userConnections = make(map[int]*websocket.Conn)
	FriendChannel = make(chan interface{})
	go loop()
}

func loop() {
	for {
		select {
		case message := <-FriendChannel:
			switch message.(type) {
			case LoginEvent:
				log.Println("Login event")
				event := message.(LoginEvent)
				AddUser(event.UserID)
				AddUserConnection(event.UserID, event.Connection)
				break
			case LogoutEvent:
				event := message.(LogoutEvent)

				for _, friendID := range GetFriends(event.UserID) {
					if connection, ok := userConnections[friendID]; ok {
						connection.WriteJSON(fiber.Map{"type": "logout", "userID": event.UserID})
					}
				}
				RemoveUserConnection(event.UserID)
				RemoveUser(event.UserID)
				break
			case FriendListEvent:
				event := message.(FriendListEvent)
				log.Println(event.Users)
				activeFriends[event.UserID] = append(activeFriends[event.UserID], event.Users...)

				for _, friendID := range GetFriends(event.UserID) {
					if connection, ok := userConnections[friendID]; ok {
						connection.WriteJSON(fiber.Map{"type": "login", "userID": event.UserID})
						userConnections[event.UserID].WriteJSON(fiber.Map{"type": "login", "userID": friendID})
					}
				}
				break
			case FriendRequestEvent:
				event := message.(FriendRequestEvent)

				if connection, ok := userConnections[event.Friend]; ok {
					connection.WriteJSON(fiber.Map{"type": "friendRequest", "userID": event.UserID})
				}
				break
			case FriendAcceptEvent:
				event := message.(FriendAcceptEvent)
				AddFriend(event.UserID, event.NewFriend)
				if connection, ok := userConnections[event.UserID]; ok {
					connection.WriteJSON(fiber.Map{"type": "friendAccept", "userID": event.NewFriend})
				}
				break
			case FriendDeclineEvent:
				event := message.(FriendDeclineEvent)
				if connection, ok := userConnections[event.UserID]; ok {
					connection.WriteJSON(fiber.Map{"type": "friendDecline", "userID": event.Friend})
				}
				break
			}

		}
	}
}

func AddFriend(userID int, friendID int) {
	if activeFriends[userID] == nil {
		activeFriends[userID] = append(activeFriends[userID], friendID)
	}

	if activeFriends[friendID] == nil {
		activeFriends[friendID] = append(activeFriends[friendID], userID)

	}
}

func AddUser(userID int) {
	activeFriends[userID] = []int{}
}

func RemoveFriend(userID int, friendID int) {
	for i, id := range activeFriends[userID] {
		if id == friendID {
			activeFriends[userID] = append(activeFriends[userID][:i], activeFriends[userID][i+1:]...)
		}
	}

	for i, id := range activeFriends[friendID] {
		if id == userID {
			activeFriends[friendID] = append(activeFriends[friendID][:i], activeFriends[friendID][i+1:]...)
		}
	}
}

func GetFriends(userID int) []int {
	return activeFriends[userID]
}

func AddUserConnection(userID int, connection *websocket.Conn) {
	userConnections[userID] = connection
}

func GetUserConnection(userID int) *websocket.Conn {
	return userConnections[userID]
}

func RemoveUserConnection(userID int) {
	delete(userConnections, userID)
}

func RemoveUser(userID int) {
	delete(activeFriends, userID)
}
