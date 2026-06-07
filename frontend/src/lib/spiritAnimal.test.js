import { describe, expect, it } from 'vitest'
import {
  formatPhaseCountdown,
  formatSpiritAnimalJourneyCooldown,
  secondsRemainingFromPhase,
} from './spiritAnimal'

describe('spiritAnimal timings', () => {
  it('counts down from phase start and estimate', () => {
    const started = new Date(Date.now() - 15_000).toISOString()
    expect(secondsRemainingFromPhase(started, 45)).toBe(30)
  })

  it('formats remaining seconds for players', () => {
    expect(formatPhaseCountdown(42)).toBe('About 42 seconds left')
    expect(formatPhaseCountdown(75)).toBe('About 2 minutes left')
    expect(formatPhaseCountdown(0)).toBe('Almost there…')
  })

  it('formats journey cooldown for players', () => {
    expect(formatSpiritAnimalJourneyCooldown(12, null)).toBe(
      'Your next spirit animal journey opens in 12 days.',
    )
  })
})
