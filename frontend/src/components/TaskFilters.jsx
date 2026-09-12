import styles from './TaskFilters.module.css'

const FILTERS = ['all', 'active', 'completed']

export function TaskFilters({ filter, tasks, onFilterChange }) {
  const counts = {
    all: tasks.length,
    active: tasks.filter(t => !t.completed).length,
    completed: tasks.filter(t => t.completed).length,
  }

  return (
    <nav className={styles.nav} aria-label="Task filters">
      <ul className={styles.list} role="list">
        {FILTERS.map(f => (
          <li key={f}>
            <button
              className={`${styles.btn} ${filter === f ? styles.active : ''}`}
              onClick={() => onFilterChange(f)}
              aria-pressed={filter === f}
            >
              {f.charAt(0).toUpperCase() + f.slice(1)}
              <span className={styles.count}>{counts[f]}</span>
            </button>
          </li>
        ))}
      </ul>
    </nav>
  )
}

