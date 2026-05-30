import { createClient } from 'graphql-ws'
import { getGraphQLWsUrl } from './env'
import { graphqlRequest } from './graphql'

const JOIN_GAME_MUTATION = `
  mutation JoinGame($gameId: ID!) {
    joinGame(gameId: $gameId) {
      queued
      queuedCount
      sessionId
      joinUrl
    }
  }
`

const MY_QUEUE_STATUS_QUERY = `
  query MyQueueStatus($gameId: ID!) {
    myQueueStatus(gameId: $gameId) {
      queued
      sessionId
      joinUrl
      queuedCount
    }
  }
`

const SUBSCRIPTION_AUTH_QUERY = `
  query SubscriptionAuth {
    subscriptionAuth
  }
`

const LEAVE_QUEUE_MUTATION = `
  mutation LeaveQueue($gameId: ID!) {
    leaveQueue(gameId: $gameId)
  }
`

const QUEUE_UPDATED_SUBSCRIPTION = `
  subscription QueueUpdated($gameId: ID!) {
    queueUpdated(gameId: $gameId) {
      gameId
      status
      sessionId
      joinUrl
      queuedCount
    }
  }
`

let cachedSubscriptionAuth = null
let subscriptionAuthPromise = null
let wsClient = null

async function loadSubscriptionAuth() {
  if (cachedSubscriptionAuth) {
    return cachedSubscriptionAuth
  }
  if (!subscriptionAuthPromise) {
    subscriptionAuthPromise = graphqlRequest(SUBSCRIPTION_AUTH_QUERY)
      .then((data) => {
        const header = data.subscriptionAuth?.trim() || null
        cachedSubscriptionAuth = header
        return header
      })
      .finally(() => {
        subscriptionAuthPromise = null
      })
  }
  return subscriptionAuthPromise
}

/** Ensure subscriptionAuth is loaded (call after sign-in, before subscribing). */
export async function prefetchSubscriptionAuth() {
  return loadSubscriptionAuth()
}

function getSubscriptionConnectionParams() {
  return async () => {
    const header = await loadSubscriptionAuth()
    if (!header) {
      throw new Error('Sign in required for live queue updates')
    }
    return { Authorization: header }
  }
}

function getWsClient() {
  if (!wsClient) {
    wsClient = createClient({
      url: getGraphQLWsUrl(),
      connectionParams: getSubscriptionConnectionParams(),
      retryAttempts: 10,
      retryWait: async (retries) => Math.min(500 * retries, 5000),
      shouldRetry: () => true,
      lazy: false,
    })
  }
  return wsClient
}

/** Clear cached WS auth after logout or session refresh. */
export function clearSubscriptionAuthCache() {
  cachedSubscriptionAuth = null
  subscriptionAuthPromise = null
  if (wsClient) {
    wsClient.dispose()
    wsClient = null
  }
}

export async function joinGame(gameId) {
  const data = await graphqlRequest(JOIN_GAME_MUTATION, { gameId })
  return data.joinGame
}

export async function fetchMyQueueStatus(gameId) {
  const data = await graphqlRequest(MY_QUEUE_STATUS_QUERY, { gameId })
  return data.myQueueStatus
}

export async function leaveQueue(gameId) {
  const data = await graphqlRequest(LEAVE_QUEUE_MUTATION, { gameId })
  return data.leaveQueue
}

function formatSubscriptionError(err) {
  if (!err) {
    return 'Queue updates unavailable'
  }
  if (typeof err === 'string') {
    return err
  }
  if (Array.isArray(err)) {
    return err[0]?.message || 'Queue updates unavailable'
  }
  if (err.message) {
    return err.message
  }
  return 'Queue updates unavailable'
}

export async function subscribeToQueue(gameId, { onUpdate, onError } = {}) {
  await prefetchSubscriptionAuth()

  const client = getWsClient()

  return client.subscribe(
    {
      query: QUEUE_UPDATED_SUBSCRIPTION,
      variables: { gameId },
    },
    {
      next: (payload) => {
        if (payload?.errors?.length) {
          onError?.(formatSubscriptionError(payload.errors))
          return
        }
        if (payload?.data?.queueUpdated) {
          onUpdate?.(payload.data.queueUpdated)
        }
      },
      error: (err) => onError?.(formatSubscriptionError(err)),
      complete: () => {},
    },
  )
}
