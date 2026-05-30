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
