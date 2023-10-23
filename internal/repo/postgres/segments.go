package postgres

import (
	"avitoTask/internal/service"
	"avitoTask/internal/types"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentsClient struct {
	db *pgxpool.Pool
}

var _ service.SegmentsRepo = SegmentsClient{}

func NewSegmentClient(db *pgxpool.Pool) *SegmentsClient {
	return &SegmentsClient{db: db}
}

func (c SegmentsClient) Close() {
	c.db.Close()
}

func (c SegmentsClient) CreateSegment(ctx context.Context, slug types.Slug) error {
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

func (c SegmentsClient) DeleteSegment(ctx context.Context, slug types.Slug) error {
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
