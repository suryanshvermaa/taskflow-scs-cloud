/**
 * Centralized API base URL, read once from the Vite environment.
 *
 * During a Vite build (npm run build), Vite replaces `import.meta.env.VITE_*`
 * with the literal values found in the .env file.  The resulting JavaScript
 * bundle therefore contains the baked-in URL — which is why this value must
 * never be a secret.
 *
 * To target a different backend, update frontend/.env and rebuild.
 */
const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

/**
 * Thin fetch wrapper that:
 *  - Prepends BASE_URL to relative paths
 *  - Sets Content-Type: application/json for requests with a body
 *  - Parses JSON responses
 *  - Throws a structured error for non-2xx responses
 */
export async function request(path, options = {}) {
  const url = `${BASE_URL}${path}`

  const headers = {
    'Content-Type': 'application/json',
    ...options.headers,
  }

  const response = await fetch(url, { ...options, headers })

  // Parse JSON body regardless of status so we can read error messages.
  let data = null
  const contentType = response.headers.get('content-type') ?? ''
  if (contentType.includes('application/json')) {
    data = await response.json()
  }

  if (!response.ok) {
    const message = data?.error ?? `HTTP ${response.status}`
    const err = new Error(message)
    err.status = response.status
    err.data = data
    throw err
  }

  return data
}

