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
        modeKey
        displayName
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

export async function fetchGames() {
  const data = await graphqlRequest(GAMES_QUERY)
  return data.games ?? []
}
