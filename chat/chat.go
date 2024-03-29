package chat

import (
	"backend/main/chat/events"
	"backend/main/chat/objects"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	activeChats     map[uuid.UUID]objects.MessageFile
	users           map[uuid.UUID][]int
	ChatChannel     chan interface{}
	usersConnection map[int]*websocket.Conn
)

func Init() {
	activeChats = make(map[uuid.UUID]objects.MessageFile)
	ChatChannel = make(chan interface{})
	users = make(map[uuid.UUID][]int)
	usersConnection = make(map[int]*websocket.Conn)
	go loop()
}

func loop() {
	for {
		select {
		case message := <-ChatChannel:
			switch message.(type) {
			case events.GetMessagesEvent:
				event := message.(events.GetMessagesEvent)
				messages := GetMessages(event.UUID)
				event.Connection.WriteJSON(fiber.Map{"type": "message", "chatUUID": event.UUID.String(), "messages": messages})
				AddUser(event.UUID, event.UserID)
				AddUserConnection(event.UserID, event.Connection)

				break
			case events.SendMessageEvent:
				event := message.(events.SendMessageEvent)
				AddMessage(event.UUID, objects.Message{UserID: event.Connection.Locals("id").(int), Content: event.Message, Time: time.Now().Format("15:04"), Prejel: false})

				chatUUID := event.UUID
				for _, userID := range GetChatUsers(chatUUID) {
					log.Println("Sending message to user", userID)
					if connection, ok := usersConnection[userID]; ok {
						connection.WriteJSON(fiber.Map{"type": "messageSent", "chatUUID": chatUUID.String(), "message": event.Message, "time": time.Now().Format("15:04"), "userID": event.Connection.Locals("id").(int)})
					}
				}
				break
			case events.LogoutEvent:
				event := message.(events.LogoutEvent)
				chatUUIDs := GetChatUUID(event.UserID)
				for _, chatUUID := range chatUUIDs {
					for _, user := range GetChatUsers(chatUUID) {
						if connection, ok := usersConnection[user]; ok {
							connection.WriteJSON(fiber.Map{"type": "logout", "id": event.UserID})
						}
					}

					RemoveUser(chatUUID, event.UserID)
					if GetChatUsersCount(chatUUID) == 0 {
						RemoveChat(chatUUID)
					}
				}
				RemoveUserConnection(event.UserID)
			}

		}
	}

}

func AddChat(chatUUID uuid.UUID) {

	if _, ok := activeChats[chatUUID]; ok {
		log.Printf("Chat %s already exists\n", chatUUID.String())
		return
	}

	jsonFile, err := os.OpenFile(chatUUID.String()+".json", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var messages objects.MessageFile
	err = json.Unmarshal(byteValue, &messages)
	if err != nil {
		log.Printf("Error reading a %s.json. err: %s resetting..", chatUUID.String(), err)

		messages = objects.MessageFile{ReadMessages: []objects.Message{}}
		data, _ := json.MarshalIndent(messages, "", " ")
		ioutil.WriteFile(chatUUID.String()+".json", data, 0644)
	}
	activeChats[chatUUID] = messages
}

func SaveChat(chatUUID uuid.UUID) {
	if _, ok := activeChats[chatUUID]; !ok {
		log.Printf("Chat %s does not exist\n", chatUUID.String())
		return
	}

	jsonFile, err := os.Create(chatUUID.String() + ".json")
	if err != nil {
		log.Println("Error creating file:", err)
	}
	defer jsonFile.Close()

	byteValue, _ := json.Marshal(activeChats[chatUUID])
	jsonFile.Write(byteValue)

	delete(activeChats, chatUUID)
}

func GetActiveChats() map[uuid.UUID]objects.MessageFile {
	return activeChats
}

func AddMessage(chatUUID uuid.UUID, message objects.Message) {
	if entry, ok := activeChats[chatUUID]; ok {
		entry.ReadMessages = append(entry.ReadMessages, message)
		activeChats[chatUUID] = entry
	} else {
		log.Printf("Chat %s is not active when trying adding a message\n", chatUUID.String())
	}
}

func GetMessages(chatUUID uuid.UUID) []objects.Message {
	if _, ok := activeChats[chatUUID]; !ok {
		AddChat(chatUUID)
	}
	return activeChats[chatUUID].ReadMessages
}

func RemoveChat(chatUUID uuid.UUID) {
	delete(activeChats, chatUUID)
}

func GetUsers() map[uuid.UUID][]int {
	return users
}

func AddUser(chatUUID uuid.UUID, userID int) {

	if _, ok := users[chatUUID]; !ok {
		users[chatUUID] = []int{}
	}

	if IsUserInChat(chatUUID, userID) {
		return
	}

	users[chatUUID] = append(users[chatUUID], userID)
}

func RemoveUser(chatUUID uuid.UUID, userID int) {
	if _, ok := users[chatUUID]; ok {
		for i, id := range users[chatUUID] {
			if id == userID {
				users[chatUUID] = append(users[chatUUID][:i], users[chatUUID][i+1:]...)
				break
			}
		}
	}

	if len(users[chatUUID]) == 0 {
		delete(users, chatUUID)
		SaveChat(chatUUID)
	}
}

func GetChatUUID(userID int) []uuid.UUID {
	var chatUUIDs []uuid.UUID
	for chatUUID, users := range users {
		for _, id := range users {
			if id == userID {
				chatUUIDs = append(chatUUIDs, chatUUID)
				break
			}
		}
	}
	return chatUUIDs
}

func IsUserInChat(chatUUID uuid.UUID, userID int) bool {
	if _, ok := users[chatUUID]; ok {
		for _, id := range users[chatUUID] {
			if id == userID {
				return true
			}
		}
	}
	return false
}

func GetChatUsers(chatUUID uuid.UUID) []int {
	if _, ok := users[chatUUID]; ok {
		return users[chatUUID]
	}
	return []int{}
}

func GetChatUsersCount(chatUUID uuid.UUID) int {
	if _, ok := users[chatUUID]; ok {
		return len(users[chatUUID])
	}
	return 0
}

func GetChatUsersCountAll() int {
	return len(users)
}

func AddUserConnection(userID int, connection *websocket.Conn) {
	usersConnection[userID] = connection
}

func RemoveUserConnection(userID int) {
	delete(usersConnection, userID)
}
