import { graphqlRequest } from './graphql'
import { defaultModeForGame } from './games'

const MY_GAME_FIELDS = `
  id
  slug
  name
  shortDescription
  apiBaseUrl
  visibility
  contactEmail
  websiteUrl
  communityUrl
  manifestSyncedAt
  integrationChecks {
    checkId
    status
    message
    detail
    ranAt
  }
  modes {
    id
    modeKey
    displayName
    status
    queues {
      id
      status
    }
  }
`

const MY_GAMES_QUERY = `
  query MyGames {
    myGames {
      id
      slug
      name
      shortDescription
      visibility
    }
  }
`

const MY_GAME_QUERY = `
  query MyGame($id: ID!) {
    myGame(id: $id) {
      ${MY_GAME_FIELDS}
    }
  }
`

const MY_GAME_CREDENTIALS_QUERY = `
  query MyGameCredentials($id: ID!) {
    myGameCredentials(id: $id) {
      serviceToken
      webhookSecret
    }
  }
`

const REGISTER_MY_GAME = `
  mutation RegisterMyGame($input: RegisterMyGameInput!) {
    registerMyGame(input: $input) {
      connected
      connectError
      webhookSecret
      serviceToken
      game {
        id
        slug
        name
        visibility
      }
    }
  }
`

const RUN_MY_GAME_CHECKS = `
  mutation RunMyGameChecks($gameId: ID!) {
    runMyGameChecks(gameId: $gameId) {
      checkId
      status
      message
      detail
      ranAt
    }
  }
`

const INTEGRATION_GUIDE_QUERY = `
  query DeveloperIntegrationGuide {
    developerIntegrationGuide
  }
`

/** Parse developer routes from pathname. */
export function parseDeveloperRoute(pathname) {
  const path = String(pathname || '')
  if (path === '/developers' || path === '/developers/') {
    return { kind: 'landing' }
  }
  const welcomeMatch = path.match(/^\/developers\/games\/([^/]+)\/welcome\/?$/)
  if (welcomeMatch?.[1]) {
    return { kind: 'welcome', gameId: decodeURIComponent(welcomeMatch[1]) }
  }
  const dashboardMatch = path.match(/^\/developers\/games\/([^/]+)\/?$/)
  if (dashboardMatch?.[1]) {
    return { kind: 'dashboard', gameId: decodeURIComponent(dashboardMatch[1]) }
  }
  return null
}

export function suggestSlugFromName(name) {
  return String(name || '')
    .trim()
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 64)
}

export function visibilityLabel(visibility) {
  switch (visibility) {
    case 'PRIVATE_TESTING':
      return 'Private testing'
    case 'PENDING_REVIEW':
      return 'Pending review'
    case 'PUBLIC':
      return 'Public'
    default:
      return 'Draft'
  }
}

export function checkSectionTitle(checkId) {
  const section = String(checkId || '').split('.')[0]
  switch (section) {
    case 'manifest':
      return 'Manifest'
    case 'provision':
      return 'Provisioning'
    case 'jwt':
      return 'JWT verification'
    default:
      return 'Checks'
  }
}

export function checkStatusLabel(status) {
  switch (status) {
    case 'PASS':
      return 'Pass'
    case 'FAIL':
      return 'Fail'
    default:
      return 'Skipped'
  }
}

export async function fetchMyGames() {
  const data = await graphqlRequest(MY_GAMES_QUERY)
  return data.myGames ?? []
}

export async function fetchMyGame(gameId) {
  const data = await graphqlRequest(MY_GAME_QUERY, { id: gameId })
  return data.myGame ?? null
}

export async function fetchMyGameCredentials(gameId) {
  const data = await graphqlRequest(MY_GAME_CREDENTIALS_QUERY, { id: gameId })
  return data.myGameCredentials ?? null
}

export async function registerMyGame(input) {
  const data = await graphqlRequest(REGISTER_MY_GAME, { input })
  return data.registerMyGame
}

export async function runMyGameChecks(gameId) {
  const data = await graphqlRequest(RUN_MY_GAME_CHECKS, { gameId })
  return data.runMyGameChecks ?? []
}

export async function fetchDeveloperIntegrationGuide() {
  const data = await graphqlRequest(INTEGRATION_GUIDE_QUERY)
  return data.developerIntegrationGuide ?? ''
}

export function defaultModeForMyGame(game) {
  return defaultModeForGame(game)
}

export function developerDashboardPath(gameId) {
  return `/developers/games/${encodeURIComponent(gameId)}`
}

export function developerWelcomePath(gameId) {
  return `/developers/games/${encodeURIComponent(gameId)}/welcome`
}
