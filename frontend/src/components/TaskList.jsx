import { TaskItem } from './TaskItem.jsx'
import styles from './TaskList.module.css'

export function TaskList({ tasks, filter, loading, error, onReload, onToggle, onDelete }) {
  if (loading) {
    return (
      <div className={styles.state}>
        <div className={styles.spinner} aria-hidden="true" />
        <p className={styles.stateText}>Loading tasks…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className={styles.state}>
        <div className={styles.errorIcon} aria-hidden="true">!</div>
        <p className={styles.stateHeading}>Unable to connect to the backend</p>
        <p className={styles.stateText}>{error}</p>
        <button className={styles.retryBtn} onClick={onReload}>
          Retry
        </button>
      </div>
    )
  }

  const filtered = tasks.filter(t => {
    if (filter === 'active') return !t.completed
    if (filter === 'completed') return t.completed
    return true
  })

  if (filtered.length === 0) {
    const emptyMessages = {
      all: { heading: 'No tasks yet', body: 'Create your first task above to get started.' },
      active: { heading: 'No active tasks', body: 'All your tasks are completed — great work!' },
      completed: { heading: 'No completed tasks', body: 'Complete a task to see it here.' },
    }
    const { heading, body } = emptyMessages[filter]

    return (
      <div className={styles.state}>
        <div className={styles.emptyIcon} aria-hidden="true">
          <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
            <rect x="6" y="6" width="28" height="28" rx="4" stroke="var(--color-border-strong)" strokeWidth="1.5"/>
            <path d="M13 20h14M13 26h8" stroke="var(--color-border-strong)" strokeWidth="1.5" strokeLinecap="round"/>
          </svg>
        </div>
        <p className={styles.stateHeading}>{heading}</p>
        <p className={styles.stateText}>{body}</p>
      </div>
    )
  }

  return (
    <ul className={styles.list} role="list">
      {filtered.map(task => (
        <li key={task.id}>
          <TaskItem task={task} onToggle={onToggle} onDelete={onDelete} />
        </li>
      ))}
    </ul>
  )
}

