import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import ActiveQueueBanner from './ActiveQueueBanner'

describe('ActiveQueueBanner', () => {
  it('renders nothing without an active queue', () => {
    const { container } = render(
      <ActiveQueueBanner activeQueue={null} busy={false} onLeave={vi.fn()} />,
    )
    expect(container).toBeEmptyDOMElement()
  })

  it('shows waiting state with stop looking', () => {
    render(
      <ActiveQueueBanner
        activeQueue={{
          queueId: 'q1',
          gameName: 'Demo Game',
          status: 'WAITING',
          queuedCount: 3,
        }}
        busy={false}
        onLeave={vi.fn()}
      />,
    )
    expect(screen.getByText(/Looking for a group in Demo Game/)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Stop looking' })).toBeInTheDocument()
  })

  it('shows which cohort the player joined in composition games', () => {
    render(
      <ActiveQueueBanner
        activeQueue={{
          queueId: 'q1',
          gameName: 'Word Hunt',
          status: 'WAITING',
          queuedCount: 2,
          queuePath: 'ClueGiver',
          queuePathDisplayName: 'Clue Giver',
        }}
        busy={false}
        onLeave={vi.fn()}
      />,
    )
    expect(
      screen.getByText('Looking for a group in Word Hunt as Clue Giver · 2 players looking'),
    ).toBeInTheDocument()
  })

  it('shows matched state with launch link', () => {
    render(
      <ActiveQueueBanner
        activeQueue={{
          queueId: 'q1',
          gameName: 'Demo Game',
          status: 'MATCHED',
          joinUrl: 'http://game.example/play',
        }}
        busy={false}
        onLeave={vi.fn()}
      />,
    )
    expect(screen.getByText(/Your group is ready/)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Launch game' })).toHaveAttribute(
      'href',
      'http://game.example/play',
    )
  })

  it('calls onLeave from the banner action', async () => {
    const onLeave = vi.fn()
    const user = userEvent.setup()
    render(
      <ActiveQueueBanner
        activeQueue={{
          queueId: 'q1',
          gameName: 'Demo Game',
          status: 'WAITING',
          queuedCount: 1,
        }}
        busy={false}
        onLeave={onLeave}
      />,
    )
    await user.click(screen.getByRole('button', { name: 'Stop looking' }))
    expect(onLeave).toHaveBeenCalledTimes(1)
  })
})
