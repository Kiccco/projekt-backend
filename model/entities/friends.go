package entities

import "gorm.io/gorm"

type Friends struct {
	gorm.Model
	UserID   int
	FriendID int
	User     User
	Friend   User
	ChatUUID string
}

type FriendRequests struct {
	gorm.Model
	SenderID   int
	ReceiverID int
	Sender     User
	Receiver   User
}
