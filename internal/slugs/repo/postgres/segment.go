package postgres

import (
	"context"
	"errors"
	"fmt"

	"avitoTask/internal/slugs/service"
	"avitoTask/internal/slugs/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SegmentDB struct {
	db *pgxpool.Pool
}

var _ service.SegmentRepo = SegmentDB{}

func NewSegmentsRepo(db *pgxpool.Pool) *SegmentDB {
	return &SegmentDB{db: db}
}

func (c SegmentDB) CreateSegment(ctx context.Context, slug types.Slug) error {
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

func (c SegmentDB) DeleteSegment(ctx context.Context, slug types.Slug) (err error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit(ctx)
			return
		}

		err = errors.Join(err, tx.Rollback(ctx))
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
