import { render, screen } from '@testing-library/react'
import { describe, it, expect } from 'vitest'
import GameListItem from './GameListItem'

describe('GameListItem', () => {
  it('shows the game name and session count', () => {
    render(
      <ul>
        <GameListItem
          game={{
            id: 'game-1',
            name: 'Quick Match',
            activeSessions: [{ id: 'session-1' }],
          }}
        />
      </ul>,
    )

    expect(screen.getByRole('heading', { name: 'Quick Match' })).toBeInTheDocument()
    expect(screen.getByText('1 active session')).toBeInTheDocument()
  })

  it('pluralizes the session count', () => {
    render(
      <ul>
        <GameListItem
          game={{
            id: 'game-2',
            name: 'Party Lobby',
            activeSessions: [],
          }}
        />
      </ul>,
    )

    expect(screen.getByText('0 active sessions')).toBeInTheDocument()
  })
})
