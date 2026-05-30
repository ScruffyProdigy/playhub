import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import LoginForm from './LoginForm'
import { AuthProvider } from './AuthProvider'
import * as auth from '../../lib/auth'

vi.mock('../../lib/auth', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    requestMagicLink: vi.fn(),
    fetchCurrentUser: vi.fn().mockResolvedValue(null),
  }
})

function renderLoginForm() {
  return render(
    <AuthProvider>
      <LoginForm />
    </AuthProvider>,
  )
}

describe('LoginForm', () => {
  beforeEach(() => {
    vi.mocked(auth.requestMagicLink).mockReset()
    vi.mocked(auth.requestMagicLink).mockResolvedValue(true)
    vi.mocked(auth.fetchCurrentUser).mockReset()
    vi.mocked(auth.fetchCurrentUser).mockResolvedValue(null)
  })

  it('submits a magic link request for the entered email', async () => {
    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Send magic link' }))

    await waitFor(() => {
      expect(auth.requestMagicLink).toHaveBeenCalledWith('player@example.com')
      expect(screen.getByText('Check your email for your sign-in link.')).toBeInTheDocument()
    })
  })

  it('shows an error when the magic link request fails', async () => {
    vi.mocked(auth.requestMagicLink).mockRejectedValue(new Error('SMTP unavailable'))

    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Send magic link' }))

    expect(await screen.findByText('SMTP unavailable')).toBeInTheDocument()
  })

  it('refreshes the session when the user already signed in elsewhere', async () => {
    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'I already signed in' })).toBeInTheDocument()
    })

    await userEvent.click(screen.getByRole('button', { name: 'I already signed in' }))

    await waitFor(() => {
      expect(auth.fetchCurrentUser).toHaveBeenCalled()
    })
  })
})
