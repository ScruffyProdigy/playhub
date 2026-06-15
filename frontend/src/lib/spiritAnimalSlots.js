/** Journey slot icons for the tarot reading (subset of starter avatars). */
const JOURNEY_SLOT_KEYS = ['compass', 'coin', 'storm', 'campfire', 'beacon']

/** What each journey slot is asking about (from product spec). */
const SLOT_PROMPTS = {
  compass: 'What guides the traveler on their journey.',
  coin: 'The value the traveler brings to others.',
  storm: 'What repeatedly tests the traveler.',
  campfire: 'Who the traveler becomes when people gather.',
  beacon: 'What the traveler ultimately becomes known for.',
}

const JOURNEY_SLOT_ICONS = {
  compass: { name: 'Compass', imageUrl: '/avatars/compass.png' },
  coin: { name: 'Coin', imageUrl: '/avatars/coin.png' },
  storm: { name: 'Storm', imageUrl: '/avatars/storm.png' },
  campfire: { name: 'Campfire', imageUrl: '/avatars/campfire.png' },
  beacon: { name: 'Beacon', imageUrl: '/avatars/beacon.png' },
}

/** Journey slots in tarot draw order (Compass → Beacon). */
export const JOURNEY_SLOTS = JOURNEY_SLOT_KEYS.map((key) => ({
  key,
  name: JOURNEY_SLOT_ICONS[key].name,
  imageUrl: JOURNEY_SLOT_ICONS[key].imageUrl,
  prompt: SLOT_PROMPTS[key] ?? 'A chapter of your journey.',
}))

export function slotPromptForKey(slotKey) {
  return SLOT_PROMPTS[slotKey?.toLowerCase()] ?? 'A chapter of your journey.'
}

export function journeySlotForKey(slotKey) {
  const key = slotKey?.toLowerCase()
  return JOURNEY_SLOTS.find((slot) => slot.key === key) ?? null
}
