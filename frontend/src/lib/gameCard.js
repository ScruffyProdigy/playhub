/** Bump when replacing catalog/detail artwork so browsers skip stale caches. */
export const GAME_ARTWORK_VERSION = 1

function withArtworkVersion(url) {
  const trimmed = String(url || '').trim()
  if (!trimmed || trimmed.includes('?v=')) {
    return trimmed
  }
  return `${trimmed}?v=${GAME_ARTWORK_VERSION}`
}

/** Resolve a catalog game icon URL (required on every game). */
export function gameIconUrl(game) {
  const url = game?.iconUrl?.trim()
  if (!url) {
    throw new Error('game iconUrl is required')
  }
  return withArtworkVersion(url)
}

/** Resolve the catalog card hero (5:2). Uses catalogHeroUrl when set, else heroUrl. */
export function gameCatalogHeroUrl(game) {
  const url = game?.catalogHeroUrl?.trim() || game?.heroUrl?.trim()
  if (!url) {
    throw new Error('game catalog hero is required')
  }
  return withArtworkVersion(url)
}

/** Resolve a catalog game hero banner URL for the detail page (16:9 slot). */
export function gameHeroUrl(game) {
  const url = game?.heroUrl?.trim()
  if (!url) {
    throw new Error('game heroUrl is required')
  }
  return withArtworkVersion(url)
}

/** Shareable detail page path, or null when the game has no slug. */
export function gamePagePath(game) {
  const slug = game?.slug?.trim()
  if (!slug) {
    return null
  }
  return `/games/${slug}`
}

/** Absolute URL for sharing a game detail page. */
export function gamePageShareUrl(game) {
  const path = gamePagePath(game)
  if (!path) {
    return null
  }
  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin}${path}`
  }
  return path
}

/** Long description for the detail page, falling back to the catalog blurb. */
export function gameDetailDescription(game) {
  const long = game?.longDescription?.trim()
  if (long) {
    return long
  }
  return game?.shortDescription?.trim() || ''
}

/** Visible tag chips for catalog cards (max 3). */
export function gameTagChips(tags) {
  if (!Array.isArray(tags)) {
    return []
  }
  return tags.map((tag) => String(tag).trim()).filter(Boolean).slice(0, 3)
}
