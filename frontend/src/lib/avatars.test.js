import { describe, it, expect } from 'vitest'
import {
  avatarInitial,
  defaultDisplayNameInput,
  hasExistingAvatar,
  isProvisionalDisplayName,
  needsProfileSetup,
  PROVISIONAL_DISPLAY_SUFFIX,
  resolveUserAvatarUrl,
  STARTER_AVATAR_FALLBACK,
} from './avatars'

describe('avatars', () => {
  it('builds initials from display name', () => {
    expect(avatarInitial({ displayName: 'Pat' })).toBe('P')
    expect(avatarInitial({})).toBe('P')
  })

  it('ships eighteen starter avatars', () => {
    expect(STARTER_AVATAR_FALLBACK).toHaveLength(18)
  })

  it('detects provisional display names', () => {
    expect(isProvisionalDisplayName(`pat${PROVISIONAL_DISPLAY_SUFFIX}`)).toBe(true)
    expect(isProvisionalDisplayName('Pat')).toBe(false)
  })

  it('requires profile setup without avatar or provisional name', () => {
    expect(needsProfileSetup({ displayName: 'Pat', avatarKey: 'compass' })).toBe(false)
    expect(needsProfileSetup({ displayName: `pat${PROVISIONAL_DISPLAY_SUFFIX}` })).toBe(true)
    expect(needsProfileSetup({ displayName: 'Pat' })).toBe(true)
    expect(
      needsProfileSetup({
        displayName: `river${PROVISIONAL_DISPLAY_SUFFIX}`,
        avatarSource: 'SPIRIT_ANIMAL',
        avatarUrl: 'https://joinquest.cc/avatars/spirit/wolf.png',
      }),
    ).toBe(true)
    expect(
      needsProfileSetup({
        displayName: 'River',
        avatarSource: 'SPIRIT_ANIMAL',
        avatarUrl: 'https://joinquest.cc/avatars/spirit/wolf.png',
      }),
    ).toBe(false)
  })

  it('detects existing avatar from key, url, or spirit animal', () => {
    expect(hasExistingAvatar({ avatarKey: 'compass' })).toBe(true)
    expect(hasExistingAvatar({ avatarUrl: 'https://joinquest.cc/avatars/spirit/wolf.png' })).toBe(true)
    expect(hasExistingAvatar({ avatarSource: 'SPIRIT_ANIMAL' })).toBe(true)
    expect(hasExistingAvatar({ displayName: 'Pat' })).toBe(false)
  })

  it('prefills display name without provisional suffix', () => {
    expect(defaultDisplayNameInput({ displayName: `river${PROVISIONAL_DISPLAY_SUFFIX}` })).toBe('river')
    expect(defaultDisplayNameInput({ displayName: 'River' })).toBe('River')
  })

  it('resolves avatar url from key when url is missing', () => {
    expect(resolveUserAvatarUrl({ avatarKey: 'storm' })).toBe('/avatars/storm.png')
    expect(resolveUserAvatarUrl({ avatarUrl: 'https://joinquest.cc/avatars/beacon.png' })).toBe(
      'https://joinquest.cc/avatars/beacon.png',
    )
  })
})
