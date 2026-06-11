/** @returns {boolean} */
export function isLobbyDebugEnabled() {
  if (typeof window === 'undefined') {
    return false
  }
  if (window.env?.REACT_APP_LOBBY_DEBUG === 'true') {
    return true
  }
  try {
    if (window.localStorage?.getItem('lobbyDebug') === '1') {
      return true
    }
  } catch {
    // private browsing may block storage
  }
  try {
    return new URLSearchParams(window.location.search).get('lobbyDebug') === '1'
  } catch {
    return false
  }
}

/** Log queue/intent diagnostics when lobby debug is enabled. */
export function lobbyDebug(tag, detail = undefined) {
  if (!isLobbyDebugEnabled()) {
    return
  }
  if (detail === undefined) {
    console.info(`[lobby:${tag}]`)
    return
  }
  console.info(`[lobby:${tag}]`, detail)
}
