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
      seats: [{ queuePath: null }, { queuePath: null }],
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
        <GameListItem game={catalogGame} activeQueue={null} onQueueChange={vi.fn()} />
      </ul>,
    )

    expect(screen.getByRole('heading', { name: 'Rock Paper Scissors Lizard Spock' })).toBeInTheDocument()
    expect(screen.getByText('1 active session')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Look for group' })).toBeInTheDocument()
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
        <GameListItem game={catalogGame} activeQueue={null} onQueueChange={vi.fn()} />
      </ul>,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Look for group' }))

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

  it('refreshes the active-queue banner after joining', async () => {
    const onQueueChange = vi.fn()
    vi.mocked(queue.joinQueue).mockResolvedValue({ queued: true, queuedCount: 1 })

    render(
      <ul>
        <GameListItem game={catalogGame} activeQueue={null} onQueueChange={onQueueChange} />
      </ul>,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Look for group' }))

    await waitFor(() => expect(onQueueChange).toHaveBeenCalled())
  })

  it('shows a plain Look for group button for fifo modes with empty queuePaths', () => {
    const fifoGame = {
      ...catalogGame,
      modes: [
        {
          modeKey: 'duel',
          queuePaths: [{ queuePath: '', displayName: '', playersToStart: 2 }],
          seats: [{ queuePath: null }, { queuePath: null }],
          queues: [{ id: 'queue-1', name: 'Default', status: 'active' }],
        },
      ],
    }

    render(
      <ul>
        <GameListItem game={fifoGame} activeQueue={null} onQueueChange={vi.fn()} />
      </ul>,
    )

    expect(screen.queryByRole('region', { name: 'Look for group' })).not.toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Join as' })).not.toBeInTheDocument()
  })

  it('keeps fifo join controls visible while waiting in this queue', async () => {
    vi.mocked(queue.fetchMyQueueStatus).mockResolvedValue({
      queued: true,
      queuedCount: 1,
    })

    render(
      <ul>
        <GameListItem
          game={catalogGame}
          activeQueue={{
            queueId: 'queue-1',
            gameName: 'Rock Paper Scissors Lizard Spock',
            status: 'WAITING',
            queuedCount: 1,
          }}
          onQueueChange={vi.fn()}
        />
      </ul>,
    )

    await waitFor(() => {
      expect(screen.getByText('Looking…')).toBeInTheDocument()
    })
    expect(screen.queryByRole('region', { name: 'Look for group' })).not.toBeInTheDocument()
    expect(screen.queryByText('Use the banner above')).not.toBeInTheDocument()
  })

  it('keeps composition join actions visible while waiting in this queue', async () => {
    const wordHuntGame = {
      id: 'game-wh',
      name: 'Word Hunt',
      activeSessions: [],
      modes: [
        {
          modeKey: 'party',
          queuePaths: [
            { queuePath: 'ClueGiver', displayName: 'Clue Giver' },
            { queuePath: 'Guesser', displayName: 'Guesser' },
          ],
          seats: [{ queuePath: 'ClueGiver' }, { queuePath: 'Guesser' }],
          queues: [{ id: 'queue-wh', name: 'Default', status: 'active' }],
        },
      ],
    }

    vi.mocked(queue.fetchMyQueueStatus).mockResolvedValue({
      queued: true,
      queuedCount: 2,
      queuePath: 'ClueGiver',
    })

    render(
      <ul>
        <GameListItem
          game={wordHuntGame}
          activeQueue={{
            queueId: 'queue-wh',
            gameName: 'Word Hunt',
            status: 'WAITING',
            queuedCount: 2,
            queuePath: 'ClueGiver',
          }}
          onQueueChange={vi.fn()}
        />
      </ul>,
    )

    await waitFor(() => {
      expect(screen.getByText('Looking as Clue Giver…')).toBeInTheDocument()
    })
    expect(screen.getByRole('button', { name: 'Join as Guesser' })).toBeInTheDocument()
    expect(screen.queryByText('Use the banner above')).not.toBeInTheDocument()
  })

  it('shows composition join buttons inside a Look for group panel', () => {
    const compositionGame = {
      ...catalogGame,
      modes: [
        {
          modeKey: 'quick-play',
          seats: [
            { queuePath: 'DPS' },
            { queuePath: 'Tank' },
            { queuePath: 'Support' },
          ],
          queues: [{ id: 'queue-comp', name: 'Default', status: 'active' }],
        },
      ],
    }

    render(
      <ul>
        <GameListItem game={compositionGame} activeQueue={null} onQueueChange={vi.fn()} />
      </ul>,
    )

    expect(screen.getByRole('region', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as DPS' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Tank' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Support' })).toBeInTheDocument()
  })
})
