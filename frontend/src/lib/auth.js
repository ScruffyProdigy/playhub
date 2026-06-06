import { graphqlRequest } from './graphql'
import { clearSubscriptionAuthCache } from './queue'

const ME_QUERY = `
  query Me {
    me {
      id
      email
      displayName
      avatarUrl
      avatarKey
      avatarSource
      createdAt
    }
  }
`

const REQUEST_SIGN_IN_MUTATION = `
  mutation RequestSignIn($email: String!) {
    requestSignIn(email: $email)
  }
`

const COMPLETE_SIGN_IN_WITH_LINK_MUTATION = `
  mutation CompleteSignInWithLink($token: ID!) {
    completeSignInWithLink(token: $token) {
      id
      email
      displayName
      createdAt
      avatarUrl
      avatarKey
      avatarSource
    }
  }
`

const COMPLETE_SIGN_IN_WITH_CODE_MUTATION = `
  mutation CompleteSignInWithCode($email: String!, $code: String!) {
    completeSignInWithCode(email: $email, code: $code) {
      id
      email
      displayName
      avatarUrl
      avatarKey
      avatarSource
      createdAt
    }
  }
`

const LOGOUT_MUTATION = `
  mutation Logout {
    logout
  }
`

export async function fetchCurrentUser() {
  const data = await graphqlRequest(ME_QUERY)
  return data.me
}

export async function requestSignIn(email) {
  const data = await graphqlRequest(REQUEST_SIGN_IN_MUTATION, { email })
  return data.requestSignIn === true
}

export async function completeSignInWithCode(email, code) {
  const data = await graphqlRequest(COMPLETE_SIGN_IN_WITH_CODE_MUTATION, { email, code })
  clearSubscriptionAuthCache()
  return data.completeSignInWithCode
}

export async function logout() {
  const data = await graphqlRequest(LOGOUT_MUTATION)
  clearSubscriptionAuthCache()
  return data.logout === true
}

export async function completeSignInWithLink(token) {
  const data = await graphqlRequest(COMPLETE_SIGN_IN_WITH_LINK_MUTATION, { token })
  clearSubscriptionAuthCache()
  return data.completeSignInWithLink
}

const pendingCompletions = new Map()

// React StrictMode mounts twice in dev; share one in-flight completion per token.
export async function completeSignInWithLinkOnce(token) {
  const key = token.trim()
  if (!key) {
    throw new Error('Missing sign-in token')
  }

  if (pendingCompletions.has(key)) {
    return pendingCompletions.get(key)
  }

  const promise = completeSignInWithLink(key).catch((error) => {
    pendingCompletions.delete(key)
    throw error
  })
  pendingCompletions.set(key, promise)
  return promise
}
