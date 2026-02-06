package core

type User struct {
	Id           int    `json:"id"`
	Username     string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type Group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type GroupMember struct {
	UserId  int `json:"user_id"`
	GroupId int `json:"group_id"`
}

type Message struct {
	Id       int    `json:"id"`
	UserId   int    `json:"user_id"`
	GroupId  int    `json:"group_id"`
	Content  string `json:"content"`
	MediaUrl string `json:"media_url,omitempty"`
}
