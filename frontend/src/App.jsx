import { useState } from 'react'
import { Header } from './components/Header.jsx'
import { TaskForm } from './components/TaskForm.jsx'
import { TaskFilters } from './components/TaskFilters.jsx'
import { TaskList } from './components/TaskList.jsx'
import { useTasks } from './hooks/useTasks.js'
import styles from './App.module.css'

export default function App() {
  const [filter, setFilter] = useState('all')
  const { tasks, loading, error, reload, addTask, toggleTask, removeTask } = useTasks()

  return (
    <div className={styles.root}>
      <Header />

      <main className={styles.main}>
        <div className={styles.card}>
          <TaskForm onAdd={addTask} />
          <TaskFilters filter={filter} tasks={tasks} onFilterChange={setFilter} />
          <TaskList
            tasks={tasks}
            filter={filter}
            loading={loading}
            error={error}
            onReload={reload}
            onToggle={toggleTask}
            onDelete={removeTask}
          />
        </div>
      </main>

      <footer className={styles.footer}>
        <p>
          TaskFlow · SCS Cloud demo ·{' '}
          <a
            href="https://github.com/scscloud/taskflow"
            target="_blank"
            rel="noopener noreferrer"
            className={styles.footerLink}
          >
            GitHub
          </a>
        </p>
      </footer>
    </div>
  )
}

