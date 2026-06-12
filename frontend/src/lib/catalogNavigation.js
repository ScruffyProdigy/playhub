import { navigateTo } from './usePathname'

export const CATALOG_SCROLL_KEY = 'lobby.catalogScrollY'
export const CATALOG_RETURN_KEY = 'lobby.catalogReturnPending'

/** Remember catalog scroll position before opening a game detail page. */
export function saveCatalogScrollForReturn() {
  sessionStorage.setItem(CATALOG_SCROLL_KEY, String(window.scrollY))
  sessionStorage.setItem(CATALOG_RETURN_KEY, '1')
}

/** Restore catalog scroll when returning from a game detail page. */
export function restoreCatalogScrollIfPending() {
  if (sessionStorage.getItem(CATALOG_RETURN_KEY) !== '1') {
    return
  }

  sessionStorage.removeItem(CATALOG_RETURN_KEY)
  const y = Number(sessionStorage.getItem(CATALOG_SCROLL_KEY))
  sessionStorage.removeItem(CATALOG_SCROLL_KEY)

  if (!Number.isFinite(y) || y < 0) {
    return
  }

  requestAnimationFrame(() => {
    window.scrollTo({ top: y, left: 0 })
  })
}

export function navigateToGameDetail(slug) {
  const trimmed = String(slug || '').trim()
  if (!trimmed) {
    return
  }
  saveCatalogScrollForReturn()
  navigateTo(`/games/${encodeURIComponent(trimmed)}`)
}

export function navigateBackToCatalog() {
  navigateTo('/')
}
