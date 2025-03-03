package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	massageTable = "messages"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBname   string
	SSL      string
}

type MessageRepository interface {
	SaveMessage(ctx context.Context, message string) error
}

type Repository struct {
	MessageRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		MessageRepository: NewMessagePostgres(db),
	}
}

func ConnectDb(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBname, cfg.SSL))

	if err != nil {
		return nil, err
	}

	return db, nil
}
