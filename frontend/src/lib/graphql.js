import { getGraphQLUrl } from './env'

export function isTransientServerError(message) {
  const text = (message || '').toLowerCase()
  return (
    text.includes('load failed')
    || text.includes('failed to fetch')
    || text.includes('networkerror')
    || text.includes('network request failed')
    || /api request failed \(502\)|api request failed \(503\)|api request failed \(504\)/.test(text)
  )
}

export function isAuthRequiredError(message) {
  return /authentication required/i.test(message || '')
}

export async function graphqlRequest(query, variables = {}) {
  const response = await fetch(getGraphQLUrl(), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ query, variables }),
  })

  if (!response.ok) {
    let detail = ''
    try {
      const payload = await response.clone().json()
      detail = payload.errors?.[0]?.message || payload.error || ''
    } catch {
      // ignore parse errors
    }
    throw new Error(detail || `API request failed (${response.status})`)
  }

  const payload = await response.json()
  if (payload.errors?.length) {
    throw new Error(payload.errors[0]?.message || 'GraphQL request failed')
  }

  return payload.data
}
