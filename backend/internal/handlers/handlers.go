package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/scscloud/taskflow/internal/database"
	"github.com/scscloud/taskflow/internal/models"
)

// Handler holds shared dependencies for all HTTP handlers.
type Handler struct {
	db *database.DB
}

// New creates a Handler with the given database.
func New(db *database.DB) *Handler {
	return &Handler{db: db}
}

// HealthCheck handles GET /api/health.
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "taskflow-api",
	})
}

// GetTasks handles GET /api/tasks.
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.db.GetAllTasks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve tasks")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

// CreateTask handles POST /api/tasks.
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, "Task title is required")
		return
	}
	if len(title) > 255 {
		writeError(w, http.StatusBadRequest, "Task title must not exceed 255 characters")
		return
	}

	task, err := h.db.CreateTask(title, strings.TrimSpace(req.Description))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

// UpdateTask handles PUT /api/tasks/{id}.
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	if req.Completed == nil {
		writeError(w, http.StatusBadRequest, "Field 'completed' is required")
		return
	}

	task, err := h.db.UpdateTaskCompleted(id, *req.Completed)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

// DeleteTask handles DELETE /api/tasks/{id}.
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := taskIDFromPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	err = h.db.DeleteTask(id)
	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, fmt.Sprintf("Task %d not found", id))
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete task")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Task deleted"})
}

// ─── Helpers ────────────────────────────────────────────────────────────────

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// taskIDFromPath extracts the numeric ID from a path like /api/tasks/42.
func taskIDFromPath(path string) (int64, error) {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("no id segment")
	}
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}

