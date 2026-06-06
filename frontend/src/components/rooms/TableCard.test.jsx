import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi } from 'vitest'
import TableCard from './TableCard'

vi.mock('../auth/AuthProvider', () => ({
  useAuth: () => ({
    user: {
      id: 'user-1',
      displayName: 'Pat',
      avatarKey: 'campfire',
      avatarUrl: '/avatars/campfire.png',
    },
  }),
}))

describe('TableCard', () => {
  it('shows the signed-in user avatar on role-layout seats', () => {
    const table = {
      id: 'table-1',
      game: { name: 'Codenames' },
      mode: { displayName: 'Party' },
      seats: [{ seatKey: 'ClueGiver-Red', user: { id: 'user-1', displayName: 'Pat' } }],
      seatSlots: [
        {
          seatKey: 'ClueGiver-Red',
          queuePath: 'ClueGiver',
          displayName: 'Clue Giver · Red',
          user: { id: 'user-1', displayName: 'Pat' },
        },
        {
          seatKey: 'Guesser-1',
          queuePath: 'Guesser',
          displayName: 'Guesser · 1',
          user: null,
        },
      ],
      lookForGroupOptions: [],
    }

    const { container } = render(
      <TableCard table={table} busy={false} onSit={() => {}} onLeave={() => {}} onStart={() => {}} onDiscard={() => {}} />,
    )

    const avatar = container.querySelector('img.player-avatar[src="/avatars/campfire.png"]')
    expect(avatar).toBeTruthy()
    expect(screen.getByText('You')).toBeInTheDocument()
  })
})
