import { describe, expect, it } from 'vitest'
import { parseDeveloperRoute, suggestSlugFromName, visibilityLabel } from './developers'

describe('parseDeveloperRoute', () => {
  it('parses landing', () => {
    expect(parseDeveloperRoute('/developers')).toEqual({ kind: 'landing' })
  })

  it('parses dashboard and welcome', () => {
    expect(parseDeveloperRoute('/developers/games/abc-123')).toEqual({
      kind: 'dashboard',
      gameId: 'abc-123',
    })
    expect(parseDeveloperRoute('/developers/games/abc-123/welcome')).toEqual({
      kind: 'welcome',
      gameId: 'abc-123',
    })
  })
})

describe('suggestSlugFromName', () => {
  it('slugifies game names', () => {
    expect(suggestSlugFromName('My Cool Game!')).toBe('my-cool-game')
  })
})

describe('visibilityLabel', () => {
  it('maps visibility enums', () => {
    expect(visibilityLabel('PRIVATE_TESTING')).toBe('Private testing')
    expect(visibilityLabel('DRAFT')).toBe('Draft')
  })
})
