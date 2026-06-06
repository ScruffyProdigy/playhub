/** Player-facing strings (avoid “queue” in the shell UI). */

export const APP_TAGLINE = 'Find your group. Play together.'

export const LOOK_FOR_GROUP = 'Look for group'
export const LOOKING_FOR_GROUP = 'Looking…'
export const STOP_LOOKING = 'Stop looking'

export function joinAsLabel(queuePath) {
  return `Join as ${queuePath}`
}

export function waitingAsRoleLine(queuePath) {
  return `Looking as ${queuePath}…`
}

export const LAUNCH_GAME = 'Launch game'
export const LEAVE_MATCH = 'Leave match'

export const GAMES_HEADING = 'Available games'
export const GAMES_INTRO =
  'Pick a game and look for a group. Starting a new search moves you out of any other group you were waiting for.'

export function waitingForGroupLine(count) {
  const n = count ?? 0
  return `Looking for players… (${n} ${n === 1 ? 'player' : 'players'} looking)`
}

export function activeMatchBlockedLine(gameName) {
  return `You have an active match in ${gameName}. Finish or leave it before looking for another group.`
}

export function bannerWaitingLine(gameName, count, queuePathDisplayName) {
  const n = count ?? 0
  const players = `${n} ${n === 1 ? 'player' : 'players'} looking`
  const role = queuePathDisplayName?.trim()
  if (role) {
    return `Looking for a group in ${gameName} as ${role} · ${players}`
  }
  return `Looking for a group in ${gameName} · ${players}`
}

export function bannerMatchedLine(gameName) {
  return `Your group is ready — ${gameName}`
}

export const CREATE_PRIVATE_GAME = 'Create private game'
export const START_GAME = 'Start now'
export const DISCARD = 'Discard'
export const KING_LABEL = 'King'
export const LEAVE_TABLE_SEAT = 'Leave seat'

export function bannerTableSeatHint() {
  return 'Use Room below to manage seats or start the game.'
}

export function bannerTableSeatLine(gameName, modeName, seatDisplayName) {
  const role = seatDisplayName?.trim()
  if (role) {
    return `Seated at ${gameName} (${modeName}) as ${role}`
  }
  return `Seated at ${gameName} (${modeName})`
}

export function switchedFromGroupMessage(gameName) {
  return `You left the group for ${gameName} to look for a group here.`
}
