package models

import "time"

// Task represents a single task in the TaskFlow application.
type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateTaskRequest is the expected body for POST /api/tasks.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTaskRequest is the expected body for PUT /api/tasks/:id.
type UpdateTaskRequest struct {
	Completed *bool  `json:"completed"`
	Title     string `json:"title"`
}

