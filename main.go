package main

import (
	"backend/main/chat"
	"backend/main/database"
	"backend/main/friends"
	"backend/main/router"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const (
	version = "0.0.1"
	dbUser  = "postgres"
	dbPass  = "test"
	dbHost  = "localhost"
	dbPort  = 5432
	dbName  = "chat"
)

func main() {

	database.OpenConnection(dbHost, dbUser, dbPass, dbName, dbPort)
	chat.Init()
	friends.Init()
	log.Printf("Povezan z bazo\n")
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "*",
	}))

	router.SetupRoutes(app)

	log.Panic(app.Listen(":8080"))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
