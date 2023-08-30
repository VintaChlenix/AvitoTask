package postgres

import (
	"avitoTask/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (c Client) CreateUser(ctx context.Context, userID model.UserID, segmentsToAdd []model.Slug, segmentsToDelete []model.Slug) error {
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

func (c Client) segmentsExist(ctx context.Context, slugs []string) (bool, error) {
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

func (c Client) SelectActiveSegments(ctx context.Context, userID model.UserID) ([]model.Slug, error) {
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

	activeSegments := make([]model.Slug, 0)
	for rows.Next() {
		var activeSegment model.Slug
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
