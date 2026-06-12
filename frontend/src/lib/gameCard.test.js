import { describe, it, expect } from 'vitest'
import {
  gameCatalogHeroUrl,
  gameDetailDescription,
  gameHeroUrl,
  gameIconUrl,
  gamePagePath,
  gamePageShareUrl,
  gameTagChips,
} from './gameCard'

describe('gameCard', () => {
  it('gameIconUrl returns the icon path', () => {
    expect(gameIconUrl({ iconUrl: '/games/rpslr-icon.png' })).toBe('/games/rpslr-icon.png?v=1')
  })

  it('gameIconUrl throws when missing', () => {
    expect(() => gameIconUrl({})).toThrow('game iconUrl is required')
  })

  it('gameHeroUrl returns the hero path', () => {
    expect(gameHeroUrl({ heroUrl: '/games/rpslr-hero.jpg' })).toBe('/games/rpslr-hero.jpg?v=1')
  })

  it('gameHeroUrl throws when missing', () => {
    expect(() => gameHeroUrl({})).toThrow('game heroUrl is required')
  })

  it('gameCatalogHeroUrl prefers catalogHeroUrl', () => {
    expect(
      gameCatalogHeroUrl({
        catalogHeroUrl: '/games/catalog.png',
        heroUrl: '/games/rpslr-hero.jpg',
      }),
    ).toBe('/games/catalog.png?v=1')
  })

  it('gameCatalogHeroUrl falls back to heroUrl', () => {
    expect(gameCatalogHeroUrl({ heroUrl: '/games/rpslr-hero.jpg' })).toBe('/games/rpslr-hero.jpg?v=1')
  })

  it('gamePagePath returns a shareable path', () => {
    expect(gamePagePath({ slug: 'word-hunt' })).toBe('/games/word-hunt')
  })

  it('gamePagePath returns null without a slug', () => {
    expect(gamePagePath({})).toBeNull()
  })

  it('gamePageShareUrl returns an absolute URL', () => {
    expect(gamePageShareUrl({ slug: 'word-hunt' })).toMatch(/\/games\/word-hunt$/)
  })

  it('gameDetailDescription prefers longDescription', () => {
    expect(
      gameDetailDescription({
        longDescription: 'Long copy.',
        shortDescription: 'Short copy.',
      }),
    ).toBe('Long copy.')
  })

  it('gameDetailDescription falls back to shortDescription', () => {
    expect(gameDetailDescription({ shortDescription: 'Short copy.' })).toBe('Short copy.')
  })

  it('gameTagChips caps at three tags', () => {
    expect(gameTagChips(['a', 'b', 'c', 'd'])).toEqual(['a', 'b', 'c'])
  })
})
