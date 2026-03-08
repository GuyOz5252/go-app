package websocket

type WSMessage struct {
	Type              WSMessageType `json:"message_type"`
	InitiatorUserId   string        `json:"initiator_user_id,omitempty"`
	DestinationChatId string        `json:"destination_chat_id,omitempty"`
	DestinationUserId string        `json:"destination_user_id,omitempty"`
	Payload           any           `json:"payload,omitempty"`
}

type WSMessageType int

const (
	NewMessage WSMessageType = iota
	MessageServerAck
	MessageUserAck
	MessageUserReadAck
	UserTypingStart
	UserTypingEnd
	UserOnline
	UserOffline
	UserAway
	ServerError
)

type NewMessagePayload struct {
	Content   string `json:"content" mapstructure:"content"`
	MediaUrl  string `json:"media_url" mapstructure:"media_url"`
	ReplyToId string `json:"reply_to_id" mapstructure:"reply_to_id"`
}

type MessageIdPayload struct {
	MessageId string `json:"message_id"`
}

type ServerErrorPayload struct {
	Error string `json:"error"`
}
