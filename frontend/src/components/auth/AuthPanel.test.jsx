import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import AuthPanel from './AuthPanel'
import { AuthProvider } from './AuthProvider'
import { mockAuthenticatedSession, mockUnauthenticatedSession } from '../../test/setup'

function renderAuthPanel() {
  return render(
    <AuthProvider>
      <AuthPanel />
    </AuthProvider>,
  )
}

describe('AuthPanel', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/')
  })

  it('shows a loading state while checking the session', () => {
    mockUnauthenticatedSession()
    renderAuthPanel()

    expect(screen.getByText('Loading session…')).toBeInTheDocument()
  })

  it('shows the sign-in form when logged out', async () => {
    mockUnauthenticatedSession()
    renderAuthPanel()

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Get in the game' })).toBeInTheDocument()
    })
  })

  it('shows the signed-in user when authenticated', async () => {
    mockAuthenticatedSession()
    renderAuthPanel()

    await waitFor(() => {
      expect(screen.getByText('Welcome back')).toBeInTheDocument()
      expect(screen.getByText('player@example.com')).toBeInTheDocument()
    })
  })
})
