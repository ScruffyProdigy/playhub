import { graphqlRequest } from './graphql'

const MY_ACTIVE_QUEUE_QUERY = `
  query MyActiveQueue {
    myActiveQueue {
      queueId
      gameId
      gameName
      status
      queuedCount
      joinUrl
    }
  }
`

export async function fetchMyActiveQueue() {
  const data = await graphqlRequest(MY_ACTIVE_QUEUE_QUERY)
  return data.myActiveQueue ?? null
}
