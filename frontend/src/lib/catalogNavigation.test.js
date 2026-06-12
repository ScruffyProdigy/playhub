import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import {
  CATALOG_RETURN_KEY,
  CATALOG_SCROLL_KEY,
  navigateBackToCatalog,
  navigateToGameDetail,
  restoreCatalogScrollIfPending,
  saveCatalogScrollForReturn,
} from './catalogNavigation'
import * as pathname from './usePathname'

vi.mock('./usePathname', () => ({
  navigateTo: vi.fn(),
}))

describe('catalogNavigation', () => {
  beforeEach(() => {
    sessionStorage.clear()
    vi.mocked(pathname.navigateTo).mockReset()
    Object.defineProperty(window, 'scrollY', { value: 420, configurable: true, writable: true })
    window.scrollTo = vi.fn()
  })

  afterEach(() => {
    sessionStorage.clear()
  })

  it('saveCatalogScrollForReturn stores scroll position and return flag', () => {
    saveCatalogScrollForReturn()
    expect(sessionStorage.getItem(CATALOG_SCROLL_KEY)).toBe('420')
    expect(sessionStorage.getItem(CATALOG_RETURN_KEY)).toBe('1')
  })

  it('navigateToGameDetail saves scroll and navigates', () => {
    navigateToGameDetail('word-hunt')
    expect(sessionStorage.getItem(CATALOG_RETURN_KEY)).toBe('1')
    expect(pathname.navigateTo).toHaveBeenCalledWith('/games/word-hunt')
  })

  it('navigateBackToCatalog returns to the catalog route', () => {
    navigateBackToCatalog()
    expect(pathname.navigateTo).toHaveBeenCalledWith('/')
  })

  it('restoreCatalogScrollIfPending restores saved scroll once', async () => {
    saveCatalogScrollForReturn()
    restoreCatalogScrollIfPending()
    await new Promise((resolve) => requestAnimationFrame(resolve))

    expect(window.scrollTo).toHaveBeenCalledWith({ top: 420, left: 0 })
    expect(sessionStorage.getItem(CATALOG_RETURN_KEY)).toBeNull()
    expect(sessionStorage.getItem(CATALOG_SCROLL_KEY)).toBeNull()

    restoreCatalogScrollIfPending()
    expect(window.scrollTo).toHaveBeenCalledTimes(1)
  })

  it('restoreCatalogScrollIfPending is a no-op without a pending return', () => {
    restoreCatalogScrollIfPending()
    expect(window.scrollTo).not.toHaveBeenCalled()
  })
})
