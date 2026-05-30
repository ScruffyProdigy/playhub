import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import UserSessionCard from './UserSessionCard'
import { AuthProvider } from './AuthProvider'
import * as auth from '../../lib/auth'

vi.mock('../../lib/auth', () => ({
  logout: vi.fn(),
  fetchCurrentUser: vi.fn(),
}))

describe('UserSessionCard', () => {
  const user = {
    id: 'user-1',
    email: 'player@example.com',
    displayName: 'player',
    createdAt: '2026-01-01T00:00:00Z',
  }

  beforeEach(() => {
    vi.mocked(auth.logout).mockReset()
    vi.mocked(auth.logout).mockResolvedValue(true)
    vi.mocked(auth.fetchCurrentUser).mockResolvedValue(user)
  })

  function renderSessionCard() {
    return render(
      <AuthProvider>
        <UserSessionCard user={user} />
      </AuthProvider>,
    )
  }

  it('shows the signed-in user', async () => {
    renderSessionCard()

    expect(await screen.findByText('Welcome back')).toBeInTheDocument()
    expect(screen.getByText('player@example.com')).toBeInTheDocument()
    expect(screen.getByText('player')).toBeInTheDocument()
  })

  it('logs out and clears the session', async () => {
    renderSessionCard()
    await screen.findByText('Welcome back')

    await userEvent.click(screen.getByRole('button', { name: 'Log out' }))

    expect(auth.logout).toHaveBeenCalledTimes(1)
  })
})
