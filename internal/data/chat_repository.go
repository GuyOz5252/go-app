package data

import (
	"context"
	"database/sql"

	"github.com/GuyOz5252/go-app/internal/core"
)

type SqlChatRepository struct {
	db      *sql.DB
	queries map[string]string
}

func NewSqlChatRepository(db *sql.DB, queries map[string]string) core.ChatRepository {
	return &SqlChatRepository{
		db:      db,
		queries: queries,
	}
}

func (r *SqlChatRepository) GetById(ctx context.Context, id string) (*core.Chat, error) {
	chat := &core.Chat{}
	query, ok := r.queries["get_chat_by_id"]
	if !ok {
		return nil, core.ErrQueryNotConfigured
	}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&chat.Id,
		&chat.Name,
		&chat.ImageUrl,
		&chat.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, core.ErrNotFound
		}
		return nil, err
	}
	return chat, nil
}

func (r *SqlChatRepository) ListByUserId(ctx context.Context, userId string) ([]*core.ChatDto, error) {
	panic("not implemented")
}

func (r *SqlChatRepository) Create(ctx context.Context, chat *core.Chat) (string, error) {
	panic("not implemented")
}

func (r *SqlChatRepository) AddMember(ctx context.Context, chatId, userId string) error {
	panic("not implemented")
}

func (r *SqlChatRepository) RemoveMember(ctx context.Context, chatId, userId string) error {
	panic("not implemented")
}

func (r *SqlChatRepository) IsMemberInChat(ctx context.Context, chatId, userId string) (bool, error) {
	panic("not implemented")
}

func (r *SqlChatRepository) CreateMessage(ctx context.Context, message *core.ChatMessage) error {
	panic("not implemented")
}

func (r *SqlChatRepository) GetMessages(ctx context.Context, chatId string, limit, offset int) ([]*core.ChatMessage, error) {
	panic("not implemented")
}
