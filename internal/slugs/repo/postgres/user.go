package postgres

import (
	"context"
	"errors"
	"fmt"

	"avitoTask/internal/slugs/service"
	types2 "avitoTask/internal/slugs/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDB struct {
	db *pgxpool.Pool
}

var _ service.UserRepo = UserDB{}

func NewUsersRepo(db *pgxpool.Pool) *UserDB {
	return &UserDB{db: db}
}

func (c UserDB) CreateUser(ctx context.Context, userID types2.UserID, segmentsToAdd []types2.Slug, segmentsToDelete []types2.Slug) error {
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

	if len(segmentsToAdd) != 0 {
		slugsToAdd := make([]string, len(segmentsToAdd))
		for i := range segmentsToAdd {
			slugsToAdd[i] = string(segmentsToAdd[i])
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
	}

	if len(segmentsToDelete) != 0 {
		slugsToDelete := make([]string, len(segmentsToDelete))
		for i := range segmentsToDelete {
			slugsToDelete[i] = string(segmentsToDelete[i])
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
	}

	return nil
}

func (c UserDB) segmentsExist(ctx context.Context, slugs []string) (bool, error) {
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

func (c UserDB) SelectActiveSegments(ctx context.Context, userID types2.UserID) ([]types2.Slug, error) {
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

	activeSegments := make([]types2.Slug, 0)
	for rows.Next() {
		var activeSegment types2.Slug
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