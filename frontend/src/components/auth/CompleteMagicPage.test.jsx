import { render, screen, waitFor } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import CompleteMagicPage from './CompleteMagicPage'
import * as auth from '../../lib/auth'

vi.mock('../../lib/auth', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    completeMagicLoginOnce: vi.fn(),
  }
})

describe('CompleteMagicPage', () => {
  const assign = vi.fn()
  const replaceState = vi.fn()

  beforeEach(() => {
    vi.mocked(auth.completeMagicLoginOnce).mockReset()
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
    render(<CompleteMagicPage />)

    expect(await screen.findByText('Missing sign-in token. Request a new magic link.')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Back to home' })).toBeInTheDocument()
    expect(auth.completeMagicLoginOnce).not.toHaveBeenCalled()
  })

  it('shows an error when completing sign-in fails', async () => {
    window.location.search = '?token=expired-token'
    vi.mocked(auth.completeMagicLoginOnce).mockRejectedValue(
      new Error('auth: invalid or expired magic link'),
    )

    render(<CompleteMagicPage />)

    expect(await screen.findByText('auth: invalid or expired magic link')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Back to home' })).toBeInTheDocument()
    expect(assign).not.toHaveBeenCalled()
  })

  it('redirects home after a successful sign-in', async () => {
    window.location.search = '?token=valid-token'
    vi.mocked(auth.completeMagicLoginOnce).mockResolvedValue({
      id: 'user-1',
      email: 'player@example.com',
    })

    render(<CompleteMagicPage />)

    await waitFor(() => {
      expect(auth.completeMagicLoginOnce).toHaveBeenCalledWith('valid-token')
      expect(replaceState).toHaveBeenCalledWith({}, '', '/')
      expect(assign).toHaveBeenCalledWith('/')
    })
  })
})
