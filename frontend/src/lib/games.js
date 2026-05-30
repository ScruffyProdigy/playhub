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
    }
  }
`

export async function fetchGames() {
  const data = await graphqlRequest(GAMES_QUERY)
  return data.games ?? []
}
