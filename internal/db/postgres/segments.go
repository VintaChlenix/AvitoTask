package postgres

import (
	"avitoTask/internal/model"
	"context"
	"fmt"
)

func (c Client) CreateSegment(ctx context.Context, slug model.Slug) error {
	q := `
		INSERT INTO
		  segments(slug)
		VALUES
		  ($1)
	`
	if _, err := c.db.Exec(ctx, q, slug); err != nil {
		return fmt.Errorf("failed to create segment column: %w", err)
	}
	return nil
}

func (c Client) DeleteSegment(ctx context.Context, slug model.Slug) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	q := `
		DELETE FROM
		  segments
		WHERE
		  slug = $1
	`
	if _, err := c.db.Exec(ctx, q, slug); err != nil {
		return fmt.Errorf("failed to delete segment column: %w", err)
	}
	q = `
		DELETE FROM
		  users_segments
		WHERE
		  slug = $1
	`
	if _, err := c.db.Exec(ctx, q, slug); err != nil {
		return fmt.Errorf("failed to delete segment column: %w", err)
	}

	return nil
}
