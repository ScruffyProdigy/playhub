import { graphqlRequest } from './graphql'

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

export async function fetchCurrentUser() {
  const data = await graphqlRequest(ME_QUERY)
  return data.me
}

export async function requestMagicLink(email) {
  const data = await graphqlRequest(LOGIN_MAGIC_MUTATION, { email })
  return data.loginMagic === true
}

export async function completeMagicLogin(token) {
  const data = await graphqlRequest(COMPLETE_MAGIC_MUTATION, { token })
  return data.completeMagic
}
