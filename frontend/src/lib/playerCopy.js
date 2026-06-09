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
export const LEAVE_GAME = 'Leave game'
export const LEAVE_MATCH = 'Leave match'

export const GAMES_HEADING = 'Available games'
export const GAMES_INTRO =
  'Pick a game and look for a group. Starting a new search moves you out of any other group you were waiting for.'

export function waitingForGroupLine(count) {
  const n = count ?? 0
  return `Looking for players… (${n} ${n === 1 ? 'player' : 'players'} looking)`
}

export function bannerWaitingLine(gameName, count, queuePathDisplayName, formingGaps) {
  const n = count ?? 0
  const players = `${n} ${n === 1 ? 'player' : 'players'} looking`
  const role = queuePathDisplayName?.trim()
  let line
  if (role) {
    line = `Looking for a group in ${gameName} as ${role} · ${players}`
  } else {
    line = `Looking for a group in ${gameName} · ${players}`
  }
  const needLine = formatFormingGapsNeedLine(formingGaps)
  if (needLine) {
    line += ` · ${needLine}`
  }
  return line
}

export function bannerIntentPlayingLine(gameName, modeName, roleLabel) {
  const mode = modeName?.trim()
  const role = roleLabel?.trim()
  if (mode && role) {
    return `Playing ${gameName} (${mode}) as ${role}`
  }
  if (mode) {
    return `Playing ${gameName} (${mode})`
  }
  return `Playing ${gameName}`
}

export function bannerIntentPlayingHint() {
  return 'Launch when you are ready — your seat is reserved.'
}

export function bannerIntentLaunchPendingHint() {
  return 'Preparing your launch link…'
}

export function bannerIntentWaitingHint() {
  return 'We will notify you here when your group is ready.'
}

export const CREATE_PRIVATE_GAME = 'Create private game'
export const START_GAME = 'Start now'
export const DISCARD = 'Discard'
export const KING_LABEL = 'King'
export const LEAVE_TABLE_SEAT = 'Leave seat'

export function bannerTableSeatHint() {
  return 'Use Room below to manage seats or start the game.'
}

export function bannerTableBackfillHint(formingGaps) {
  const needLine = formatFormingGapsNeedLine(formingGaps)
  if (needLine) {
    return `Your table is queued to start — ${needLine.charAt(0).toLowerCase()}${needLine.slice(1)}`
  }
  return 'Your table is queued to start — waiting for players from the lobby.'
}

/** e.g. "Need 1 Clue Giver, 3 Guessers" */
export function formatFormingGapsNeedLine(formingGaps) {
  const parts = formatFormingGapParts(formingGaps)
  if (parts.length === 0) {
    return null
  }
  return `Need ${parts.join(', ')}`
}

/** e.g. "Need 1 Clue Giver, 3 Guessers from the lobby" */
export function formatFormingGapsFromLobbyLine(formingGaps) {
  const parts = formatFormingGapParts(formingGaps)
  if (parts.length === 0) {
    return null
  }
  return `Need ${parts.join(', ')} from the lobby`
}

function formatFormingGapParts(formingGaps) {
  if (!Array.isArray(formingGaps)) {
    return []
  }
  return formingGaps
    .filter((gap) => (gap?.needed ?? 0) > 0)
    .map((gap) => {
      const count = gap.needed
      const label = gap.displayName?.trim() || gap.queuePath?.trim() || 'player'
      return `${count} ${label}`
    })
}

export function bannerTableStartedLine(gameName, modeName) {
  return `Your game is ready — ${gameName} (${modeName})`
}

export function bannerTableStartedHint() {
  return 'Launch when you are ready — your seat is reserved.'
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
