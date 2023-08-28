package postgres

import (
	"avitoTask/internal/model"
	"fmt"
)

func (c Client) CreateUser(userID model.UserID, segmentsToAdd []model.Slug, segmentsToDelete []model.Slug) error {
	q := `
		SELECT
		  * 
		FROM 
		  users_segments
		WHERE
		  user_id = $1
	`
	exists, err := c.rowExists(q, userID)
	if err != nil {
		return fmt.Errorf("failed to check if row exists: %w", err)
	}

	if len(segmentsToAdd) != 0 {
		if !exists {
			q := fmt.Sprintf("INSERT INTO users_segments(user_id, %s) VALUES($1, %s)", getSegmentsString(segmentsToAdd), getOnesString(segmentsToAdd))
			if _, err := c.db.Exec(q, userID); err != nil {
				return fmt.Errorf("failed to insert new user segments: %w", err)
			}
		}

		q = fmt.Sprintf("UPDATE users_segments SET %s WHERE user_id = $1", getAddingString(segmentsToAdd))
		if _, err := c.db.Exec(q, userID); err != nil {
			return fmt.Errorf("failed to insert user segments: %w", err)
		}
	}

	if len(segmentsToDelete) != 0 {
		q = fmt.Sprintf("UPDATE users_segments SET %s WHERE user_id = $1", getDeletingString(segmentsToDelete))
		if _, err := c.db.Exec(q, userID); err != nil {
			return fmt.Errorf("failed to insert user segments: %w", err)
		}
	}

	return nil
}

func (c Client) rowExists(query string, args ...interface{}) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	if err := c.db.QueryRow(query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to scan query row: %w", err)
	}
	return exists, nil
}

func getSegmentsString(segments []model.Slug) string {
	out := ""
	for _, segment := range segments {
		out += fmt.Sprintf("%s, ", segment)
	}
	return out[:len(out)-2]
}

func getOnesString(segments []model.Slug) string {
	out := ""
	for range segments {
		out += fmt.Sprintf("%d, ", 1)
	}
	return out[:len(out)-2]
}

func getAddingString(segments []model.Slug) string {
	out := ""
	for _, segment := range segments {
		out += fmt.Sprintf("%s = %d, ", segment, 1)
	}
	return out[:len(out)-2]
}

func getDeletingString(segments []model.Slug) string {
	out := ""
	for _, segment := range segments {
		out += fmt.Sprintf("%s = NULL, ", segment)
	}
	return out[:len(out)-2]
}

func (c Client) SelectActiveSegments(userID model.UserID) ([]model.Slug, error) {
	q := `
		SELECT
		  column_name
		FROM
		  information_schema.columns
		WHERE
		  table_name = 'users_segments'
	`
	rows, err := c.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slugs := make([]model.Slug, 0)
	for rows.Next() {
		var slug model.Slug
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("failed to parse slug: %w", err)
		}
		slugs = append(slugs, slug)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse slugs: %w", err)
	}
	slugs = slugs[1:]

	var activeSegments []model.Slug

	for i, slug := range slugs {
		q = fmt.Sprintf("SELECT %s FROM users_segments WHERE user_id = $1", slug)

		var slugValue int
		if err := c.db.QueryRow(q, userID).Scan(&slugValue); err != nil {
			return nil, fmt.Errorf("failed to parse slug values: %w", err)
		}

		if slugValue == 1 {
			activeSegments = append(activeSegments, slugs[i])
		}
	}

	return activeSegments, nil
}

func getColumnsString(slugs []model.Slug) string {
	out := ""
	for _, slug := range slugs {
		out += fmt.Sprintf("%s, ", slug)
	}
	return out[:len(out)-2]
}
