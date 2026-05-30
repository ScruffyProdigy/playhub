import { graphqlRequest } from './graphql'
import { clearSubscriptionAuthCache } from './queue'

const ME_QUERY = `
  query Me {
    me {
      id
      email
      displayName
      createdAt
    }
  }
`

const LOGIN_MAGIC_MUTATION = `
  mutation LoginMagic($email: String!) {
    loginMagic(email: $email)
  }
`

const COMPLETE_MAGIC_MUTATION = `
  mutation CompleteMagic($token: ID!) {
    completeMagic(token: $token) {
      id
      email
      displayName
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

export async function requestMagicLink(email) {
  const data = await graphqlRequest(LOGIN_MAGIC_MUTATION, { email })
  return data.loginMagic === true
}

export async function logout() {
  const data = await graphqlRequest(LOGOUT_MUTATION)
  clearSubscriptionAuthCache()
  return data.logout === true
}

export async function completeMagicLogin(token) {
  const data = await graphqlRequest(COMPLETE_MAGIC_MUTATION, { token })
  clearSubscriptionAuthCache()
  return data.completeMagic
}

const pendingCompletions = new Map()

// React StrictMode mounts twice in dev; share one in-flight completion per token.
export async function completeMagicLoginOnce(token) {
  const key = token.trim()
  if (!key) {
    throw new Error('Missing sign-in token')
  }

  if (pendingCompletions.has(key)) {
    return pendingCompletions.get(key)
  }

  const promise = completeMagicLogin(key).catch((error) => {
    pendingCompletions.delete(key)
    throw error
  })
  pendingCompletions.set(key, promise)
  return promise
}
