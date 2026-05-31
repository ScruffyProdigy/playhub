import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import CompleteSignInPage from './CompleteSignInPage'
import * as auth from '../../lib/auth'

vi.mock('../../lib/auth', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    completeSignInWithLinkOnce: vi.fn(),
  }
})

describe('CompleteSignInPage', () => {
  const assign = vi.fn()
  const replaceState = vi.fn()

  beforeEach(() => {
    vi.mocked(auth.completeSignInWithLinkOnce).mockReset()
    assign.mockReset()
    replaceState.mockReset()

    vi.stubGlobal('location', {
      ...window.location,
      assign,
      search: '',
    })
    window.history.replaceState = replaceState
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('shows an error when the sign-in token is missing', async () => {
    render(<CompleteSignInPage />)

    expect(await screen.findByText('Missing sign-in token. Request a new sign-in email.')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Back to home' })).toBeInTheDocument()
    expect(auth.completeSignInWithLinkOnce).not.toHaveBeenCalled()
  })

  it('shows an error when completing sign-in fails', async () => {
    window.location.search = '?token=expired-token'
    vi.mocked(auth.completeSignInWithLinkOnce).mockRejectedValue(
      new Error('Invalid or expired sign-in link. Request a new sign-in email.'),
    )

    render(<CompleteSignInPage />)

    expect(await screen.findByText('Invalid or expired sign-in link. Request a new sign-in email.')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Back to home' })).toBeInTheDocument()
    expect(assign).not.toHaveBeenCalled()
  })

  it('redirects home after a successful sign-in', async () => {
    window.location.search = '?token=valid-token'
    vi.mocked(auth.completeSignInWithLinkOnce).mockResolvedValue({
      id: 'user-1',
      email: 'player@example.com',
    })

    render(<CompleteSignInPage />)

    await waitFor(() => {
      expect(auth.completeSignInWithLinkOnce).toHaveBeenCalledWith('valid-token')
      expect(replaceState).toHaveBeenCalledWith({}, '', '/')
      expect(assign).toHaveBeenCalledWith('/')
    })
  })
})
