package postgres

import (
	"avitoTask/internal/model"
	"fmt"
)

func (c Client) CreateUser(userID model.UserID, segmentsToAdd []model.Slug, segmentsToDelete []model.Slug) error {
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

	exist, err := c.segmentsExist(segmentsToAdd)
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
	for _, segment := range segmentsToAdd {
		if _, err := c.db.Exec(q, userID, segment); err != nil {
			return fmt.Errorf("failed to insert user segments: %w", err)
		}
	}

	exist, err = c.segmentsExist(segmentsToDelete)
	if err != nil {
		return fmt.Errorf("failed to check if segments exist to delete are: %w", err)
	}
	if !exist {
		return fmt.Errorf("trying to delete non existing segments: %w", err)
	}

	q = `
		DELETE FROM
		  users_segments
		WHERE
		  user_id = $1 AND slug = $2
	`
	for _, segment := range segmentsToDelete {
		if _, err := c.db.Exec(q, userID, segment); err != nil {
			return fmt.Errorf("failed to insert user segments: %w", err)
		}
	}

	return nil
}

func (c Client) segmentsExist(segments []model.Slug) (bool, error) {
	slugs := make([]string, len(segments))
	for i := range segments {
		slugs[i] = string(segments[i])
	}
	q := `
		SELECT
		  *
		FROM
		  segments
		WHERE
		  slug = any($1)
	`

	rows, err := c.db.Query(q, slugs)
	if err != nil {
		return false, fmt.Errorf("failed to select existing segments: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		count++
	}

	if count != len(segments) {
		return false, nil
	}

	return true, nil
}

func (c Client) SelectActiveSegments(userID model.UserID) ([]model.Slug, error) {
	q := `
		SELECT
		  slug
		FROM
		  users_segments
		WHERE
		  user_id = $1
	`
	rows, err := c.db.Query(q, userID)
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
