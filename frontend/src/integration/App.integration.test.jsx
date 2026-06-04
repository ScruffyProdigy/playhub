import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import App from '../App'
import { mockAuthenticatedSession, mockUnauthenticatedSession } from '../test/setup'

describe('App Integration Tests', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/')
  })

  describe('User Journey: Content Discovery', () => {
    it('presents branding and sign-in when logged out', async () => {
      mockUnauthenticatedSession()
      render(<App />)

      expect(screen.getByRole('heading', { name: 'JoinQuest' })).toBeInTheDocument()
      expect(screen.getByText('Find your group. Play together.')).toBeInTheDocument()

      await waitFor(() => {
        expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
      })
    })

    it('shows account details when logged in', async () => {
      mockAuthenticatedSession()
      render(<App />)

      await waitFor(() => {
        expect(screen.getByText('player@example.com')).toBeInTheDocument()
        expect(screen.getByRole('heading', { name: 'Available games' })).toBeInTheDocument()
        expect(screen.getByRole('heading', { name: 'Rock Paper Scissors Lizard Spock' })).toBeInTheDocument()
      })
    })
  })
})
