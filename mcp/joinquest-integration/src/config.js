export const DEFAULT_JOINQUEST_API_URL = 'https://joinquest.cc/graphql'

export function loadConfig() {
  const apiUrl = (process.env.JOINQUEST_API_URL || DEFAULT_JOINQUEST_API_URL).trim()
  const apiKey = (process.env.JOINQUEST_API_KEY || '').trim()
  const session = (process.env.JOINQUEST_SESSION || process.env.JOINQUEST_SESSION_COOKIE || '').trim()
  const cookieName = (process.env.JOINQUEST_SESSION_COOKIE_NAME || 'lobby_session').trim()

  if (apiKey) {
    return { apiUrl, authHeader: `Bearer ${apiKey}` }
  }

  if (session) {
    const cookieValue = session.includes('=') ? session : `${cookieName}=${session}`
    return { apiUrl, cookieValue }
  }

  throw new Error(
    'JOINQUEST_API_KEY is required. Generate one on your developer dashboard (Connect AI assistant), then set JOINQUEST_API_KEY in the MCP server env.',
  )
}
