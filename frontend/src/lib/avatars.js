import { graphqlRequest } from './graphql'

export const PROVISIONAL_DISPLAY_SUFFIX = ' (new)'

export const USER_AVATAR_FIELDS = `
  id
  displayName
  avatarUrl
  avatarKey
`

export const STARTER_AVATARS_QUERY = `
  query StarterAvatars {
    starterAvatars {
      key
      name
      slot
      imageUrl
    }
  }
`

const UPDATE_PLAYER_PROFILE = `
  mutation UpdatePlayerProfile($displayName: String!, $avatarKey: ID) {
    updatePlayerProfile(displayName: $displayName, avatarKey: $avatarKey) {
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

/** Static fallback when the catalog query is unavailable (tests, offline). */
export const STARTER_AVATAR_FALLBACK = [
  { key: 'compass', name: 'Compass', slot: 'Compass', imageUrl: '/avatars/compass.png' },
  { key: 'coin', name: 'Coin', slot: 'Coin', imageUrl: '/avatars/coin.png' },
  { key: 'storm', name: 'Storm', slot: 'Storm', imageUrl: '/avatars/storm.png' },
  { key: 'campfire', name: 'Campfire', slot: 'Campfire', imageUrl: '/avatars/campfire.png' },
  { key: 'beacon', name: 'Beacon', slot: 'Beacon', imageUrl: '/avatars/beacon.png' },
]

export function resolveUserAvatarUrl(user) {
  const direct = user?.avatarUrl?.trim()
  if (direct) {
    return direct
  }
  const key = user?.avatarKey?.trim()
  if (!key) {
    return null
  }
  const known = STARTER_AVATAR_FALLBACK.find((item) => item.key === key)
  return known?.imageUrl ?? `/avatars/${key}.png`
}

export function isProvisionalDisplayName(name) {
  return Boolean(name?.trim().endsWith(PROVISIONAL_DISPLAY_SUFFIX))
}

export function hasExistingAvatar(user) {
  return Boolean(user?.avatarKey?.trim())
    || Boolean(user?.avatarUrl?.trim())
    || user?.avatarSource === 'SPIRIT_ANIMAL'
}

export function needsProfileSetup(user) {
  return !hasExistingAvatar(user) || isProvisionalDisplayName(user?.displayName)
}

export function defaultDisplayNameInput(user) {
  const name = user?.displayName?.trim() || ''
  if (isProvisionalDisplayName(name)) {
    return name.slice(0, -PROVISIONAL_DISPLAY_SUFFIX.length).trim()
  }
  return name
}

export async function fetchStarterAvatars() {
  try {
    const data = await graphqlRequest(STARTER_AVATARS_QUERY)
    return data.starterAvatars?.length ? data.starterAvatars : STARTER_AVATAR_FALLBACK
  } catch {
    return STARTER_AVATAR_FALLBACK
  }
}

export async function updatePlayerProfile(displayName, avatarKey) {
  const variables = { displayName }
  if (avatarKey?.trim()) {
    variables.avatarKey = avatarKey.trim()
  }
  const data = await graphqlRequest(UPDATE_PLAYER_PROFILE, variables)
  return data.updatePlayerProfile
}

export function avatarInitial(user) {
  const name = user?.displayName?.trim() || 'Player'
  return name.charAt(0).toUpperCase()
}
