package postgres

import (
	"avitoTask/internal/model"
	"fmt"
)

func (c Client) CreateSegment(slug model.Slug) error {
	q := fmt.Sprintf("ALTER TABLE users_segments\n ADD COLUMN IF NOT EXISTS user_id INTEGER PRIMARY KEY,\n ADD COLUMN IF NOT EXISTS %s INTEGER NOT NULL DEFAULT 0;", slug)
	if _, err := c.db.Exec(q); err != nil {
		return fmt.Errorf("failed to create segment column: %w", err)
	}
	return nil
}

func (c Client) DeleteSegment(slug model.Slug) error {
	q := fmt.Sprintf("ALTER TABLE users_segments\n DROP COLUMN %s", slug)
	if _, err := c.db.Exec(q); err != nil {
		return fmt.Errorf("failed to delete segment column: %w", err)
	}
	return nil
}
