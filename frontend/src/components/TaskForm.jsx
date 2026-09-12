import { useState } from 'react'
import styles from './TaskForm.module.css'

export function TaskForm({ onAdd }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState(null)

  async function handleSubmit(e) {
    e.preventDefault()
    const trimmedTitle = title.trim()

    if (!trimmedTitle) {
      setError('Title is required.')
      return
    }

    setSubmitting(true)
    setError(null)

    try {
      await onAdd({ title: trimmedTitle, description: description.trim() })
      setTitle('')
      setDescription('')
    } catch (err) {
      setError(err.message || 'Failed to create task. Please try again.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className={styles.section} aria-labelledby="form-heading">
      <h2 id="form-heading" className={styles.heading}>
        Add a task
      </h2>

      <form className={styles.form} onSubmit={handleSubmit} noValidate>
        <div className={styles.field}>
          <label htmlFor="task-title" className={styles.label}>
            Title <span className={styles.required} aria-hidden="true">*</span>
          </label>
          <input
            id="task-title"
            type="text"
            className={styles.input}
            placeholder="What needs to be done?"
            value={title}
            onChange={e => { setTitle(e.target.value); setError(null) }}
            disabled={submitting}
            maxLength={255}
            aria-required="true"
            aria-invalid={error ? 'true' : undefined}
            aria-describedby={error ? 'form-error' : undefined}
          />
        </div>

        <div className={styles.field}>
          <label htmlFor="task-desc" className={styles.label}>
            Description <span className={styles.optional}>(optional)</span>
          </label>
          <textarea
            id="task-desc"
            className={`${styles.input} ${styles.textarea}`}
            placeholder="Add more detail…"
            value={description}
            onChange={e => setDescription(e.target.value)}
            disabled={submitting}
            rows={3}
          />
        </div>

        {error && (
          <p id="form-error" className={styles.errorMsg} role="alert">
            {error}
          </p>
        )}

        <div className={styles.actions}>
          <button
            type="submit"
            className={styles.submitBtn}
            disabled={submitting}
          >
            {submitting ? 'Adding…' : 'Add Task'}
          </button>
        </div>
      </form>
    </section>
  )
}

