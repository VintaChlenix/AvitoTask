package postgres

import (
	"avitoTask/internal/db"
	"fmt"
	"github.com/jackc/pgx"
)

type Client struct {
	db *pgx.Conn
}

var _ db.DB = Client{}

func NewClient(connectionString string) (*Client, error) {
	connConfig, err := pgx.ParseConnectionString(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection string: %w", err)
	}
	conn, err := pgx.Connect(connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}
	return &Client{db: conn}, nil
}
