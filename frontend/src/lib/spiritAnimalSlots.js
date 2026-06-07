import { STARTER_AVATAR_FALLBACK } from './avatars'

/** What each journey slot is asking about (from product spec). */
const SLOT_PROMPTS = {
  compass: 'What guides the traveler on their journey.',
  coin: 'The value the traveler brings to others.',
  storm: 'What repeatedly tests the traveler.',
  campfire: 'Who the traveler becomes when people gather.',
  beacon: 'What the traveler ultimately becomes known for.',
}

/** Journey slots in tarot draw order (Compass → Beacon). */
export const JOURNEY_SLOTS = STARTER_AVATAR_FALLBACK.map((avatar) => ({
  key: avatar.key,
  name: avatar.name,
  imageUrl: avatar.imageUrl,
  prompt: SLOT_PROMPTS[avatar.key] ?? 'A chapter of your journey.',
}))

export function slotPromptForKey(slotKey) {
  return SLOT_PROMPTS[slotKey?.toLowerCase()] ?? 'A chapter of your journey.'
}

export function journeySlotForKey(slotKey) {
  const key = slotKey?.toLowerCase()
  return JOURNEY_SLOTS.find((slot) => slot.key === key) ?? null
}
