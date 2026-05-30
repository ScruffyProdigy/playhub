import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import App from './App'
import { mockAuthenticatedSession, mockUnauthenticatedSession } from './test/setup'

describe('App Component', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/')
  })

  it('renders branding and login when signed out', async () => {
    mockUnauthenticatedSession()
    render(<App />)

    expect(screen.getByText('PlayHub')).toBeInTheDocument()
    expect(screen.getByText('Your Gaming Hub - Queue, Play, Trade')).toBeInTheDocument()

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })
  })

  it('shows the signed-in user when authenticated', async () => {
    mockAuthenticatedSession()
    render(<App />)

    await waitFor(() => {
      expect(screen.getByText('Welcome back')).toBeInTheDocument()
      expect(screen.getByText('player@example.com')).toBeInTheDocument()
    })
  })

  it('renders the magic-link completion page on /auth/complete', async () => {
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
          completeMagic: {
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
})
