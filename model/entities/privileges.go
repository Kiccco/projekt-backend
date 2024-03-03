package entities

type Privileges struct {
	UserID int `gorm:"primarykey"`
	User   User
	Admin  bool
}
