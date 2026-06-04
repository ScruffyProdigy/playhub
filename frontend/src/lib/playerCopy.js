/** Player-facing strings (avoid “queue” in the shell UI). */

export const APP_TAGLINE = 'Find your group. Play together.'

export const LOOK_FOR_GROUP = 'Look for group'
export const LOOKING_FOR_GROUP = 'Looking…'
export const STOP_LOOKING = 'Stop looking'

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

export function bannerWaitingLine(gameName, count) {
  const n = count ?? 0
  return `Looking for a group in ${gameName} · ${n} ${n === 1 ? 'player' : 'players'} looking`
}

export function bannerMatchedLine(gameName) {
  return `Your group is ready — ${gameName}`
}

export function switchedFromGroupMessage(gameName) {
  return `You left the group for ${gameName} to look for a group here.`
}
