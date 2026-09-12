import { useState, useEffect } from 'react'
import { checkHealth } from '../api/tasks.js'

/**
 * Polls the backend health endpoint and returns the current connection status.
 * status: 'connecting' | 'connected' | 'offline'
 */
export function useHealth(intervalMs = 30_000) {
  const [status, setStatus] = useState('connecting')

  useEffect(() => {
    let cancelled = false

    async function ping() {
      try {
        await checkHealth()
        if (!cancelled) setStatus('connected')
      } catch {
        if (!cancelled) setStatus('offline')
      }
    }

    ping()
    const id = setInterval(ping, intervalMs)
    return () => {
      cancelled = true
      clearInterval(id)
    }
  }, [intervalMs])

  return status
}

