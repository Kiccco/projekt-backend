package main

import (
	"backend/main/database"
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
	log.Printf("Povezan z bazo\n")
	log.Println(database.DB)
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	router.SetupRoutes(app)

	log.Printf("Server laufa z verzijo: %s\n", version)
	log.Panic(app.Listen(":8080"))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
}
