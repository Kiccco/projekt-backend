package handler

import (
	"backend/main/database"
	"backend/main/model/entities"
	"backend/main/model/responses"

	"github.com/gofiber/fiber/v2"
)

func SearchUser(c *fiber.Ctx) error {
	user := c.AllParams()
	var users []entities.User
	db := database.DB

	db.Where("username LIKE ?", "%"+user["name"]+"%").Find(&users)

	var res = []responses.StructUser{}

	for i := 0; i < len(users); i++ {
		res = append(res, responses.StructUser{ID: users[i].ID, User: users[i].Username, Mail: users[i].Mail})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "count": len(users), "result": res})
}
