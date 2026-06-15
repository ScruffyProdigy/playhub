import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { AuthProvider } from '../auth/AuthProvider'
import DeveloperLandingPage from './DeveloperLandingPage'

vi.mock('../../lib/developers', async (importOriginal) => {
  const actual = await importOriginal()
  return {
    ...actual,
    fetchMyDeveloperApiKeys: vi.fn().mockResolvedValue([]),
  }
})

function renderLandingPage() {
  return render(
    <AuthProvider>
      <DeveloperLandingPage />
    </AuthProvider>,
  )
}

describe('DeveloperLandingPage', () => {
  beforeEach(() => {
    window.history.replaceState(null, '', '/developers')
  })

  afterEach(() => {
    window.history.replaceState(null, '', '/')
  })

  it('opens manual registration when ?path=manual is in the URL', () => {
    window.history.replaceState(null, '', '/developers?path=manual')

    renderLandingPage()

    expect(screen.getByRole('heading', { name: /register your game/i })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: /how do you want to get started/i })).not.toBeInTheDocument()
  })

  it('opens AI assistant setup when ?path=ai is in the URL', () => {
    window.history.replaceState(null, '', '/developers?path=ai')

    renderLandingPage()

    expect(screen.getByRole('heading', { name: /connect an ai assistant/i })).toBeInTheDocument()
    expect(screen.queryByRole('heading', { name: /how do you want to get started/i })).not.toBeInTheDocument()
  })

  it('navigates to manual registration from route picker', async () => {
    const user = userEvent.setup()
    renderLandingPage()

    await user.click(screen.getByRole('button', { name: /register in the browser/i }))

    expect(screen.getByRole('heading', { name: /register your game/i })).toBeInTheDocument()
    expect(window.location.search).toBe('?path=manual')
  })
})
