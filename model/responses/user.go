package responses

type StructUser struct {
	ID   uint   `json:"id"`
	User string `json:"user"`
	Mail string `json:"mail"`
}

type StructFriend struct {
	ID       uint   `json:"id"`
	Username string `json:"user"`
	Mail     string `json:"mail"`
	ChatUUID string `json:"chatUUID"`
}
