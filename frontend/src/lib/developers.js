import { graphqlRequest } from './graphql'
import { defaultModeForGame } from './games'

const MY_GAME_FIELDS = `
  id
  slug
  name
  shortDescription
  longDescription
  howToPlay
  tags
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
    seats {
      seatKey
      queuePath
    }
    queues {
      id
      status
    }
  }
`

const CATALOG_TAG_TAXONOMY_QUERY = `
  query CatalogTagTaxonomy {
    catalogTagTaxonomy {
      id
      label
      description
    }
  }
`

const UPDATE_MY_GAME_METADATA = `
  mutation UpdateMyGameMetadata($input: UpdateMyGameMetadataInput!) {
    updateMyGameMetadata(input: $input) {
      ${MY_GAME_FIELDS}
    }
  }
`

const REQUEST_PUBLIC_RELEASE = `
  mutation RequestPublicRelease($gameId: ID!) {
    requestPublicRelease(gameId: $gameId) {
      id
      visibility
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

const MY_DEVELOPER_API_KEYS_QUERY = `
  query MyDeveloperApiKeys {
    myDeveloperApiKeys {
      id
      name
      keyPrefix
      createdAt
      lastUsedAt
    }
  }
`

const CREATE_DEVELOPER_API_KEY = `
  mutation CreateDeveloperApiKey($name: String) {
    createDeveloperApiKey(name: $name) {
      secret
      apiKey {
        id
        name
        keyPrefix
        createdAt
        lastUsedAt
      }
    }
  }
`

const REVOKE_DEVELOPER_API_KEY = `
  mutation RevokeDeveloperApiKey($id: ID!) {
    revokeDeveloperApiKey(id: $id)
  }
`

export const REQUIRED_INTEGRATION_CHECKS = [
  'manifest.reach_api',
  'manifest.status',
  'manifest.launch_urls_on_provision',
  'manifest.game_modes',
  'manifest.sync_freshness',
  'provision.happy_path',
  'provision.idempotent_repush',
  'provision.auth',
  'provision.missing_auth',
  'provision.launch_urls',
  'provision.launch_url_no_jwt',
  'jwt.jwks',
  'jwt.claim_happy_path',
  'jwt.wrong_audience',
  'jwt.unknown_match',
  'jwt.wrong_issuer',
  'jwt.expired',
  'jwt.invalid_token',
  'jwt.wrong_seat',
]

export const CHECK_FIX_HINTS = {
  'manifest.reach_api':
    'JoinQuest must reach your public HTTPS API. Check the URL is live and not localhost.',
  'manifest.status':
    'GET /api/v1/status should return game info and launchUrlsOnProvision: true.',
  'manifest.launch_urls_on_provision':
    'GET /api/v1/status must include launchUrlsOnProvision: true.',
  'manifest.game_modes':
    'GET /api/v1/game-modes needs valid seatTemplate JSON for each mode.',
  'manifest.sync_freshness':
    'Re-sync your manifest — run checks again after deploying API changes.',
  'provision.happy_path':
    'POST /api/v1/matches should return launch URLs for each seated player.',
  'provision.idempotent_repush':
    'Re-posting the same externalMatchId should succeed (idempotent provision).',
  'provision.auth':
    'Accept the service token on Authorization: Bearer … when provisioning matches.',
  'provision.missing_auth':
    'Reject provision requests with no Authorization header.',
  'provision.banlist':
    'Return HTTP 403 with bannedLobbyUserIds when a player is banned.',
  'provision.launch_urls':
    'Every seated player needs a launchUrls entry (or a launchUrlTemplate).',
  'provision.launch_url_no_jwt':
    'Launch URL bases must not include a JWT — JoinQuest adds token= later.',
  'jwt.jwks': 'Publish JWKS at {lobby}/.well-known/jwks.json.',
  'jwt.claim_happy_path': 'POST /api/v1/matches/{id}/claim must accept Lobby seat JWTs.',
  'jwt.wrong_audience': 'Reject JWTs whose aud does not match your API base URL.',
  'jwt.unknown_match': 'Return 404 when the claim URL match id is unknown or mismatched.',
  'jwt.wrong_issuer': 'Reject JWTs whose iss does not match the match lobbyId.',
  'jwt.expired': 'Reject expired seat tokens with 401/403.',
  'jwt.invalid_token': 'Reject malformed tokens with 401/403.',
  'jwt.wrong_seat': 'Reject tokens that claim another player\'s reserved seat.',
}

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

export async function fetchCatalogTagTaxonomy() {
  const data = await graphqlRequest(CATALOG_TAG_TAXONOMY_QUERY)
  return data.catalogTagTaxonomy ?? []
}

export async function updateMyGameMetadata(input) {
  const data = await graphqlRequest(UPDATE_MY_GAME_METADATA, { input })
  return data.updateMyGameMetadata
}

export async function requestPublicRelease(gameId) {
  const data = await graphqlRequest(REQUEST_PUBLIC_RELEASE, { gameId })
  return data.requestPublicRelease
}

export function checkFixHint(checkId) {
  return CHECK_FIX_HINTS[checkId] ?? 'See the integration guide below for details.'
}

export function integrationNextSteps(game) {
  if (!game) {
    return []
  }

  const checksById = new Map((game.integrationChecks ?? []).map((c) => [c.checkId, c]))
  const hasChecks = (game.integrationChecks ?? []).length > 0
  const requiredPass = REQUIRED_INTEGRATION_CHECKS.every((id) => checksById.get(id) === 'PASS')
  const hasFailures = (game.integrationChecks ?? []).some((c) => c.status === 'FAIL')
  const connected = game.visibility !== 'DRAFT'
  const hasMetadata =
    Boolean(game.shortDescription?.trim()) &&
    Boolean(game.longDescription?.trim()) &&
    Array.isArray(game.tags) &&
    game.tags.length > 0
  const canRelease = canRequestPublicRelease(game)

  const steps = [
    {
      id: 'connect',
      label: 'Connect your game API',
      done: connected,
      hint: connected
        ? 'Your API is reachable and modes are synced.'
        : 'Deploy a public HTTPS URL with /healthz and /api/v1/game-modes — localhost will not work.',
    },
    {
      id: 'checks',
      label: 'Run integration checks',
      done: hasChecks && !hasFailures && requiredPass,
      hint: !connected
        ? 'Connect your API first.'
        : !hasChecks
          ? 'Click Run all checks on this dashboard.'
          : hasFailures
            ? 'Fix the failing checks below, then run checks again.'
            : 'All required checks passed.',
    },
    {
      id: 'metadata',
      label: 'Complete catalog listing',
      done: hasMetadata,
      hint: hasMetadata
        ? 'Short description, long description, and tags are set.'
        : 'Fill in catalog copy so players know what your game is about.',
    },
    {
      id: 'test',
      label: 'Create a test table',
      done: false,
      hint: 'Invite friends to your room and play before going public.',
      optional: !connected,
    },
    {
      id: 'release',
      label: 'Request public release',
      done: game.visibility === 'PENDING_REVIEW' || game.visibility === 'PUBLIC',
      hint: canRelease
        ? 'Ready when you are — we do a quick review before catalog listing.'
        : 'Complete checks and catalog metadata first.',
    },
  ]

  let foundCurrent = false
  return steps.map((step) => {
    if (step.optional) {
      return { ...step, status: 'upcoming' }
    }
    if (step.done) {
      return { ...step, status: 'done' }
    }
    if (!foundCurrent) {
      foundCurrent = true
      return { ...step, status: 'current' }
    }
    return { ...step, status: 'upcoming' }
  })
}

function joinquestMcpEnv({ apiKey }) {
  return {
    JOINQUEST_API_KEY: apiKey || '<generate-on-dashboard>',
  }
}

const MCP_NPX_PACKAGE = '@joinquest/mcp-integration'
export const INSTALL_DEV_SCRIPT_URL =
  'https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-dev.sh'
export const INSTALL_DEV_SCRIPT_GITHUB =
  'https://github.com/scruffyprodigy/playhub/blob/main/scripts/install-joinquest-dev.sh'

export function buildInstallDevCommand({ apiKey, client = 'cursor' }) {
  const flag =
    client === 'claude'
      ? '--claude'
      : client === 'claude-desktop'
        ? '--claude-desktop'
        : client === 'all'
          ? '--all'
          : client === 'skill-only'
            ? '--skill-only'
            : '--cursor'
  if (client === 'skill-only') {
    return `curl -fsSL ${INSTALL_DEV_SCRIPT_URL} | sh -s -- ${flag}`
  }
  const key = apiKey || 'lq_dev_PASTE_YOUR_KEY'
  return `export JOINQUEST_API_KEY=${key}
curl -fsSL ${INSTALL_DEV_SCRIPT_URL} | sh -s -- ${flag}`
}

export function buildInstallDevInspectCommand({ apiKey, client = 'cursor' }) {
  const flag =
    client === 'claude'
      ? '--claude'
      : client === 'claude-desktop'
        ? '--claude-desktop'
        : client === 'all'
          ? '--all'
          : client === 'skill-only'
            ? '--skill-only'
            : '--cursor'
  if (client === 'skill-only') {
    return `curl -fsSL ${INSTALL_DEV_SCRIPT_URL} -o install-joinquest-dev.sh
less install-joinquest-dev.sh
bash install-joinquest-dev.sh ${flag}`
  }
  const key = apiKey || 'lq_dev_PASTE_YOUR_KEY'
  return `curl -fsSL ${INSTALL_DEV_SCRIPT_URL} -o install-joinquest-dev.sh
less install-joinquest-dev.sh
export JOINQUEST_API_KEY=${key}
bash install-joinquest-dev.sh ${flag}`
}

export function buildClaudeMcpAddCommand({ apiKey }) {
  const key = apiKey || 'lq_dev_PASTE_HERE'
  return `claude mcp add --scope project --transport stdio \\
  --env JOINQUEST_API_KEY=${key} \\
  joinquest-integration -- npx -y ${MCP_NPX_PACKAGE}`
}

export function buildMcpServerConfig({ apiKey, clientId = 'cursor' }) {
  const env = joinquestMcpEnv({ apiKey })
  const args =
    clientId === 'cursor'
      ? ['--yes', '--package', MCP_NPX_PACKAGE, 'joinquest-integration-mcp-cursor']
      : ['-y', MCP_NPX_PACKAGE]
  return {
    mcpServers: {
      'joinquest-integration': {
        type: 'stdio',
        command: 'npx',
        args,
        env,
      },
    },
  }
}

export async function fetchMyDeveloperApiKeys() {
  const data = await graphqlRequest(MY_DEVELOPER_API_KEYS_QUERY)
  return data.myDeveloperApiKeys ?? []
}

export async function createDeveloperApiKey(name) {
  const data = await graphqlRequest(CREATE_DEVELOPER_API_KEY, { name: name ?? null })
  return data.createDeveloperApiKey
}

export async function revokeDeveloperApiKey(id) {
  const data = await graphqlRequest(REVOKE_DEVELOPER_API_KEY, { id })
  return data.revokeDeveloperApiKey
}

export function canRequestPublicRelease(game) {
  if (!game || game.visibility !== 'PRIVATE_TESTING') {
    return false
  }
  const hasShort = Boolean(game.shortDescription?.trim())
  const hasLong = Boolean(game.longDescription?.trim())
  const hasTags = Array.isArray(game.tags) && game.tags.length > 0
  const requiredChecks = REQUIRED_INTEGRATION_CHECKS
  const checksById = new Map((game.integrationChecks ?? []).map((c) => [c.checkId, c.status]))
  const checksPass = requiredChecks.every((id) => checksById.get(id) === 'PASS')
  return hasShort && hasLong && hasTags && checksPass
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
