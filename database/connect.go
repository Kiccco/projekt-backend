package database

import (
	"backend/main/model/entities"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenConnection(host string, user string, pass string, dbname string, port uint) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Panicf("Napaka pri povezavi do podatkovne baze: %s", err)
	}
	DB = db

	DB.AutoMigrate(&entities.User{}, &entities.Privileges{})

}
