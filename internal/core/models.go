package core

import "time"

type User struct {
	Id           string `json:"id"`
	Username     string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
}

type Chat struct {
	Id            string    `json:"id"`
	Name          string    `json:"name"`
	ChatMemberIds []string  `json:"chat_member_ids"`
	CreatedAt     time.Time `json:"created_at"`
}

type ChatMember struct {
	UserId string `json:"user_id"`
	ChatId string `json:"chat_id"`
}

type Message struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	ChatId    string `json:"chat_id"`
	Content   string `json:"content"`
	MediaUrl  string `json:"media_url,omitempty"`
	ReplyToId string `json:"reply_to_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
