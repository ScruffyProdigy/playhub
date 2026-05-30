import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameListItem from './GameListItem'
import * as queue from '../../lib/queue'

vi.mock('../../lib/queue', () => ({
  joinGame: vi.fn(),
  fetchMyQueueStatus: vi.fn(),
  leaveQueue: vi.fn(),
  prefetchSubscriptionAuth: vi.fn().mockResolvedValue('Bearer test'),
  subscribeToQueue: vi.fn().mockResolvedValue(() => {}),
}))

describe('GameListItem', () => {
  beforeEach(() => {
    vi.mocked(queue.joinGame).mockReset()
    vi.mocked(queue.fetchMyQueueStatus).mockReset()
    vi.mocked(queue.leaveQueue).mockReset()
    vi.mocked(queue.subscribeToQueue).mockReset()
    vi.mocked(queue.subscribeToQueue).mockReturnValue(() => {})
    vi.mocked(queue.fetchMyQueueStatus).mockResolvedValue({ queued: false })
  })

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
    expect(screen.getByRole('button', { name: 'Join queue' })).toBeInTheDocument()
  })

  it('shows launch when the subscription reports a match', async () => {
    vi.mocked(queue.joinGame).mockResolvedValue({ queued: true, queuedCount: 1 })
    let onUpdate
    vi.mocked(queue.subscribeToQueue).mockImplementation(async (_gameId, handlers) => {
      onUpdate = handlers.onUpdate
      return () => {}
    })

    render(
      <ul>
        <GameListItem
          game={{
            id: 'game-1',
            name: 'Quick Match',
            activeSessions: [],
          }}
        />
      </ul>,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Join queue' }))

    await waitFor(() => expect(onUpdate).toBeDefined())

    act(() => {
      onUpdate({
        status: 'MATCHED',
        joinUrl: 'http://localhost:5174/?match=session-1&token=eyJ.test',
        queuedCount: 0,
      })
    })

    expect(await screen.findByRole('link', { name: 'Launch game' })).toHaveAttribute(
      'href',
      'http://localhost:5174/?match=session-1&token=eyJ.test',
    )
  })
})
