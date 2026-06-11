import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest'
import { isLobbyDebugEnabled, lobbyDebug } from './lobbyDebug'

describe('lobbyDebug', () => {
  beforeEach(() => {
    window.env = { REACT_APP_LOBBY_DEBUG: '' }
    window.localStorage.clear()
    vi.spyOn(console, 'info').mockImplementation(() => {})
  })

  afterEach(() => {
    vi.restoreAllMocks()
    delete window.location
    window.location = new URL('https://joinquest.cc/')
  })

  it('is off by default', () => {
    expect(isLobbyDebugEnabled()).toBe(false)
  })

  it('enables via REACT_APP_LOBBY_DEBUG', () => {
    window.env.REACT_APP_LOBBY_DEBUG = 'true'
    expect(isLobbyDebugEnabled()).toBe(true)
  })

  it('enables via localStorage', () => {
    window.localStorage.setItem('lobbyDebug', '1')
    expect(isLobbyDebugEnabled()).toBe(true)
  })

  it('enables via ?lobbyDebug=1', () => {
    window.location = new URL('https://joinquest.cc/?lobbyDebug=1')
    expect(isLobbyDebugEnabled()).toBe(true)
  })

  it('logs when enabled', () => {
    window.env.REACT_APP_LOBBY_DEBUG = 'true'
    lobbyDebug('test:tag', { ok: true })
    expect(console.info).toHaveBeenCalledWith('[lobby:test:tag]', { ok: true })
  })

  it('does not log when disabled', () => {
    lobbyDebug('test:tag', { ok: true })
    expect(console.info).not.toHaveBeenCalled()
  })
})
