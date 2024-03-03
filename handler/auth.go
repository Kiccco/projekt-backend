package handler

import (
	"backend/main/database"
	"backend/main/model/entities"
	"backend/main/model/requests"
	"errors"
	"log"
	"net/mail"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(c *fiber.Ctx) error {

	var token requests.UserLogin

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	var user entities.User
	err := db.Where(&entities.User{Username: token.Username}).First(&user).Error

	// preveri če userja NI notr u bazi ter če je kakšna druga napaka
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "authbadc"})
	} else if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	log.Println(token)
	//primerja geslo v bazi in tistega k ga je poslau
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(token.Password))
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "authbadc"})
	}

	//jwt token
	claims := jwt.MapClaims{
		"username": token.Username,
		"email":    user.Mail,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	} //TODO: dodat kasn rank/role ma

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	res, err := jwtToken.SignedString([]byte("tojevelikporazinupamdabokmalbolje"))
	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "code": "auth", "token": res})
}

func Register(c *fiber.Ctx) error {

	var token requests.UserCreate

	if err := c.BodyParser(&token); err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "reqbad"})
	}

	db := database.DB
	var user entities.User

	//Prever če je username zaseden
	err := db.Where(&entities.User{Username: token.Username}).First(&user).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "userdupl"})
	}
	log.Println(token)
	//prever če je mail veljaven
	_, err = mail.ParseAddress(token.Mail)

	if err != nil {
		log.Println(err)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "mailbadc"})
	}

	//Prever če je email že v uporabi
	err = db.Where(&entities.User{Mail: token.Mail}).First(&entities.User{}).Error
	if err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "error", "code": "maildupl"})
	}

	//Če je use gud poj ustvar nov racun, najprej naredi hash od gesla
	pass, err := bcrypt.GenerateFromPassword([]byte(token.Password), 14)
	if err != nil {
		log.Printf("Napaka pri bcrypt.GenerateFromPassword(): %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	//zapiš ga notr u bazo
	rows := db.Create(&entities.User{Username: token.Username, Mail: token.Mail, Password: string(pass)}).RowsAffected
	if rows <= 0 {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "code": "acct"})
}
