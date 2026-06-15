import { readFileSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { describe, expect, it } from 'vitest'
import {
  REQUIRED_INTEGRATION_CHECKS,
  canRequestPublicRelease,
  integrationNextSteps,
  parseDeveloperRoute,
  suggestSlugFromName,
  visibilityLabel,
} from './developers'

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), '../../..')

function parseGoRequiredChecks(source) {
  const match = source.match(/requiredPassChecks = \[\]string\{([\s\S]*?)\}/)
  if (!match) throw new Error('requiredPassChecks block not found')
  return [...match[1].matchAll(/"([^"]+)"/g)].map((m) => m[1])
}

describe('buildMcpServerConfig', () => {
  it('uses mcpServers for Cursor', async () => {
    const { buildMcpServerConfig } = await import('./developers')
    const config = buildMcpServerConfig({ apiKey: 'lq_dev_test', clientId: 'cursor' })
    expect(config.mcpServers['joinquest-integration'].args).toContain('joinquest-integration-mcp-cursor')
  })

  it('uses servers for Copilot', async () => {
    const { buildMcpServerConfig } = await import('./developers')
    const config = buildMcpServerConfig({ apiKey: 'lq_dev_test', clientId: 'copilot' })
    expect(config.servers['joinquest-integration'].env.JOINQUEST_API_KEY).toBe('lq_dev_test')
    expect(config.mcpServers).toBeUndefined()
  })

  it('uses env interpolation for Windsurf', async () => {
    const { buildMcpServerConfig } = await import('./developers')
    const config = buildMcpServerConfig({ apiKey: 'lq_dev_test', clientId: 'windsurf' })
    expect(config.mcpServers['joinquest-integration'].env.JOINQUEST_API_KEY).toBe('${env:JOINQUEST_API_KEY}')
  })
})

describe('buildInstallDevCommand', () => {
  it('maps platform flags', async () => {
    const { buildInstallDevCommand } = await import('./developers')
    expect(buildInstallDevCommand({ apiKey: 'k', client: 'copilot' })).toContain('--copilot')
    expect(buildInstallDevCommand({ apiKey: 'k', client: 'roo' })).toContain('--roo')
  })
})

describe('REQUIRED_INTEGRATION_CHECKS', () => {
  it('matches backend requiredPassChecks', () => {
    const goSource = readFileSync(join(repoRoot, 'backend/internal/store/developer.go'), 'utf8')
    expect(REQUIRED_INTEGRATION_CHECKS).toEqual(parseGoRequiredChecks(goSource))
  })
})

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
    'manifest.launch_urls_on_provision',
    'manifest.game_modes',
    'manifest.sync_freshness',
    'provision.happy_path',
    'provision.idempotent_repush',
    'provision.auth',
    'provision.missing_auth',
    'provision.launch_urls',
    'provision.launch_url_no_jwt',
    'jwt.jwks',
    'jwt.claim_happy_path',
    'jwt.wrong_audience',
    'jwt.unknown_match',
    'jwt.wrong_issuer',
    'jwt.expired',
    'jwt.invalid_token',
    'jwt.wrong_seat',
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
