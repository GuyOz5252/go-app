package data

import (
	"context"
	"database/sql"

	"github.com/GuyOz5252/go-app/internal/core"
)

type SqlChatRepository struct {
	db      *sql.DB
	cache   core.Cache
	queries map[string]string
}

func NewSqlChatRepository(db *sql.DB, queries map[string]string) core.ChatRepository {
	return &SqlChatRepository{
		db:      db,
		cache:   nil,
		queries: queries,
	}
}

func (r *SqlChatRepository) GetById(ctx context.Context, id string) (*core.Chat, error) {
	chat := &core.Chat{}
	query, ok := r.queries["get_by_id"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}
	var imageUrl sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&chat.Id,
		&chat.Name,
		&imageUrl,
		&chat.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	if imageUrl.Valid {
		chat.ImageUrl = imageUrl.String
	}
	return chat, nil
}

func (r *SqlChatRepository) ListByUserId(ctx context.Context, userId string) ([]*core.ChatDto, error) {
	query, ok := r.queries["list_by_user_id"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []*core.ChatDto
	for rows.Next() {
		chat := &core.ChatDto{}
		var imageUrl sql.NullString
		if err := rows.Scan(&chat.Id, &chat.Name, &imageUrl); err != nil {
			return nil, err
		}
		if imageUrl.Valid {
			chat.ImageUrl = imageUrl.String
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if chats == nil {
		chats = make([]*core.ChatDto, 0)
	}
	return chats, nil
}

func (r *SqlChatRepository) Create(ctx context.Context, chat *core.Chat) (string, error) {
	query, ok := r.queries["create"]
	if !ok {
		return "", core.ErrQueryNotConfigured
	}

	var imageUrl sql.NullString
	if chat.ImageUrl != "" {
		imageUrl = sql.NullString{String: chat.ImageUrl, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, chat.Name, imageUrl, chat.CreatedAt).Scan(&chat.Id)
	if err != nil {
		return "", err
	}
	return chat.Id, nil
}

func (r *SqlChatRepository) AddMember(ctx context.Context, chatId, userId string) error {
	query, ok := r.queries["add_member"]
	if !ok {
		return core.ErrQueryNotConfigured
	}

	_, err := r.db.ExecContext(ctx, query, chatId, userId, "member")
	return err
}

func (r *SqlChatRepository) RemoveMember(ctx context.Context, chatId, userId string) error {
	query, ok := r.queries["remove_member"]
	if !ok {
		return core.ErrQueryNotConfigured
	}

	_, err := r.db.ExecContext(ctx, query, chatId, userId)
	return err
}

func (r *SqlChatRepository) IsMemberInChat(ctx context.Context, chatId, userId string) (bool, error) {
	query, ok := r.queries["is_member_in_chat"]
	if !ok {
		return false, core.ErrQueryNotConfigured
	}

	var exists bool
	err := r.db.QueryRowContext(ctx, query, chatId, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *SqlChatRepository) CreateMessage(ctx context.Context, message *core.ChatMessage) error {
	query, ok := r.queries["create_message"]
	if !ok {
		return core.ErrQueryNotConfigured
	}

	var mediaUrl sql.NullString
	if message.MediaUrl != "" {
		mediaUrl = sql.NullString{String: message.MediaUrl, Valid: true}
	}

	var replyToId sql.NullString
	if message.ReplyToId != "" {
		replyToId = sql.NullString{String: message.ReplyToId, Valid: true}
	}

	err := r.db.QueryRowContext(ctx, query, message.Id, message.UserId, message.ChatId, message.Content, mediaUrl, replyToId, message.CreatedAt).Scan(&message.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *SqlChatRepository) GetMessages(ctx context.Context, chatId string, limit, offset int) ([]*core.ChatMessage, error) {
	query, ok := r.queries["get_messages"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}

	rows, err := r.db.QueryContext(ctx, query, chatId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*core.ChatMessage
	for rows.Next() {
		msg := &core.ChatMessage{}
		var mediaUrl sql.NullString
		var replyToId sql.NullString
		err := rows.Scan(&msg.Id, &msg.UserId, &msg.ChatId, &msg.Content, &mediaUrl, &replyToId, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		if mediaUrl.Valid {
			msg.MediaUrl = mediaUrl.String
		}
		if replyToId.Valid {
			msg.ReplyToId = replyToId.String
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *SqlChatRepository) GetMembers(ctx context.Context, chatId string) ([]string, error) {
	query, ok := r.queries["get_members"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}

	rows, err := r.db.QueryContext(ctx, query, chatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []string
	for rows.Next() {
		var userId string
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		members = append(members, userId)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}
