import { graphqlRequest } from './graphql'

export const TABLE_FIELDS = `
  id
  createdAt
  canStart
  canDiscard
  game {
    id
    name
  }
  mode {
    id
    displayName
  }
  king {
    id
    displayName
  }
  seats {
    seatKey
    seatedAt
    user {
      id
      displayName
    }
  }
  seatSlots {
    seatKey
    queuePath
    displayName
    teamPrefix
    user {
      id
      displayName
    }
  }
  lookForGroupOptions {
    queueId
    queueName
    visible
    enabled
  }
`

export const MY_TABLE_SEAT_QUERY = `
  query MyTableSeat {
    myTableSeat {
      tableId
      roomId
      inviteCode
      gameId
      gameName
      modeName
      seatKey
      seatDisplayName
    }
  }
`

const CREATE_PRIVATE_TABLE = `
  mutation CreatePrivateTable($gameId: ID!, $modeId: ID!) {
    createPrivateTable(gameId: $gameId, modeId: $modeId) {
      ${TABLE_FIELDS}
    }
  }
`

const SIT_AT_TABLE = `
  mutation SitAtTable($tableId: ID!, $seatKey: String!) {
    sitAtTable(tableId: $tableId, seatKey: $seatKey) {
      ${TABLE_FIELDS}
    }
  }
`

const LEAVE_TABLE = `
  mutation LeaveTable($tableId: ID!) {
    leaveTable(tableId: $tableId)
  }
`

const DISCARD_TABLE = `
  mutation DiscardTable($tableId: ID!) {
    discardTable(tableId: $tableId)
  }
`

const START_TABLE = `
  mutation StartTable($tableId: ID!) {
    startTable(tableId: $tableId) {
      queued
      sessionId
      joinUrl
    }
  }
`

export async function fetchMyTableSeat() {
  const data = await graphqlRequest(MY_TABLE_SEAT_QUERY)
  return data.myTableSeat ?? null
}

export async function createPrivateTable(gameId, modeId) {
  const data = await graphqlRequest(CREATE_PRIVATE_TABLE, { gameId, modeId })
  return data.createPrivateTable
}

export async function sitAtTable(tableId, seatKey) {
  const data = await graphqlRequest(SIT_AT_TABLE, { tableId, seatKey })
  return data.sitAtTable
}

export async function leaveTable(tableId) {
  const data = await graphqlRequest(LEAVE_TABLE, { tableId })
  return data.leaveTable
}

export async function discardTable(tableId) {
  const data = await graphqlRequest(DISCARD_TABLE, { tableId })
  return data.discardTable
}

export async function startTable(tableId) {
  const data = await graphqlRequest(START_TABLE, { tableId })
  return data.startTable
}

/** Group seat slots by team prefix for column layout (e.g. Team-1 vs Team-2). */
export function groupSeatSlotsByTeam(seatSlots = []) {
  const teams = new Map()
  const ungrouped = []
  for (const slot of seatSlots) {
    const prefix = slot.teamPrefix?.trim()
    if (!prefix) {
      ungrouped.push(slot)
      continue
    }
    if (!teams.has(prefix)) {
      teams.set(prefix, [])
    }
    teams.get(prefix).push(slot)
  }
  const sortedTeams = [...teams.entries()].sort(([a], [b]) => a.localeCompare(b))
  return { teams: sortedTeams, ungrouped }
}

/**
 * Layout helper: team columns, role sections (Word Hunt), or flat fifo seats.
 */
export function groupSeatSlotsForDisplay(seatSlots = []) {
  const { teams, ungrouped } = groupSeatSlotsByTeam(seatSlots)
  if (teams.length > 0) {
    return { kind: 'teams', teams, ungrouped }
  }

  const byRole = new Map()
  const noPath = []
  for (const slot of seatSlots) {
    const path = slot.queuePath?.trim()
    if (!path) {
      noPath.push(slot)
      continue
    }
    if (!byRole.has(path)) {
      byRole.set(path, [])
    }
    byRole.get(path).push(slot)
  }

  if (byRole.size > 1) {
    const roles = [...byRole.entries()].sort(([a], [b]) => a.localeCompare(b))
    return {
      kind: 'roles',
      roles: roles.map(([path, slots]) => {
        const title = slots[0]?.displayName?.split(' · ')[0]?.trim() || path
        return [title, slots]
      }),
      noPath,
    }
  }

  return { kind: 'flat', slots: seatSlots.length ? seatSlots : noPath }
}

/**
 * Short seat label within a section header (e.g. "Red" under "Clue Giver").
 * Returns null when the section title carries enough context (numbered guesser seats).
 */
export function seatLabelInSection(slot, sectionTitle) {
  const full = slot?.displayName?.trim() || ''
  if (!full) {
    return null
  }

  const parts = full.split(' · ').map((part) => part.trim()).filter(Boolean)
  if (parts.length >= 2) {
    const role = parts[0]
    const suffix = parts.slice(1).join(' · ')
    if (sectionTitle && role.toLowerCase() === sectionTitle.trim().toLowerCase()) {
      if (/^\d+$/.test(suffix)) {
        return null
      }
      return suffix
    }
    if (/^\d+$/.test(suffix)) {
      return null
    }
    return suffix
  }

  if (/^\d+$/.test(full)) {
    return null
  }

  if (sectionTitle && full.toLowerCase() === sectionTitle.trim().toLowerCase()) {
    return null
  }

  return full
}

export function sectionTitleForSeat(seatKey, seatSlots = []) {
  const layout = groupSeatSlotsForDisplay(seatSlots)
  if (layout.kind === 'roles') {
    for (const [title, slots] of layout.roles) {
      if (slots.some((slot) => slot.seatKey === seatKey)) {
        return title
      }
    }
  }
  if (layout.kind === 'teams') {
    for (const [prefix, slots] of layout.teams) {
      if (slots.some((slot) => slot.seatKey === seatKey)) {
        return prefix.replace('-', ' ')
      }
    }
  }
  return null
}

export function displayName(user) {
  return user?.displayName?.trim() || 'Player'
}

/** Merge table.seats users into seatSlots when slot.user is missing (partial updates). */
export function enrichTableSeats(table) {
  if (!table) {
    return table
  }
  const usersByKey = Object.fromEntries(
    (table.seats ?? []).map((seat) => [seat.seatKey, seat.user ?? null]),
  )
  const seatSlots = (table.seatSlots ?? []).map((slot) => ({
    ...slot,
    user: slot.user ?? usersByKey[slot.seatKey] ?? null,
  }))
  return { ...table, seatSlots }
}

/** Merge a table subscription/mutation payload into prior room state. */
export function mergeTableRecord(prevTable, updatedTable) {
  if (!updatedTable?.id) {
    return prevTable
  }
  if (!prevTable) {
    return updatedTable
  }
  return {
    ...prevTable,
    ...updatedTable,
    ...(updatedTable.seats != null ? { seats: updatedTable.seats } : {}),
    ...(updatedTable.seatSlots != null ? { seatSlots: updatedTable.seatSlots } : {}),
  }
}

export function mySeatKeyOnTable(table, userId) {
  if (!userId) {
    return ''
  }
  for (const seat of table?.seats ?? []) {
    if (seat.user?.id === userId) {
      return seat.seatKey
    }
  }
  for (const slot of table?.seatSlots ?? []) {
    if (slot.user?.id === userId) {
      return slot.seatKey
    }
  }
  return ''
}

export function mySeatDisplayName(table, userId) {
  const seatKey = mySeatKeyOnTable(table, userId)
  if (!seatKey) {
    return ''
  }
  const slot = (table?.seatSlots ?? []).find((s) => s.seatKey === seatKey)
  if (!slot) {
    return seatKey
  }
  const sectionTitle = sectionTitleForSeat(seatKey, table?.seatSlots ?? [])
  const short = seatLabelInSection(slot, sectionTitle)
  if (short) {
    return sectionTitle ? `${sectionTitle} · ${short}` : short
  }
  return sectionTitle || slot.displayName?.trim() || seatKey
}

export function isKing(table, userId) {
  return Boolean(table?.king?.id && userId && table.king.id === userId)
}
