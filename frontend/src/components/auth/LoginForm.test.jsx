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
    requestSignIn: vi.fn(),
    completeSignInWithCode: vi.fn(),
    fetchCurrentUser: vi.fn().mockResolvedValue(null),
  }
})

vi.mock('../../lib/authBroadcast', () => ({
  onAuthComplete: vi.fn(() => () => {}),
}))

function renderLoginForm() {
  return render(
    <AuthProvider>
      <LoginForm />
    </AuthProvider>,
  )
}

describe('LoginForm', () => {
  beforeEach(() => {
    vi.mocked(auth.requestSignIn).mockReset()
    vi.mocked(auth.requestSignIn).mockResolvedValue(true)
    vi.mocked(auth.completeSignInWithCode).mockReset()
    vi.mocked(auth.fetchCurrentUser).mockReset()
    vi.mocked(auth.fetchCurrentUser).mockResolvedValue(null)
  })

  it('submits a sign-in email request and shows the code step', async () => {
    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    await waitFor(() => {
      expect(auth.requestSignIn).toHaveBeenCalledWith('player@example.com')
      expect(screen.getByRole('heading', { name: 'Enter your code' })).toBeInTheDocument()
    })
    expect(screen.getByLabelText('Sign-in code')).not.toBeDisabled()
  })

  it('shows an error when the sign-in email request fails', async () => {
    vi.mocked(auth.requestSignIn).mockRejectedValue(new Error('SMTP unavailable'))

    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    expect(await screen.findByText('SMTP unavailable')).toBeInTheDocument()
  })

  it('submits a 6-digit code to complete sign-in', async () => {
    vi.mocked(auth.completeSignInWithCode).mockResolvedValue({
      id: 'user-1',
      email: 'player@example.com',
    })

    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    const codeInput = await screen.findByLabelText('Sign-in code')
    expect(codeInput).not.toBeDisabled()
    await userEvent.type(codeInput, '123456')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    await waitFor(() => {
      expect(auth.completeSignInWithCode).toHaveBeenCalledWith('player@example.com', '123456')
    })
  })

  it('shows a friendly error when the sign-in code is invalid', async () => {
    vi.mocked(auth.completeSignInWithCode).mockRejectedValue(
      new Error('Invalid or expired code. Try again or use the sign-in link in your email.'),
    )

    renderLoginForm()
    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Sign in' })).toBeInTheDocument()
    })

    await userEvent.type(screen.getByLabelText('Email'), 'player@example.com')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    const codeInput = await screen.findByLabelText('Sign-in code')
    expect(codeInput).not.toBeDisabled()
    await userEvent.type(codeInput, '123456')
    await userEvent.click(screen.getByRole('button', { name: 'Continue' }))

    expect(
      await screen.findByText('Invalid or expired code. Try again or use the sign-in link in your email.'),
    ).toBeInTheDocument()
  })
})
