import { graphqlRequest } from './graphql'

const RETURN_DESTINATION_QUERY = `
  query ReturnDestination($matchId: ID) {
    returnDestination(matchId: $matchId) {
      path
      kind
    }
  }
`

export async function fetchReturnDestination(matchId) {
  const variables = matchId ? { matchId } : {}
  const data = await graphqlRequest(RETURN_DESTINATION_QUERY, variables)
  return data.returnDestination
}
