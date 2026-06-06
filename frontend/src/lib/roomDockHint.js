const MEMBER_KEY = 'joinquest-room-member'
const CODE_KEY = 'joinquest-room-code'

export function readRoomMemberHint() {
  try {
    return sessionStorage.getItem(MEMBER_KEY) === '1'
  } catch {
    return false
  }
}

export function readRoomInviteHint() {
  try {
    return sessionStorage.getItem(CODE_KEY)?.trim().toUpperCase() || ''
  } catch {
    return ''
  }
}

export function writeRoomDockHint({ member, inviteCode = '' } = {}) {
  try {
    if (member) {
      sessionStorage.setItem(MEMBER_KEY, '1')
      const code = inviteCode.trim().toUpperCase()
      if (code) {
        sessionStorage.setItem(CODE_KEY, code)
      }
    } else {
      sessionStorage.removeItem(MEMBER_KEY)
      sessionStorage.removeItem(CODE_KEY)
    }
  } catch {
    // ignore storage errors (private mode, etc.)
  }
}
