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
  mutation UpdatePlayerProfile($displayName: String!, $avatarKey: ID!) {
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

const SELECT_STARTER_AVATAR = `
  mutation SelectStarterAvatar($key: ID!) {
    selectStarterAvatar(key: $key) {
      id
      displayName
      avatarUrl
      avatarKey
      avatarSource
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

export function isProvisionalDisplayName(name) {
  return Boolean(name?.trim().endsWith(PROVISIONAL_DISPLAY_SUFFIX))
}

export function needsProfileSetup(user) {
  return !user?.avatarKey || isProvisionalDisplayName(user?.displayName)
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
  const data = await graphqlRequest(UPDATE_PLAYER_PROFILE, { displayName, avatarKey })
  return data.updatePlayerProfile
}

export async function selectStarterAvatar(key) {
  const data = await graphqlRequest(SELECT_STARTER_AVATAR, { key })
  return data.selectStarterAvatar
}

export function avatarInitial(user) {
  const name = user?.displayName?.trim() || 'Player'
  return name.charAt(0).toUpperCase()
}
