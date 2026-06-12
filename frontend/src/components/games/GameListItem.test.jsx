import { render, screen, waitFor, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameListItem from './GameListItem'
import * as queue from '../../lib/queue'

vi.mock('../rooms/ActiveRoomProvider', () => ({
  useActiveRoom: () => ({
    refresh: vi.fn(),
    openRoom: vi.fn(),
  }),
}))

vi.mock('../../lib/tables', () => ({
  createPrivateTable: vi.fn(),
}))

vi.mock('../../lib/queue', () => ({
  joinQueue: vi.fn(),
  fetchMyQueueStatus: vi.fn(),
  leaveQueue: vi.fn(),
  prefetchSubscriptionAuth: vi.fn().mockResolvedValue('Bearer test'),
  subscribeToQueue: vi.fn().mockResolvedValue(() => {}),
}))

const activeMode = {
  id: 'mode-1',
  modeKey: 'duel',
  displayName: 'Duel',
  status: 'active',
  seats: [{ queuePath: null }, { queuePath: null }],
  queues: [{ id: 'queue-1', name: 'Default', status: 'active' }],
}

const catalogGame = {
  id: 'game-1',
  slug: 'rock-paper-scissors-lizard-robot',
  name: 'Rock Paper Scissors Lizard Robot',
  iconUrl: '/games/rpslr-icon.png',
  heroUrl: '/games/rpslr-hero.jpg',
  shortDescription: 'Best-of-five duel with move cooldown.',
  tags: ['competitive', '1v1'],
  activeSessions: [{ id: 'session-1' }],
  modes: [activeMode],
}

function renderItem(props) {
  return render(
    <ul>
      <GameListItem
        game={catalogGame}
        activeIntent={null}
        activeTableSeat={null}
        onQueueChange={vi.fn()}
        onTableChange={vi.fn()}
        {...props}
      />
    </ul>,
  )
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

  it('shows the game name and catalog details', () => {
    const { container } = renderItem()

    expect(screen.getByRole('heading', { name: 'Rock Paper Scissors Lizard Robot' })).toBeInTheDocument()
    expect(screen.getByText('Best-of-five duel with move cooldown.')).toBeInTheDocument()
    expect(screen.getByText('competitive')).toBeInTheDocument()
    expect(container.querySelector('.game-list-item__hero')).toHaveAttribute(
      'src',
      '/games/rpslr-hero.jpg?v=1',
    )
    expect(screen.getByRole('link', { name: 'Rock Paper Scissors Lizard Robot' })).toHaveAttribute(
      'href',
      '/games/rock-paper-scissors-lizard-robot',
    )
    expect(screen.getByRole('button', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Create private game' })).toBeInTheDocument()
  })

  it('uses catalogHeroUrl on the card when set', () => {
    const { container } = renderItem({
      game: {
        ...catalogGame,
        catalogHeroUrl: '/games/rpslr-catalog-hero.png',
      },
    })

    expect(container.querySelector('.game-list-item__hero')).toHaveAttribute(
      'src',
      '/games/rpslr-catalog-hero.png?v=1',
    )
  })

  it('shows launch when the subscription reports a match', async () => {
    vi.mocked(queue.joinQueue).mockResolvedValue({ queued: true, queuedCount: 1 })
    let onUpdate
    vi.mocked(queue.subscribeToQueue).mockImplementation(async (_queueId, handlers) => {
      onUpdate = handlers.onUpdate
      return () => {}
    })

    renderItem()

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

  it('notifies intent after joining', async () => {
    const onQueueJoined = vi.fn()
    vi.mocked(queue.joinQueue).mockResolvedValue({ queued: true, queuedCount: 1 })

    renderItem({ onQueueJoined })

    await userEvent.click(screen.getByRole('button', { name: 'Look for group' }))

    await waitFor(() =>
      expect(onQueueJoined).toHaveBeenCalledWith(
        'queue-1',
        { queued: true, queuedCount: 1 },
        expect.objectContaining({ gameId: 'game-1', gameName: catalogGame.name }),
      ),
    )
  })

  it('shows waiting controls when activeIntent matches this queue', () => {
    renderItem({
      activeIntent: {
        queueId: 'queue-1',
        gameId: 'game-1',
        gameName: catalogGame.name,
        status: 'WAITING',
        queuedCount: 2,
      },
    })

    expect(screen.getByText('Looking…')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Stop looking' })).toBeInTheDocument()
  })

  it('shows a plain Look for group button for fifo modes with empty queuePaths', () => {
    const fifoGame = {
      ...catalogGame,
      modes: [
        {
          ...activeMode,
          queuePaths: [{ queuePath: '', displayName: '', playersToStart: 2 }],
        },
      ],
    }

    render(
      <ul>
        <GameListItem
          game={fifoGame}
          activeIntent={null}
          activeTableSeat={null}
          onQueueChange={vi.fn()}
          onTableChange={vi.fn()}
        />
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

    renderItem({
      activeIntent: {
        queueId: 'queue-1',
        gameName: 'Rock Paper Scissors Lizard Robot',
        status: 'WAITING',
        queuedCount: 1,
      },
    })

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
      iconUrl: '/games/word-hunt-icon.png',
      heroUrl: '/games/word-hunt-hero.jpg',
      activeSessions: [],
      modes: [
        {
          id: 'mode-wh',
          modeKey: 'party',
          displayName: 'Party',
          status: 'active',
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
          activeIntent={{
            queueId: 'queue-wh',
            gameName: 'Word Hunt',
            status: 'WAITING',
            queuedCount: 2,
            queuePath: 'ClueGiver',
          }}
          activeTableSeat={null}
          onQueueChange={vi.fn()}
          onTableChange={vi.fn()}
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
          id: 'mode-comp',
          modeKey: 'quick-play',
          displayName: 'Quick Play',
          status: 'active',
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
        <GameListItem
          game={compositionGame}
          activeIntent={null}
          activeTableSeat={null}
          onQueueChange={vi.fn()}
          onTableChange={vi.fn()}
        />
      </ul>,
    )

    expect(screen.getByRole('region', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as DPS' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Tank' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Support' })).toBeInTheDocument()
  })
})
