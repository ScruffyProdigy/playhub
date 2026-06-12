import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { onTabVisible } from './tabVisibility'

describe('onTabVisible', () => {
  beforeEach(() => {
    Object.defineProperty(document, 'visibilityState', {
      configurable: true,
      value: 'visible',
    })
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('calls callback when tab becomes visible', () => {
    const callback = vi.fn()
    onTabVisible(callback)

    Object.defineProperty(document, 'visibilityState', { configurable: true, value: 'hidden' })
    document.dispatchEvent(new Event('visibilitychange'))
    expect(callback).not.toHaveBeenCalled()

    Object.defineProperty(document, 'visibilityState', { configurable: true, value: 'visible' })
    document.dispatchEvent(new Event('visibilitychange'))
    expect(callback).toHaveBeenCalledTimes(1)
  })
})
