import { describe, it, expect, beforeEach, afterEach } from 'vitest'
import { readRoomInviteHint, readRoomMemberHint, writeRoomDockHint } from './roomDockHint'

describe('roomDockHint', () => {
  beforeEach(() => {
    sessionStorage.clear()
  })

  afterEach(() => {
    sessionStorage.clear()
  })

  it('starts false when unset', () => {
    expect(readRoomMemberHint()).toBe(false)
  })

  it('persists membership across reads', () => {
    writeRoomDockHint({ member: true, inviteCode: 'abc123' })
    expect(readRoomMemberHint()).toBe(true)
    expect(readRoomInviteHint()).toBe('ABC123')
  })

  it('clears on write false', () => {
    writeRoomDockHint({ member: true, inviteCode: 'abc123' })
    writeRoomDockHint({ member: false })
    expect(readRoomMemberHint()).toBe(false)
    expect(readRoomInviteHint()).toBe('')
  })
})
