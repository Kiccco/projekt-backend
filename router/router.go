package router

import (
	"backend/main/handler"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	auth := api.Group("/auth")

	auth.Post("/login", handler.Login)
	auth.Post("/register", handler.Register)

	app.Use(jwtware.New(jwtware.Config{SigningKey: jwtware.SigningKey{Key: []byte("tojevelikporazinupamdabokmalbolje")}}))

}
