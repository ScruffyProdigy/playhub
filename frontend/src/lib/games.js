import { graphqlRequest } from './graphql'

const GAMES_QUERY = `
  query Games {
    games {
      id
      name
      createdAt
      activeSessions {
        id
      }
      modes {
        id
        modeKey
        displayName
        status
        queuePaths {
          queuePath
          displayName
          playersToStart
        }
        seats {
          queuePath
        }
        queues {
          id
          name
          playersToStart
          status
        }
      }
    }
  }
`

/** Pick the first active queue for a game (typically the default queue). */
export function defaultQueueForGame(game) {
  for (const mode of game?.modes ?? []) {
    for (const queue of mode?.queues ?? []) {
      if (queue.status === 'active') {
        return queue
      }
    }
  }
  return null
}

/** Mode that owns the default active queue for this game. */
export function defaultModeForGame(game) {
  for (const mode of game?.modes ?? []) {
    for (const queue of mode?.queues ?? []) {
      if (queue.status === 'active') {
        return mode
      }
    }
  }
  return null
}

/**
 * Join UI options derived from expanded seat queue paths.
 * @returns {{ kind: 'fifo', paths: [] } | { kind: 'composition', paths: string[] }}
 */
export function joinGroupOptionsForMode(mode) {
  const queuePaths = (mode?.queuePaths ?? []).filter(
    (entry) => (entry.queuePath?.trim() ?? '') !== '',
  )
  if (queuePaths.length > 0) {
    return {
      kind: 'composition',
      paths: queuePaths.map((entry) => ({
        queuePath: entry.queuePath,
        displayName: entry.displayName || entry.queuePath,
      })),
    }
  }

  const paths = new Set()
  for (const seat of mode?.seats ?? []) {
    const path = seat.queuePath?.trim() ?? ''
    if (path) {
      paths.add(path)
    }
  }
  const sorted = [...paths].sort()
  if (sorted.length === 0) {
    return { kind: 'fifo', paths: [] }
  }
  return {
    kind: 'composition',
    paths: sorted.map((queuePath) => ({ queuePath, displayName: queuePath })),
  }
}

export function joinGroupOptionsForGame(game) {
  return joinGroupOptionsForMode(defaultModeForGame(game))
}

export async function fetchGames() {
  const data = await graphqlRequest(GAMES_QUERY)
  return data.games ?? []
}
