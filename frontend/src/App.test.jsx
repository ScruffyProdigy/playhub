import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import App from './App'
import { SIGN_IN_HEADING } from './lib/playerCopy'
import { mockAuthenticatedSession, mockUnauthenticatedSession } from './test/setup'

describe('App Component', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/')
  })

  it('renders branding and login when signed out', async () => {
    mockUnauthenticatedSession()
    render(<App />)

    expect(screen.getByRole('heading', { level: 1, name: 'JoinQuest' })).toBeInTheDocument()
    expect(screen.getByText('Find your group. Play together.')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: SIGN_IN_HEADING })).toBeInTheDocument()
    })
    expect(screen.queryByRole('heading', { name: 'Room' })).not.toBeInTheDocument()
  })

  it('shows the signed-in user when authenticated', async () => {
    mockAuthenticatedSession()
    render(<App />)

    expect(await screen.findByText('Welcome back')).toBeInTheDocument()
    expect(screen.getByText('player@example.com')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: 'Log out' })).toBeInTheDocument()
    expect(await screen.findByRole('heading', { name: 'Available games' })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /get started for developers/i })).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: 'Create room' })).not.toBeInTheDocument()
    expect(await screen.findByRole('heading', { name: 'Rock Paper Scissors Lizard Robot' })).toBeInTheDocument()
  })

  it('renders the sign-in link completion page on /auth/complete', async () => {
    window.history.replaceState({}, '', '/auth/complete?token=test-token')
    const assign = vi.fn()
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: { ...window.location, assign, pathname: '/auth/complete', search: '?token=test-token' },
    })

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        data: {
          completeSignInWithLink: {
            id: 'user-1',
            email: 'player@example.com',
            displayName: 'player',
            createdAt: '2026-01-01T00:00:00Z',
          },
        },
      }),
    })

    render(<App />)

    expect(await screen.findByText('Signing you in')).toBeInTheDocument()
  })

  it('renders the OAuth error page on /auth/oauth/complete', async () => {
    window.history.replaceState({}, '', '/auth/oauth/complete?error=provider_error')
    Object.defineProperty(window, 'location', {
      configurable: true,
      value: {
        ...window.location,
        assign: vi.fn(),
        pathname: '/auth/oauth/complete',
        search: '?error=provider_error',
      },
    })

    mockUnauthenticatedSession()

    render(<App />)

    expect(await screen.findByText(/Could not sign in with that provider/i)).toBeInTheDocument()
    expect(screen.queryByText('Signing you in')).not.toBeInTheDocument()
  })
})
