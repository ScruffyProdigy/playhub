import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameLobby from './GameLobby'
import { AuthProvider } from '../auth/AuthProvider'
import { mockAuthenticatedSession, mockUnauthenticatedSession } from '../../test/setup'
import * as games from '../../lib/games'

vi.mock('../../lib/games', () => ({
  fetchGames: vi.fn(),
}))

function renderGameLobby() {
  return render(
    <AuthProvider>
      <GameLobby />
    </AuthProvider>,
  )
}

describe('GameLobby', () => {
  beforeEach(() => {
    vi.mocked(games.fetchGames).mockReset()
    vi.mocked(games.fetchGames).mockResolvedValue([
      {
        id: 'game-1',
        name: 'Quick Match',
        createdAt: '2026-01-01T00:00:00Z',
        activeSessions: [{ id: 'session-1' }],
      },
      {
        id: 'game-2',
        name: 'Party Lobby',
        createdAt: '2026-01-02T00:00:00Z',
        activeSessions: [],
      },
    ])
  })

  it('does not render when the user is logged out', async () => {
    mockUnauthenticatedSession()
    renderGameLobby()

    await waitFor(() => {
      expect(screen.queryByRole('heading', { name: 'Available games' })).not.toBeInTheDocument()
    })
    expect(games.fetchGames).not.toHaveBeenCalled()
  })

  it('loads and shows available games for signed-in users', async () => {
    mockAuthenticatedSession()
    renderGameLobby()

    expect(await screen.findByRole('heading', { name: 'Available games' })).toBeInTheDocument()
    expect(await screen.findByRole('heading', { name: 'Quick Match' })).toBeInTheDocument()
    expect(screen.getByRole('heading', { name: 'Party Lobby' })).toBeInTheDocument()
    expect(screen.getByText('1 active session')).toBeInTheDocument()
    expect(screen.getByText('0 active sessions')).toBeInTheDocument()
    expect(games.fetchGames).toHaveBeenCalledTimes(1)
  })

  it('shows an error when loading games fails', async () => {
    mockAuthenticatedSession()
    vi.mocked(games.fetchGames).mockRejectedValue(new Error('API unavailable'))

    renderGameLobby()

    expect(await screen.findByText('API unavailable')).toBeInTheDocument()
  })
})
