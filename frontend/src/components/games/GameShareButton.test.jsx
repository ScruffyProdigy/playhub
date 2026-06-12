import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import GameShareButton from './GameShareButton'

describe('GameShareButton', () => {
  beforeEach(() => {
    vi.stubGlobal('navigator', {
      ...navigator,
      share: undefined,
      clipboard: { writeText: vi.fn().mockResolvedValue(undefined) },
    })
  })

  it('copies the share URL when Web Share is unavailable', async () => {
    render(<GameShareButton game={{ slug: 'word-hunt', name: 'Word Hunt' }} />)

    await userEvent.click(screen.getByRole('button', { name: 'Share' }))

    await waitFor(() => {
      expect(navigator.clipboard.writeText).toHaveBeenCalledWith(
        expect.stringMatching(/\/games\/word-hunt$/),
      )
    })
    expect(screen.getByText('Link copied!')).toBeInTheDocument()
  })

  it('uses native share when available', async () => {
    const share = vi.fn().mockResolvedValue(undefined)
    vi.stubGlobal('navigator', {
      ...navigator,
      share,
      clipboard: { writeText: vi.fn() },
    })

    render(<GameShareButton game={{ slug: 'word-hunt', name: 'Word Hunt' }} />)
    await userEvent.click(screen.getByRole('button', { name: 'Share' }))

    expect(share).toHaveBeenCalledWith(
      expect.objectContaining({
        title: 'Word Hunt',
        url: expect.stringMatching(/\/games\/word-hunt$/),
      }),
    )
  })
})
