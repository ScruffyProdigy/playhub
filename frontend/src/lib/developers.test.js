import { describe, expect, it } from 'vitest'
import {
  canRequestPublicRelease,
  integrationNextSteps,
  parseDeveloperRoute,
  suggestSlugFromName,
  visibilityLabel,
} from './developers'

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

describe('integrationNextSteps', () => {
  it('marks connect as current for draft games', () => {
    const steps = integrationNextSteps({ visibility: 'DRAFT', integrationChecks: [] })
    expect(steps[0].status).toBe('current')
    expect(steps[0].id).toBe('connect')
  })

  it('marks checks current after API connect', () => {
    const steps = integrationNextSteps({
      visibility: 'PRIVATE_TESTING',
      integrationChecks: [],
    })
    expect(steps[0].status).toBe('done')
    expect(steps[1].status).toBe('current')
  })
})

describe('canRequestPublicRelease', () => {
  const passChecks = [
    'manifest.reach_api',
    'manifest.status',
    'manifest.game_modes',
    'manifest.sync_freshness',
    'provision.happy_path',
    'provision.auth',
    'provision.launch_urls',
  ].map((checkId) => ({ checkId, status: 'PASS' }))

  it('requires metadata and passing checks', () => {
    expect(
      canRequestPublicRelease({
        visibility: 'PRIVATE_TESTING',
        shortDescription: 'Short',
        longDescription: 'Long',
        tags: ['party'],
        integrationChecks: passChecks,
      }),
    ).toBe(true)
    expect(
      canRequestPublicRelease({
        visibility: 'PRIVATE_TESTING',
        shortDescription: 'Short',
        longDescription: '',
        tags: ['party'],
        integrationChecks: passChecks,
      }),
    ).toBe(false)
  })
})
