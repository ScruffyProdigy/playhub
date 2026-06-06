import { describe, it, expect } from 'vitest'
import { avatarInitial, STARTER_AVATAR_FALLBACK } from './avatars'

describe('avatars', () => {
  it('builds initials from display name', () => {
    expect(avatarInitial({ displayName: 'Pat' })).toBe('P')
    expect(avatarInitial({})).toBe('P')
  })

  it('ships five journey starter avatars', () => {
    expect(STARTER_AVATAR_FALLBACK).toHaveLength(5)
    expect(STARTER_AVATAR_FALLBACK.map((item) => item.key)).toEqual([
      'compass',
      'coin',
      'storm',
      'campfire',
      'beacon',
    ])
  })
})
