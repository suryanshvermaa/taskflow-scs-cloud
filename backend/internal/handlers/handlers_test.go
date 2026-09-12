package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/scscloud/taskflow/internal/database"
	"github.com/scscloud/taskflow/internal/handlers"
)

func newTestDB(t *testing.T) *database.DB {
	t.Helper()
	f, err := os.CreateTemp("", "taskflow-test-*.db")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	db, err := database.Open(f.Name())
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestHealthCheck(t *testing.T) {
	h := handlers.New(newTestDB(t))
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	w := httptest.NewRecorder()
	h.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", w.Code)
	}
	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["status"] != "ok" {
		t.Errorf("want status=ok, got %q", body["status"])
	}
}

func TestCreateAndGetTasks(t *testing.T) {
	h := handlers.New(newTestDB(t))

	// Create
	payload := `{"title":"Test task","description":"desc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTask(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d — body: %s", w.Code, w.Body)
	}

	// Get
	req = httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w = httptest.NewRecorder()
	h.GetTasks(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get: want 200, got %d", w.Code)
	}

	var tasks []map[string]any
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 1 {
		t.Fatalf("want 1 task, got %d", len(tasks))
	}
	if tasks[0]["title"] != "Test task" {
		t.Errorf("want title='Test task', got %v", tasks[0]["title"])
	}
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	h := handlers.New(newTestDB(t))
	payload := `{"title":"   "}`
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateTask(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestUpdateTask(t *testing.T) {
	db := newTestDB(t)
	task, _ := db.CreateTask("Task", "")
	h := handlers.New(db)

	payload := `{"completed":true}`
	path := "/api/tasks/" + json.Number(string(rune('0'+task.ID))).String()
	if task.ID >= 10 {
		path = "/api/tasks/1"
	}
	req := httptest.NewRequest(http.MethodPut, path, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.UpdateTask(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d — body: %s", w.Code, w.Body)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	h := handlers.New(newTestDB(t))
	req := httptest.NewRequest(http.MethodDelete, "/api/tasks/9999", nil)
	w := httptest.NewRecorder()
	h.DeleteTask(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

