import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameListItem from './GameListItem'
import * as queue from '../../lib/queue'

vi.mock('../../lib/queue', () => ({
  joinQueue: vi.fn(),
  fetchMyQueueStatus: vi.fn(),
  leaveQueue: vi.fn(),
  prefetchSubscriptionAuth: vi.fn().mockResolvedValue('Bearer test'),
  subscribeToQueue: vi.fn().mockResolvedValue(() => {}),
}))

const catalogGame = {
  id: 'game-1',
  name: 'Rock Paper Scissors Lizard Spock',
  activeSessions: [{ id: 'session-1' }],
  modes: [
    {
      modeKey: 'duel',
      queues: [{ id: 'queue-1', name: 'Default', status: 'active' }],
    },
  ],
}

describe('GameListItem', () => {
  beforeEach(() => {
    vi.mocked(queue.joinQueue).mockReset()
    vi.mocked(queue.fetchMyQueueStatus).mockReset()
    vi.mocked(queue.leaveQueue).mockReset()
    vi.mocked(queue.subscribeToQueue).mockReset()
    vi.mocked(queue.subscribeToQueue).mockReturnValue(() => {})
    vi.mocked(queue.fetchMyQueueStatus).mockResolvedValue({ queued: false })
  })

  it('shows the game name and session count', () => {
    render(
      <ul>
        <GameListItem game={catalogGame} />
      </ul>,
    )

    expect(screen.getByRole('heading', { name: 'Rock Paper Scissors Lizard Spock' })).toBeInTheDocument()
    expect(screen.getByText('1 active session')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join queue' })).toBeInTheDocument()
  })

  it('shows launch when the subscription reports a match', async () => {
    vi.mocked(queue.joinQueue).mockResolvedValue({ queued: true, queuedCount: 1 })
    let onUpdate
    vi.mocked(queue.subscribeToQueue).mockImplementation(async (_queueId, handlers) => {
      onUpdate = handlers.onUpdate
      return () => {}
    })

    render(
      <ul>
        <GameListItem game={catalogGame} />
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
