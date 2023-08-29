package postgres

import (
	"avitoTask/internal/model"
	"fmt"
)

func (c Client) CreateSegment(slug model.Slug) error {
	q := `
		INSERT INTO
		  segments(slug)
		VALUES
		  ($1)
	`
	if _, err := c.db.Exec(q, slug); err != nil {
		return fmt.Errorf("failed to create segment column: %w", err)
	}
	return nil
}

func (c Client) DeleteSegment(slug model.Slug) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	q := `
		DELETE FROM
		  segments
		WHERE
		  slug = $1
	`
	if _, err := c.db.Exec(q, slug); err != nil {
		return fmt.Errorf("failed to delete segment column: %w", err)
	}
	q = `
		DELETE FROM
		  users_segments
		WHERE
		  slug = $1
	`
	if _, err := c.db.Exec(q, slug); err != nil {
		return fmt.Errorf("failed to delete segment column: %w", err)
	}

	return nil
}
