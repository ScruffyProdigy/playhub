import { render } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import useWaitForSignIn from './useWaitForSignIn'
import * as authBroadcast from '../../lib/authBroadcast'

vi.mock('../../lib/authBroadcast', () => ({
  onAuthComplete: vi.fn(() => () => {}),
}))

function TestHarness({ enabled, onSignedIn }) {
  useWaitForSignIn({ enabled, onSignedIn })
  return null
}

describe('useWaitForSignIn', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    vi.mocked(authBroadcast.onAuthComplete).mockReset()
    vi.mocked(authBroadcast.onAuthComplete).mockReturnValue(() => {})
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('does nothing when disabled', () => {
    const onSignedIn = vi.fn()
    render(<TestHarness enabled={false} onSignedIn={onSignedIn} />)

    expect(authBroadcast.onAuthComplete).not.toHaveBeenCalled()
    vi.advanceTimersByTime(5000)
    expect(onSignedIn).not.toHaveBeenCalled()
  })

  it('subscribes to auth broadcast when enabled', () => {
    const onSignedIn = vi.fn()
    render(<TestHarness enabled={true} onSignedIn={onSignedIn} />)

    expect(authBroadcast.onAuthComplete).toHaveBeenCalledWith(onSignedIn)
  })

  it('polls for session updates while enabled', () => {
    const onSignedIn = vi.fn()
    render(<TestHarness enabled={true} onSignedIn={onSignedIn} />)

    vi.advanceTimersByTime(2500)
    expect(onSignedIn).toHaveBeenCalledTimes(1)
    vi.advanceTimersByTime(2500)
    expect(onSignedIn).toHaveBeenCalledTimes(2)
  })

  it('cleans up broadcast subscription on unmount', () => {
    const unsubscribe = vi.fn()
    vi.mocked(authBroadcast.onAuthComplete).mockReturnValue(unsubscribe)

    const onSignedIn = vi.fn()
    const view = render(<TestHarness enabled={true} onSignedIn={onSignedIn} />)
    view.unmount()

    expect(unsubscribe).toHaveBeenCalled()
  })
})
