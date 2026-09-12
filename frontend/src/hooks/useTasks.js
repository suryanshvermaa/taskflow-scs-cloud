import { useState, useEffect, useCallback } from 'react'
import { fetchTasks, createTask, updateTask, deleteTask } from '../api/tasks.js'

/**
 * useTasks manages the full lifecycle of task data:
 *  - loading from the backend
 *  - creating, updating, and deleting tasks
 *  - surfacing loading / error state to the UI
 */
export function useTasks() {
  const [tasks, setTasks] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  const load = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const data = await fetchTasks()
      setTasks(data)
    } catch (err) {
      setError(err.message || 'Failed to load tasks')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  const addTask = useCallback(async ({ title, description }) => {
    const task = await createTask({ title, description })
    setTasks(prev => [task, ...prev])
    return task
  }, [])

  const toggleTask = useCallback(async (id, completed) => {
    const updated = await updateTask(id, { completed })
    setTasks(prev => prev.map(t => (t.id === id ? updated : t)))
    return updated
  }, [])

  const removeTask = useCallback(async (id) => {
    await deleteTask(id)
    setTasks(prev => prev.filter(t => t.id !== id))
  }, [])

  return { tasks, loading, error, reload: load, addTask, toggleTask, removeTask }
}

