package postgres

import (
	"context"
	"fmt"

	"avitoTask/internal/service"
	"avitoTask/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentsRepo struct {
	db *pgxpool.Pool
}

var _ service.SegmentsRepo = SegmentsRepo{}

func NewSegmentsRepo(db *pgxpool.Pool) *SegmentsRepo {
	return &SegmentsRepo{db: db}
}

func (c SegmentsRepo) CreateSegment(ctx context.Context, slug types.Slug) error {
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

func (c SegmentsRepo) DeleteSegment(ctx context.Context, slug types.Slug) error {
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
