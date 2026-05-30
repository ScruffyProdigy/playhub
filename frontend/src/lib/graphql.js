import { getGraphQLUrl } from './env'

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
