import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, beforeEach } from 'vitest'
import { AuthProvider, useAuth } from './AuthProvider'
import { mockAuthenticatedSession } from '../../test/setup'

function AuthConsumer() {
  const { user, loading } = useAuth()

  if (loading) {
    return <p>Loading</p>
  }

  return <p>{user ? `Signed in as ${user.email}` : 'Signed out'}</p>
}

describe('AuthProvider', () => {
  beforeEach(() => {
    window.history.replaceState({}, '', '/')
  })

  it('shares the current user with useAuth', async () => {
    mockAuthenticatedSession()

    render(
      <AuthProvider>
        <AuthConsumer />
      </AuthProvider>,
    )

    await waitFor(() => {
      expect(screen.getByText('Signed in as player@example.com')).toBeInTheDocument()
    })
  })
})
