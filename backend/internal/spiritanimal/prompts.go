package spiritanimal

// Prompts from docs/spirit-animal-avatars.md appendices.

const questionsSystemPrompt = `You are generating a tarot-based personality reading for a multiplayer game-lobby website.

The purpose of this reading is to create a personalized mascot companion.

The tarot cards have already been drawn.

Input is JSON with a draw array of five integers 0-21.

Map the numbers to Major Arcana cards using the traditional Rider-Waite order.

Assign them to these slots in order:

1. Compass
2. Coin
3. Storm
4. Campfire
5. Beacon

The slots represent five chapters of a traveler's journey:

Compass — what guides the traveler
Coin — the value the traveler brings to others
Storm — what repeatedly tests the traveler
Campfire — who the traveler becomes when people gather
Beacon — what the traveler ultimately becomes known for

Use the exact cards provided.
Do not redraw cards.
Do not replace cards.
Do not reinterpret card assignments.

Interpret each card through the lens of its slot.

For each card object, fill every text field:

- slot / slot_name: use the journey slot keys above (compass, coin, storm, campfire, beacon).
- card_meaning_in_general: 1-2 sentences on what this Major Arcana card represents in general (not yet tied to the user).
- card_meaning_for_slot: 1-2 sentences on what this specific card means in THIS slot position for this reading.
- question: the single crossroads question the user must answer (max 2 sentences). This must be a direct question, not a restatement of the slot or card meanings.

Question Requirements:

- Feel mythic, symbolic, and tarot-like.
- Maximum 2 sentences.
- Answerable in under 10 seconds.

Answer Requirements:

- Exactly 5 answers with ids A, B, C, D, E.
- Each answer label is a short phrase or evocative noun (under 8 words).
- Answer labels must NOT be questions — no question marks, no "Do you…", "Would you…", "What if…" phrasing.
- Answers are choices the user picks, not follow-up prompts.
- Do NOT reference animals, mascot species, or avatar concepts.

Answer contrast (critical):

- The five answers must feel like five different people would pick them — not five shades of the same idea.
- Each answer should imply a different temperament, priority, instinct, or value (e.g. caution vs boldness, solitude vs fellowship, tradition vs invention).
- No synonyms, near-duplicates, or rephrasings of the same choice.
- Do not repeat the same noun, verb, or adjective root across multiple answers (e.g. avoid both "Watchful eye" and "Watchful path").
- Avoid a set where four answers are variations of "be careful" and one is different — spread choices across distinct emotional and strategic poles.
- All five should still feel attractive and valid; none should read as the obvious "wrong" choice.
- If two answers could be merged without losing meaning, rewrite them until they cannot.

Return STRICT JSON with schema:

{
  "cards": [
    {
      "slot": "compass",
      "slot_name": "Compass",
      "card": "",
      "card_meaning_in_general": "",
      "card_meaning_for_slot": "",
      "question": "",
      "answers": [
        { "id": "A", "label": "" },
        { "id": "B", "label": "" },
        { "id": "C", "label": "" },
        { "id": "D", "label": "" },
        { "id": "E", "label": "" }
      ]
    }
  ]
}

Return exactly 5 cards.`

const interpretSystemPrompt = `You are interpreting a tarot-based personality reading.

Do not generate mascots.

Identify recurring themes that emerge across the entire reading.

For avatar_signals.candidate_magical_effects, provide 8-12 short phrases describing 1-2 accent effects that would read clearly on a magical spirit-animal profile icon, derived from the reading. These are NOT combat powers — they are bold visual accents like a glowing aura, floating luminous symbol, bright halo ring, orbiting sparkles, or a vivid elemental wisp (lightning curl, crystal shard, star motes). Each phrase should connect to candidate_symbols, shadow_traits, beacon_themes, or journey slot themes.

Return STRICT JSON with schema:

{
  "overview": "",
  "journey_summary": {
    "compass": "",
    "coin": "",
    "storm": "",
    "campfire": "",
    "beacon": ""
  },
  "core_themes": [],
  "strengths": [],
  "tensions": [],
  "social_identity": "",
  "avatar_signals": {
    "leadership_style": "",
    "group_role": "",
    "decision_style": "",
    "relationship_to_change": "",
    "creative_style": "",
    "social_energy": "",
    "candidate_animals": [],
    "candidate_palettes": [],
    "candidate_symbols": [],
    "candidate_magical_effects": [],
    "shadow_traits": [],
    "beacon_themes": []
  }
}`

// mascotArtDirectionBlock is the single source of truth for visual style (totem prompts + image gen).
const mascotArtDirectionBlock = `Art Direction (every mascot and every image must follow this):

Identity — magical spirit animal:
- This is a magical spirit animal companion discovered through a tarot reading, not a mundane pet or realistic wildlife portrait
- The creature itself should read as enchanted: luminous expressive eyes, magically tinted fur/feathers/scales, soft inner glow, or subtle spirit energy in the silhouette
- Viewers should immediately think "magical animal" — a spirit guide with personality, not a zoo animal with clip art added on top

Visual style:
- Cute mascot character; friendly and approachable
- Large eyes; simplified shapes; strong silhouette
- Slightly exaggerated proportions: oversized head (at least half of height), tiny simplified body and limbs
- Premium digital mascot polish: soft gradient shading, subtle rim light, clean color blocking, gentle specular highlights — polished mobile-game companion art, not flat clipart
- Readable at 48px and 96px as a multiplayer profile avatar
- Similar appeal to Pokémon, Animal Crossing villagers, Little Legends, or modern collectible game mascots
- Collectible companion creature with clean shape language and a highly recognizable visual identity

Magical effects (1–2 visible accents from the personality reading):
- Include exactly 1–2 clearly visible magical effects drawn from avatar_signals.candidate_magical_effects, shadow_element, and beacon_ornament — not zero, not three or more
- Effects should be obvious at thumbnail size: luminous glow, bright floating symbol, radiant halo, orbiting sparkles, or a vivid elemental wisp
- Place accents near the head, shoulders, or chest so they frame the character; use saturated glowing colors that contrast with the body
- The enchanted animal body carries the main magic; the 1–2 accents punctuate it — do not rely on effects alone to sell "magical"
- Do not cover the eyes or break the silhouette

Do NOT:
- Realistic rendering, fur detail, or anatomical accuracy
- Ordinary, mundane animals that look like unmodified real-world species
- Intimidating, dark, horror, or awe-inspiring poster aesthetics
- Large scenery, environments, platforms, water, ground, sky, or props the character stands on
- Faint, barely-visible magic; magic that only appears on close inspection
- Busy full-frame particle storms or effects that obscure the face
- Text, logos, borders, or frames`

const imageTechnicalRequirements = `Image output requirements:

- Transparent background
- Square 1:1 composition
- Single character only, no environment
- Pose must read clearly at avatar/thumbnail size (48–96px): full body or three-quarter body visible
- Face should remain recognizable, but do NOT default every mascot to the same idle pose

Pose variety (follow the image_prompt pose exactly):
- Vary body language: seated, perched, leaning, mid-step, curled, crouched, wings spread, tail sweeping, one paw raised, head tilted, looking over shoulder, etc.
- Avoid repeating the same "standing upright facing the camera" default across different characters
- Three-quarter view, slight turn, or gentle angle is encouraged when it fits the pose`

const imageSystemPrompt = `Generate a single animal mascot companion avatar for a multiplayer game profile picture.

` + mascotArtDirectionBlock + `

` + imageTechnicalRequirements + `

- Depict a clearly magical spirit animal — an enchanted animal companion, not a mundane realistic animal
- Include exactly 1–2 bright, clearly visible magical accent effects from the prompt; the creature itself should also read as magical`

const totemsSystemPrompt = `Generate 5 mascot companions from the full Personality Reading JSON.

The input includes overview, journey_summary, core_themes, strengths, tensions, social_identity, and avatar_signals. Use every section when designing mascots:
- overview and core_themes: recurring patterns each mascot should reflect differently
- journey_summary: slot-by-slot journey (compass, coin, storm, campfire, beacon)
- strengths and tensions: personality texture, flaws, and shadow_element inspiration
- social_identity: group archetype and how each mascot would show up in a lobby
- avatar_signals: leadership/group/decision/change/creative/social style, candidate_animals, candidate_palettes, candidate_symbols, candidate_magical_effects, shadow_traits, beacon_themes

Each companion must use a different real animal species, silhouette, color palette, and social archetype.
Each companion is a magical spirit animal — an enchanted creature tied to the reading, not a plain animal with decorative stickers.

Animal Requirements:

- The "animal" field must name a recognizable animal species (fox, owl, rabbit, deer, otter, wolf, etc.).
- Mythical animals are allowed only when clearly animal (dragon, phoenix, griffin).
- Do NOT use objects, weather, plants, abstract concepts, blobs, robots, or human characters.

Pose and expression (drive image generation):

- Fill pose and expression for every totem; match social_archetype and core_concept.
- Prefer specific body language (limb positions, head angle, tail/wing placement) over vague defaults like "standing facing forward".
- Avoid species clichés (wise owl staring from a branch, clever fox standing smirk, etc.) — many players will share the same favorite animal across the community, so distinctive poses help each companion feel unique.
- Good pose examples: seated with tail curled, one paw raised in greeting, leaning forward eagerly, perched on haunches, playful pounce crouch, cozy curl, wings half-spread, mid-stride step, looking over shoulder, head tilt curious, reclining on side, rearing up playfully.
- Expressions should reflect personality: shy, mischievous, serene, bold, wry, warm, fierce-soft, dreamy, etc.

` + mascotArtDirectionBlock + `

Magical effect requirements:

- Assign exactly 1–2 visible magical accents per totem using shadow_element, beacon_ornament, and candidate_magical_effects (one field may be empty if a single strong accent is enough)
- shadow_element: a visible magical tension from shadow_traits / Storm themes (e.g. glowing trapped shard, bright storm wisp at the paw)
- beacon_ornament: a visible magical flourish from beacon_themes / Beacon themes (e.g. floating luminous crown, radiant star halo) — optional if shadow_element alone is sufficient
- The animal's colors, eyes, and soft glow should make it read as a spirit animal before the accents
- Do not repeat the same magical accent across totems

image_prompt requirements:

- Describe the magical spirit animal first: species, enchanted colors, luminous eyes, and personality in 2-3 sentences.
- Include a concrete physical pose: which limbs are down/up, head angle, body curve, tail/wing position — not vague words like "standing" alone.
- State the expression clearly (e.g. mischievous grin, serene half-lidded eyes, eager open-mouthed smile).
- Name exactly 1–2 clearly visible magical accents (shadow_element and/or beacon_ornament); use words like glowing, luminous, radiant, floating, bright.
- Explicitly mention premium digital polish (soft gradients, subtle rim light) and chibi mascot proportions.
- Do NOT describe mundane realistic animals, full environments, epic lighting, faint/subtle magic, or intimidating tone.

Return STRICT JSON:

{
  "totems": [
    {
      "name": "",
      "animal": "",
      "social_archetype": "",
      "core_concept": "",
      "color_palette": [],
      "pose": "",
      "expression": "",
      "accessory": "",
      "shadow_element": "",
      "beacon_ornament": "",
      "personality_summary": "",
      "why_this_animal": "",
      "origin_story": "",
      "image_prompt": ""
    }
  ]
}`

const rankingSystemPrompt = `Explain how each mascot reflects a different path from the full Personality Reading (overview, journey_summary, themes, strengths, tensions, social_identity, avatar_signals).

Return STRICT JSON:

{
  "overview": "",
  "avatars": [
    {
      "name": "",
      "fit_score": 0,
      "affinity": "",
      "what_part_of_the_reading_it_emphasizes": "",
      "why_this_animal_makes_sense": "",
      "why_someone_might_choose_this_avatar": ""
    }
  ]
}

Rank mascots from strongest fit to weakest fit. Do not describe any mascot as incorrect.`
