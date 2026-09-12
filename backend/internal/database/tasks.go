package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/scscloud/taskflow/internal/models"
)

// GetAllTasks returns every task ordered by creation date descending.
func (db *DB) GetAllTasks() ([]models.Task, error) {
	const q = `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC`

	rows, err := db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		t, err := scanTask(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	// Return an empty slice instead of nil so JSON encodes as [].
	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// GetTaskByID returns a single task or sql.ErrNoRows.
func (db *DB) GetTaskByID(id int64) (models.Task, error) {
	const q = `
		SELECT id, title, description, completed, created_at, updated_at
		FROM tasks WHERE id = ?`

	row := db.QueryRow(q, id)
	return scanTask(row)
}

// CreateTask inserts a new task and returns the full record.
func (db *DB) CreateTask(title, description string) (models.Task, error) {
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	const q = `
		INSERT INTO tasks (title, description, completed, created_at, updated_at)
		VALUES (?, ?, 0, ?, ?)
		RETURNING id, title, description, completed, created_at, updated_at`

	row := db.QueryRow(q, title, description, nowStr, nowStr)
	return scanTask(row)
}

// UpdateTaskCompleted changes the completed state and bumps updated_at.
func (db *DB) UpdateTaskCompleted(id int64, completed bool) (models.Task, error) {
	completedInt := 0
	if completed {
		completedInt = 1
	}
	now := time.Now().UTC().Format(time.RFC3339)

	const q = `
		UPDATE tasks SET completed = ?, updated_at = ?
		WHERE id = ?
		RETURNING id, title, description, completed, created_at, updated_at`

	row := db.QueryRow(q, completedInt, now, id)
	return scanTask(row)
}

// DeleteTask removes a task. Returns sql.ErrNoRows if not found.
func (db *DB) DeleteTask(id int64) error {
	res, err := db.Exec(`DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// scanner is satisfied by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanTask(s scanner) (models.Task, error) {
	var t models.Task
	var completedInt int
	var createdStr, updatedStr string

	err := s.Scan(&t.ID, &t.Title, &t.Description, &completedInt, &createdStr, &updatedStr)
	if err != nil {
		return t, err
	}

	t.Completed = completedInt == 1

	t.CreatedAt, err = time.Parse(time.RFC3339, createdStr)
	if err != nil {
		return t, fmt.Errorf("parse created_at: %w", err)
	}
	t.UpdatedAt, err = time.Parse(time.RFC3339, updatedStr)
	if err != nil {
		return t, fmt.Errorf("parse updated_at: %w", err)
	}

	return t, nil
}

