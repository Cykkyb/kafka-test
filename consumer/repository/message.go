package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{
		db: db,
	}
}

func (a *MessagePostgres) SaveMessage(ctx context.Context, massage string) error {
	query := fmt.Sprintf(`INSERT INTO %s (message_data) VALUES ($1) RETURNING id`, massageTable)
	_, err := a.db.ExecContext(ctx, query, massage)
	return err
}
