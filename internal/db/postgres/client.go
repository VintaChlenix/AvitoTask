package postgres

import (
	"avitoTask/internal/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	db *pgxpool.Pool
}

var _ db.DB = Client{}

func NewClient(connectionString string) (*Client, error) {
	conn, err := pgxpool.New(context.TODO(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &Client{db: conn}, nil
}
