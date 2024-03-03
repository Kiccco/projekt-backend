package entities

import "gorm.io/gorm"

type UserInfo struct {
	FirstName string
	LastName  string
	Phone     string
	Birth     float32
}

type User struct {
	gorm.Model
	Username string
	Mail     string
	Password string
	Phone    string
	Birth    float32
}
