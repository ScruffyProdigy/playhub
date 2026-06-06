import { graphqlRequest } from './graphql'

export const USER_AVATAR_FIELDS = `
  id
  displayName
  avatarUrl
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

export async function fetchStarterAvatars() {
  try {
    const data = await graphqlRequest(STARTER_AVATARS_QUERY)
    return data.starterAvatars?.length ? data.starterAvatars : STARTER_AVATAR_FALLBACK
  } catch {
    return STARTER_AVATAR_FALLBACK
  }
}

export async function selectStarterAvatar(key) {
  const data = await graphqlRequest(SELECT_STARTER_AVATAR, { key })
  return data.selectStarterAvatar
}

export function avatarInitial(user) {
  const name = user?.displayName?.trim() || 'Player'
  return name.charAt(0).toUpperCase()
}
