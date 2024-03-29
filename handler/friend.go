package handler

import (
	"backend/main/database"
	"backend/main/model/entities"
	"backend/main/model/requests"
	"backend/main/model/responses"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetFriends(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	db := database.DB
	var users []responses.StructFriend
	db.Table("users").Select("users.id, users.username, users.mail, friends.chat_uuid").Joins("join friends on friends.user_id = users.id").Where("friends.friend_id = ? AND friends.deleted_at IS NULL", claims["id"]).Scan(&users)
	res := []responses.StructUser{}

	for _, user := range users {
		res = append(res, responses.StructUser{ID: user.ID, User: user.Username, Mail: user.Mail})
	}

	return c.Status(fiber.StatusOK).JSON(users)

}

func AddFriend(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.AddFriend

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	//treba bo ksn channel narest da opozori drug user
	users := []entities.Friends{{UserID: int(claims["id"].(float64)), FriendID: token.ID}, {UserID: token.ID, FriendID: int(claims["id"].(float64))}}
	db := database.DB
	rowsAffected := db.Create(&users).RowsAffected
	if rowsAffected < 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func RemoveFriend(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.RemoveFriend

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	rowsAffected := db.Where("user_id = ?", int(claims["id"].(float64))).Delete(&entities.Friends{}).RowsAffected
	if rowsAffected <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	rowsAffected = db.Where("user_id = ?", token.ID).Delete(&entities.Friends{}).RowsAffected
	if rowsAffected <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func GetFriendRequests(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	db := database.DB
	var users []entities.User
	db.Model(&entities.User{}).Select("users.id, users.username, users.mail").Joins("join friend_requests on friend_requests.sender_id = users.id").Where("friend_requests.receiver_id = ? AND friend_requests.deleted_at IS NULL", claims["id"]).Scan(&users)

	res := []responses.StructUser{}

	for _, user := range users {
		res = append(res, responses.StructUser{ID: user.ID, User: user.Username, Mail: user.Mail})
	}

	return c.Status(fiber.StatusOK).JSON(res)

}

func SendFriendRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.AddFriend

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	userr := entities.FriendRequests{SenderID: int(claims["id"].(float64)), ReceiverID: token.ID}
	db := database.DB
	rowsAffected := db.Create(&userr).RowsAffected
	if rowsAffected < 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func RemoveFriendRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.AcceptFriendRequest

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	rowsAffected := db.Where("sender_id = ? AND receiver_id = ?", int(claims["id"].(float64)), token.ID).Delete(&entities.FriendRequests{}).RowsAffected
	if rowsAffected <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func AcceptFriendReq(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.AcceptFriendRequest

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	rowsAffected := db.Where("sender_id = ? AND receiver_id = ?", token.ID, int(claims["id"].(float64))).Delete(&entities.FriendRequests{}).RowsAffected
	if rowsAffected <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	uuid := uuid.New().String()

	users := []entities.Friends{{UserID: int(claims["id"].(float64)), FriendID: token.ID, ChatUUID: uuid}, {UserID: token.ID, FriendID: int(claims["id"].(float64)), ChatUUID: uuid}}
	rowsAffected = db.Create(&users).RowsAffected
	if rowsAffected < 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

func DeclineFriendReq(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	var token requests.DeclineFriendRequest

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	rowsAffected := db.Where("sender_id = ? AND receiver_id = ?", token.ID, int(claims["id"].(float64))).Delete(&entities.FriendRequests{}).RowsAffected
	if rowsAffected <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
