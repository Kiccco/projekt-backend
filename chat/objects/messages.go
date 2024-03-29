package objects

type MessageFile struct {
	ReadMessages []Message `json:"messages"`
}

type Message struct {
	UserID  int    `json:"id"`
	Content string `json:"msg"`
	Time    string `json:"ura"`
	Prejel  bool   `json:"prejel"`
}
