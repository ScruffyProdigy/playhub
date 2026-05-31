import { createClient } from 'graphql-ws'
import { getGraphQLWsUrl } from './env'
import { graphqlRequest } from './graphql'

const JOIN_QUEUE_MUTATION = `
  mutation JoinQueue($queueId: ID!) {
    joinQueue(queueId: $queueId) {
      queued
      queuedCount
      sessionId
      joinUrl
    }
  }
`

const MY_QUEUE_STATUS_QUERY = `
  query MyQueueStatus($queueId: ID!) {
    myQueueStatus(queueId: $queueId) {
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
  mutation LeaveQueue($queueId: ID!) {
    leaveQueue(queueId: $queueId)
  }
`

const QUEUE_UPDATED_SUBSCRIPTION = `
  subscription QueueUpdated($queueId: ID!) {
    queueUpdated(queueId: $queueId) {
      gameId
      queueId
      status
      sessionId
      joinUrl
      queuedCount
      message
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

export async function joinQueue(queueId) {
  const data = await graphqlRequest(JOIN_QUEUE_MUTATION, { queueId })
  return data.joinQueue
}

export async function fetchMyQueueStatus(queueId) {
  const data = await graphqlRequest(MY_QUEUE_STATUS_QUERY, { queueId })
  return data.myQueueStatus
}

export async function leaveQueue(queueId) {
  const data = await graphqlRequest(LEAVE_QUEUE_MUTATION, { queueId })
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

export async function subscribeToQueue(queueId, { onUpdate, onError } = {}) {
  await prefetchSubscriptionAuth()

  const client = getWsClient()

  return client.subscribe(
    {
      query: QUEUE_UPDATED_SUBSCRIPTION,
      variables: { queueId },
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
