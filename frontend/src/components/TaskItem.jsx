import { useState } from 'react'
import styles from './TaskItem.module.css'

/** Format an ISO date string to something readable. */
function formatDate(iso) {
  const d = new Date(iso)
  return d.toLocaleDateString('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

export function TaskItem({ task, onToggle, onDelete }) {
  const [toggling, setToggling] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [error, setError] = useState(null)

  async function handleToggle() {
    setToggling(true)
    setError(null)
    try {
      await onToggle(task.id, !task.completed)
    } catch (err) {
      setError(err.message || 'Update failed. Try again.')
    } finally {
      setToggling(false)
    }
  }

  async function handleDelete() {
    if (!window.confirm(`Delete "${task.title}"?`)) return
    setDeleting(true)
    setError(null)
    try {
      await onDelete(task.id)
    } catch (err) {
      setDeleting(false)
      setError(err.message || 'Delete failed. Try again.')
    }
  }

  return (
    <article
      className={`${styles.item} ${task.completed ? styles.completed : ''}`}
      aria-label={`Task: ${task.title}${task.completed ? ' (completed)' : ''}`}
    >
      <div className={styles.top}>
        {/* Completion indicator */}
        <button
          className={styles.checkBtn}
          onClick={handleToggle}
          disabled={toggling || deleting}
          aria-label={task.completed ? 'Mark as active' : 'Mark as complete'}
          title={task.completed ? 'Mark as active' : 'Mark as complete'}
        >
          <span className={styles.checkCircle} aria-hidden="true">
            {task.completed && (
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path
                  d="M2 6.5l2.8 2.8 5.2-5.6"
                  stroke="currentColor"
                  strokeWidth="1.75"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            )}
          </span>
        </button>

        <div className={styles.content}>
          <h3 className={styles.title}>{task.title}</h3>
          {task.description && (
            <p className={styles.description}>{task.description}</p>
          )}
        </div>
      </div>

      {error && (
        <p className={styles.errorMsg} role="alert">
          {error}
        </p>
      )}

      <div className={styles.footer}>
        <span className={styles.date}>
          <time dateTime={task.created_at}>{formatDate(task.created_at)}</time>
        </span>
        <div className={styles.actions}>
          <button
            className={`${styles.actionBtn} ${styles.toggleBtn}`}
            onClick={handleToggle}
            disabled={toggling || deleting}
          >
            {toggling ? '…' : task.completed ? 'Undo' : 'Complete'}
          </button>
          <button
            className={`${styles.actionBtn} ${styles.deleteBtn}`}
            onClick={handleDelete}
            disabled={deleting || toggling}
            aria-label={`Delete task: ${task.title}`}
          >
            {deleting ? '…' : 'Delete'}
          </button>
        </div>
      </div>
    </article>
  )
}

