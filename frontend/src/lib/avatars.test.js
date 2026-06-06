import { describe, it, expect } from 'vitest'
import {
  avatarInitial,
  defaultDisplayNameInput,
  isProvisionalDisplayName,
  needsProfileSetup,
  PROVISIONAL_DISPLAY_SUFFIX,
  STARTER_AVATAR_FALLBACK,
} from './avatars'

describe('avatars', () => {
  it('builds initials from display name', () => {
    expect(avatarInitial({ displayName: 'Pat' })).toBe('P')
    expect(avatarInitial({})).toBe('P')
  })

  it('ships five journey starter avatars', () => {
    expect(STARTER_AVATAR_FALLBACK).toHaveLength(5)
  })

  it('detects provisional display names', () => {
    expect(isProvisionalDisplayName(`pat${PROVISIONAL_DISPLAY_SUFFIX}`)).toBe(true)
    expect(isProvisionalDisplayName('Pat')).toBe(false)
  })

  it('requires profile setup without avatar or provisional name', () => {
    expect(needsProfileSetup({ displayName: 'Pat', avatarKey: 'compass' })).toBe(false)
    expect(needsProfileSetup({ displayName: `pat${PROVISIONAL_DISPLAY_SUFFIX}` })).toBe(true)
    expect(needsProfileSetup({ displayName: 'Pat' })).toBe(true)
  })

  it('prefills display name without provisional suffix', () => {
    expect(defaultDisplayNameInput({ displayName: `river${PROVISIONAL_DISPLAY_SUFFIX}` })).toBe('river')
    expect(defaultDisplayNameInput({ displayName: 'River' })).toBe('River')
  })
})
