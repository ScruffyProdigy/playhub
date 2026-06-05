import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import GameQueueActions from './GameQueueActions'

describe('GameQueueActions', () => {
  it('shows a single Look for group button for fifo modes', () => {
    render(
      <GameQueueActions
        joinOptions={{ kind: 'fifo', paths: [] }}
        queueState="idle"
        busy={false}
        onJoin={vi.fn()}
        onLeave={vi.fn()}
      />,
    )

    expect(screen.getByRole('button', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.queryByRole('region', { name: 'Look for group' })).not.toBeInTheDocument()
  })

  it('shows role buttons inside a Look for group panel for composition modes', () => {
    render(
      <GameQueueActions
        joinOptions={{ kind: 'composition', paths: ['DPS', 'Support', 'Tank'] }}
        queueState="idle"
        busy={false}
        onJoin={vi.fn()}
        onLeave={vi.fn()}
      />,
    )

    expect(screen.getByRole('region', { name: 'Look for group' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as DPS' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Support' })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Join as Tank' })).toBeInTheDocument()
  })

  it('calls onJoin with the selected queue path', async () => {
    const onJoin = vi.fn()
    render(
      <GameQueueActions
        joinOptions={{ kind: 'composition', paths: ['DPS', 'Tank'] }}
        queueState="idle"
        busy={false}
        onJoin={onJoin}
        onLeave={vi.fn()}
      />,
    )

    await userEvent.click(screen.getByRole('button', { name: 'Join as Tank' }))
    expect(onJoin).toHaveBeenCalledWith('Tank')
  })

  it('shows waiting state with selected role in composition panel', () => {
    render(
      <GameQueueActions
        joinOptions={{ kind: 'composition', paths: ['DPS'] }}
        queueState="waiting"
        selectedQueuePath="DPS"
        busy={false}
        onJoin={vi.fn()}
        onLeave={vi.fn()}
      />,
    )

    expect(screen.getByText('Looking as DPS…')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Stop looking' })).toBeInTheDocument()
  })
})
