import { request } from './client.js'

/** Check backend connectivity. */
export function checkHealth() {
  return request('/api/health')
}

/** Fetch all tasks. */
export function fetchTasks() {
  return request('/api/tasks')
}

/** Create a new task. */
export function createTask({ title, description }) {
  return request('/api/tasks', {
    method: 'POST',
    body: JSON.stringify({ title, description }),
  })
}

/** Toggle a task's completed state. */
export function updateTask(id, { completed }) {
  return request(`/api/tasks/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ completed }),
  })
}

/** Permanently delete a task. */
export function deleteTask(id) {
  return request(`/api/tasks/${id}`, { method: 'DELETE' })
}

