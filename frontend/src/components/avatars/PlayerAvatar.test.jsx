import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import PlayerAvatar from './PlayerAvatar'

describe('PlayerAvatar', () => {
  it('renders image when avatarUrl is set', () => {
    const { container } = render(
      <PlayerAvatar user={{ displayName: 'Pat', avatarUrl: '/avatars/compass.png' }} />,
    )
    expect(container.querySelector('img.player-avatar')).toBeTruthy()
  })

  it('falls back to initial without avatarUrl', () => {
    render(<PlayerAvatar user={{ displayName: 'River' }} />)
    expect(screen.getByText('R')).toBeTruthy()
  })
})
