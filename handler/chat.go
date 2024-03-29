package handler

import (
	"backend/main/chat"
	"backend/main/chat/events"
	"backend/main/friends"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ChatMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChatUUID struct {
	UUID string `json:"uuid"`
}

type ChatMessageData struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
}

type FriendList struct {
	Friends []int `json:"friends"`
}

type FriendRequest struct {
	Friend int `json:"friend"`
}

type FriendAccept struct {
	NewFriend int `json:"friend"`
}

func Chat(c *websocket.Conn) {

	token := c.Query("token")

	validateToken, erro := validateToken(token)
	if erro != nil {
		log.Println("Token validation failed")
		c.Close()
		return
	}

	c.Locals("user", validateToken["username"])
	c.Locals("id", int(validateToken["id"].(float64)))
	friends.FriendChannel <- friends.LoginEvent{UserID: int(validateToken["id"].(float64)), Connection: c}

	for {
		var chatMessage ChatMessage
		err := c.ReadJSON(&chatMessage)

		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
				log.Println("Websocket closed")
				friends.FriendChannel <- friends.LogoutEvent{UserID: int(validateToken["id"].(float64))}
			} else {

				log.Println("read:", err)
			}
			break
		}

		switch chatMessage.Type {
		case "getMessages":
			var chatUUID ChatUUID
			err := json.Unmarshal([]byte(chatMessage.Data), &chatUUID)

			if err != nil {
				log.Println("error parsing json type: getMessages, error:", err)
				break
			}

			uuid, err := uuid.Parse(chatUUID.UUID)
			if err != nil {
				log.Printf("error parsing uuid: %s, error: %s", chatUUID.UUID, err)
				break
			}

			chat.ChatChannel <- events.GetMessagesEvent{Connection: c, UUID: uuid, UserID: int(validateToken["id"].(float64))}
			break
		case "sendMessage":
			var msg ChatMessageData
			err := json.Unmarshal(chatMessage.Data, &msg)
			if err != nil {
				log.Println("error parsing json type: sendMessage, error:", err)
				break
			}

			uuid, err := uuid.Parse(msg.UUID)
			if err != nil {
				log.Printf("error parsing uuid: %s, error: %s", msg.UUID, err)
				break
			}

			chat.ChatChannel <- events.SendMessageEvent{Connection: c, UUID: uuid, Message: msg.Message}
			break
		//Kle se hendlam use zadeve v povezavi z frendi... Žalost.
		case "friendList":
			var friendList FriendList
			err := json.Unmarshal(chatMessage.Data, &friendList)
			if err != nil {
				log.Println("error parsing json type: friendList, error:", err)
				break
			}
			friends.FriendChannel <- friends.FriendListEvent{UserID: int(validateToken["id"].(float64)), Users: friendList.Friends}
			break
		case "friendRequest":
			var friendRequest FriendRequest
			err := json.Unmarshal(chatMessage.Data, &friendRequest)
			if err != nil {
				log.Println("error parsing json type: friendRequest, error:", err)
				break
			}
			friends.FriendChannel <- friends.FriendRequestEvent{UserID: int(validateToken["id"].(float64)), Friend: friendRequest.Friend}
			break
		case "friendAccept":
			var friendAccept FriendAccept
			err := json.Unmarshal(chatMessage.Data, &friendAccept)
			if err != nil {
				log.Println("error parsing json type: friendAccept, error:", err)
				break
			}
			friends.FriendChannel <- friends.FriendAcceptEvent{UserID: int(validateToken["id"].(float64)), NewFriend: friendAccept.NewFriend}
		case "friendDecline":
			var friendDecline FriendAccept
			err := json.Unmarshal(chatMessage.Data, &friendDecline)
			if err != nil {
				log.Println("error parsing json type: friendDecline, error:", err)
				break
			}
			friends.FriendChannel <- friends.FriendDeclineEvent{UserID: int(validateToken["id"].(float64)), Friend: friendDecline.NewFriend}
		default:
			log.Printf("Received unknown message type from %s", c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["username"])
		}
	}
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

	signingKey := []byte("tojevelikporazinupamdabokmalbolje")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		return claims, nil
	} else {
		return nil, err
	}
}
