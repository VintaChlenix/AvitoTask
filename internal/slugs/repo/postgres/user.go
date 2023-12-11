package postgres

import (
	"context"
	"errors"
	"fmt"

	"avitoTask/internal/slugs/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *User {
	return &User{db: db}
}

func (c User) CreateUser(ctx context.Context, userID types.UserID, segmentsToAdd []types.Slug, segmentsToDelete []types.Slug) (err error) {
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

	if err := c.addSegments(ctx, userID, segmentsToAdd); err != nil {
		return fmt.Errorf("failed to add segments: %w", err)
	}

	if err := c.deleteSegments(ctx, userID, segmentsToDelete); err != nil {
		return fmt.Errorf("failed to delete segments: %w", err)
	}

	return nil
}

func (c User) addSegments(ctx context.Context, userID types.UserID, segments []types.Slug) error {
	if len(segments) == 0 {
		return nil
	}
	slugsToAdd := make([]string, len(segments))
	for i := range segments {
		slugsToAdd[i] = string(segments[i])
	}
	exist, err := c.segmentsExist(ctx, slugsToAdd)
	if err != nil {
		return fmt.Errorf("failed to check if segments to add are exist: %w", err)
	}
	if !exist {
		return fmt.Errorf("trying to insert non existing segments: %w", err)
	}
	q := `
INSERT INTO
  users_segments(user_id, slug)
VALUES
  ($1, $2)
ON CONFLICT DO NOTHING
	`
	batch := &pgx.Batch{}
	for _, slugToAdd := range slugsToAdd {
		batch.Queue(q, userID, slugToAdd)
	}
	br := c.db.SendBatch(ctx, batch)
	ct, err := br.Exec()
	if err != nil {
		return fmt.Errorf("failed to insert user segments: %w", err)
	}
	defer br.Close()
	ct.RowsAffected()

	return nil
}

func (c User) deleteSegments(ctx context.Context, userID types.UserID, segments []types.Slug) error {
	if len(segments) == 0 {
		return nil
	}
	slugsToDelete := make([]string, len(segments))
	for i := range segments {
		slugsToDelete[i] = string(segments[i])
	}
	exist, err := c.segmentsExist(ctx, slugsToDelete)
	if err != nil {
		return fmt.Errorf("failed to check if segments exist to delete are: %w", err)
	}
	if !exist {
		return fmt.Errorf("trying to delete non existing segments: %w", err)
	}
	q := `
DELETE FROM
  users_segments
WHERE
  user_id = $1 AND slug = any($2)
	`
	if _, err := c.db.Exec(ctx, q, userID, slugsToDelete); err != nil {
		return fmt.Errorf("failed to insert user segments: %w", err)
	}

	return nil
}

func (c User) segmentsExist(ctx context.Context, slugs []string) (bool, error) {
	q := `
SELECT
  *
FROM
  segments
WHERE
  slug = any($1)
	`

	rows, err := c.db.Query(ctx, q, slugs)
	if err != nil {
		return false, fmt.Errorf("failed to select existing segments: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}
	if count != len(slugs) {
		return false, nil
	}

	return true, nil
}

func (c User) SelectActiveSegments(ctx context.Context, userID types.UserID) ([]types.Slug, error) {
	q := `
SELECT
  slug
FROM
  users_segments
WHERE
  user_id = $1
	`
	rows, err := c.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to select user active segments: %w", err)
	}
	defer rows.Close()

	activeSegments := make([]types.Slug, 0)
	for rows.Next() {
		var activeSegment types.Slug
		if err := rows.Scan(&activeSegment); err != nil {
			return nil, fmt.Errorf("failed to parse slug: %w", err)
		}
		activeSegments = append(activeSegments, activeSegment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse slugs: %w", err)
	}

	return activeSegments, nil
}
