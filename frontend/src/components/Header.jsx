import { useHealth } from '../hooks/useHealth.js'
import styles from './Header.module.css'

const STATUS_LABELS = {
  connecting: 'Connecting…',
  connected: 'Connected',
  offline: 'Offline',
}

export function Header() {
  const status = useHealth()

  return (
    <header className={styles.header}>
      <div className={styles.brand}>
        <svg
          className={styles.logo}
          width="28"
          height="28"
          viewBox="0 0 32 32"
          fill="none"
          aria-hidden="true"
          xmlns="http://www.w3.org/2000/svg"
        >
          <rect width="32" height="32" rx="6" fill="var(--color-primary)" />
          <path
            d="M8 16.5l5.5 5.5 10.5-11"
            stroke="#fff"
            strokeWidth="2.5"
            strokeLinecap="round"
            strokeLinejoin="round"
          />
        </svg>
        <div>
          <h1 className={styles.title}>TaskFlow</h1>
          <p className={styles.subtitle}>Mini task manager</p>
        </div>
      </div>

      <div
        className={`${styles.statusBadge} ${styles[status]}`}
        aria-label={`Backend status: ${STATUS_LABELS[status]}`}
        role="status"
      >
        <span className={styles.dot} aria-hidden="true" />
        <span className={styles.statusText}>
          Backend&nbsp;·&nbsp;{STATUS_LABELS[status]}
        </span>
      </div>
    </header>
  )
}

