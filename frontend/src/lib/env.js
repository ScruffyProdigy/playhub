export function getApiBaseUrl() {
  const configured = window.env?.REACT_APP_API_BASE_URL?.trim()
  if (configured) {
    return configured.replace(/\/$/, '')
  }
  return ''
}

export function getGraphQLUrl() {
  const base = getApiBaseUrl()
  if (base) {
    return `${base}/graphql`
  }
  return '/graphql'
}

/**
 * WebSocket endpoint for GraphQL subscriptions.
 * Connects to the backend directly (not through the Vite proxy) so the connection stays
 * stable and auth is sent via connection_init Authorization from subscriptionAuth.
 */
export function getGraphQLWsUrl() {
  const base = getApiBaseUrl()
  const path = '/graphql'
  if (base) {
    const wsBase = base.replace(/^http:\/\//, 'ws://').replace(/^https:\/\//, 'wss://')
    return `${wsBase}${path}`
  }
  return `ws://127.0.0.1:8080${path}`
}
