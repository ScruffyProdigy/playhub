import { describe, expect, it } from 'vitest'
import { parseRoomInviteCode, roomShareText } from './rooms'

describe('parseRoomInviteCode', () => {
  it('extracts code from room path', () => {
    expect(parseRoomInviteCode('/room/abc123')).toBe('ABC123')
    expect(parseRoomInviteCode('/room/X7K2M9/')).toBe('X7K2M9')
  })

  it('returns null for non-room paths', () => {
    expect(parseRoomInviteCode('/')).toBeNull()
    expect(parseRoomInviteCode('/return')).toBeNull()
  })
})

describe('roomShareText', () => {
  it('includes join url', () => {
    expect(roomShareText('https://joinquest.cc/room/ABC123')).toContain('https://joinquest.cc/room/ABC123')
  })
})
