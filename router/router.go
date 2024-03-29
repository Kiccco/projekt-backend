package router

import (
	"backend/main/handler"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	auth := api.Group("/auth")
	api.Get("/search/:name", handler.SearchUser)

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	api.Use("/chat", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	api.Get("/chat", websocket.New(handler.Chat))

	app.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte("tojevelikporazinupamdabokmalbolje")}}))

	api.Get("/getFriends", handler.GetFriends)
	api.Post("/addFriend", handler.SendFriendRequest)
	api.Post("/removeFriend", handler.RemoveFriend)

	api.Get("/getFriendRequests", handler.GetFriendRequests)
	api.Post("/acceptFriendReq", handler.AcceptFriendReq)
	api.Post("/removeFriendReq", handler.RemoveFriendRequest)
	api.Post("/declineFriendReq", handler.DeclineFriendReq)

}
