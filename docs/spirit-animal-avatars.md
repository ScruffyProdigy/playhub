# Spirit animal avatars (planned)

**Status:** TODO — not implemented. Design spec for profile avatars at tables and elsewhere.

## Goals

- Give each player a recognizable visual identity at tables (and in room chat, banners, etc.).
- Offer **~25 starter avatars** to pick from directly.
- Offer **Find my spirit animal** — a tarot-guided, AI-assisted flow that generates **5 mascot options** for the user to choose from.

Avatars replace (or sit beside) display names in seat UI; seat role labels (e.g. `Clue Giver · Red`) stay as text.

## Avatar sources

| Source | Count | When |
|--------|-------|------|
| Starter catalog | ~25 | User picks from grid anytime |
| Spirit animal flow | 5 generated options | User completes tarot + personality flow |

Stored on user profile: `avatar_url` (public link for games), optional `avatar_key` / `avatar_source`, and for spirit-animal users a persisted **`avatar_readings`** row with personality + totem prompts for regeneration. See [spirit-animal-avatars.md](./spirit-animal-avatars.md).

---

## Find my spirit animal — flow overview

```text
Draw 5 tarot cards (0–21)
    → LLM: card reading + 5 questions (JSON)
    → User answers A–E per question
    → LLM: personality / avatar_signals (JSON)
    → LLM: 5 totem concepts + image_prompt each (JSON)
    → Image gen: 5 mascot images (transparent, 1:1)
    → LLM: ranked fit explanations (JSON)
    → Present: personality overview → mascot overview → each mascot in rank order
```

---

## Step 1 — Draw cards

Use a **random-without-replacement** draw (not independent `Math.random` per slot):

- **5 distinct integers** from **0–21** (22 Major Arcana, Rider–Waite order, zero-indexed).
- Map index → card name via fixed lookup table (The Fool = 0 … The World = 21).

Assign drawn cards **in order** to slots:

1. **Compass** — what guides the traveler  
2. **Coin** — the value the traveler brings to others  
3. **Storm** — what repeatedly tests the traveler  
4. **Campfire** — who the traveler becomes when people gather  
5. **Beacon** — what the traveler ultimately becomes known for  

Persist `{ draw: [n,n,n,n,n], slots: [...] }` on the server session so the reading cannot be redrawn mid-flow.

---

## Step 2 — Generate questions (LLM)

**Input:** `{ "draw": [17, 3, 21, 8, 14] }` (example)

**Prompt:** (fixed system prompt — see appendix A below)

**Output schema:**

```json
{
  "cards": [
    {
      "slot": "compass",
      "slot_name": "Compass",
      "card": "The Star",
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
```

Requirements (from product):

- Questions: mythic, symbolic, tarot-like; max 2 sentences; answerable in &lt;10s; not personality-test or therapeutic tone.
- Answers: exactly 5; all attractive; no joke/good/bad; visually evocative; **no animals, mascot species, or avatar concepts** in questions/answers.
- Use **exact** drawn cards; do not redraw or reassign.

---

## Step 3 — User answers

UI: one question per screen (or compact wizard), five answers each (A–E).

Collect: `["D","B","C","A","E"]` aligned to card order.

---

## Step 4 — Interpret reading (LLM)

**Input:**

```text
Cards:
[PASTE card JSON from step 2]

User Answers:
["D","B","C","A","E"]
```

**Prompt:** appendix B

**Output schema:**

```json
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
    "shadow_traits": [],
    "beacon_themes": []
  }
}
```

---

## Step 5 — Generate totem concepts (LLM)

**Input:** `[PASTE avatar_signals JSON]`

**Prompt:** appendix C

**Output schema:**

```json
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
}
```

Requirements: 5 totems; each different animal, silhouette, palette, social archetype; no repeated motifs/accessories/ornaments/poses/emotional tones; cute mascot / Little Legends / AC villager tone; readable at 48px and 96px.

---

## Step 6 — Generate images (image model)

For **each** totem, call image generation with appendix D template, substituting `[PASTE IMAGE_PROMPT]` from that totem’s `image_prompt`.

Requirements: transparent background, 1:1, no environment/text/border; large head, simplified body, strong silhouette.

Store resulting URLs (or upload to object storage); associate with totem id.

---

## Step 7 — Rank and explain (LLM)

**Input:**

```text
Journey Signals:
[PASTE avatar_signals JSON]

Mascots:
[PASTE totems JSON]
```

**Prompt:** appendix E

**Output schema:**

```json
{
  "overview": "",
  "avatars": [
    {
      "name": "",
      "fit_score": 0,
      "affinity": "",
      "what_part_of_the_reading_it_emphasizes": "",
      "why_this_animal_makes_sense": "",
      "why_someone_might_choose_this_avatar"
    }
  ]
}
```

Rank strongest → weakest fit; all paths valid; none “incorrect.”

---

## Step 8 — Present to user

1. **Personality overview** — from step 4 (`overview`, `journey_summary`, themes).
2. **Mascot overview** — from step 7 (`overview`).
3. **Each mascot, ranked** — image + fit copy (`affinity`, emphasis, why choose).

User picks one → save as profile avatar (see **Persistence** below).

---

## Game site access (GraphQL)

Partner games already resolve players via authenticated `player(id:)` using the provision **`serviceToken`** (see [lobby-protocol-handoff.md](./lobby-protocol-handoff.md), [api.md](./api.md)).

**Requirement:** the selected or generated avatar must be available on that lookup — at minimum a **public URL** games can fetch and render.

### `PublicPlayer` schema (planned)

Extend the existing public profile (no email):

```graphql
type PublicPlayer {
  id: ID!
  displayName: String
  avatarUrl: String       # absolute HTTPS URL to current avatar image
  avatarSource: AvatarSource
}

enum AvatarSource {
  STARTER
  SPIRIT_ANIMAL
  NONE
}
```

**Example (game server):**

```graphql
query Player($id: ID!) {
  player(id: $id) {
    id
    displayName
    avatarUrl
    avatarSource
  }
}
```

**Notes:**

- `avatarUrl` is the **current** image (starter asset or latest generated render). Games should treat it as an opaque URL; cache with normal HTTP semantics.
- Personality text and image prompts are **not** exposed on `PublicPlayer` — games get the portrait only. Full reading data stays on JoinQuest for regeneration and profile UI.
- `User` / `me` (signed-in player) may expose richer fields (`personalitySummary`, reading history) for the JoinQuest shell; keep game-facing surface minimal.
- `users.avatar_url` already exists in the initial schema; wire store + GraphQL when implementing.

---

## Persistence & regeneration

Spirit-animal avatars must be **re-renderable** when art direction changes (new prompt template, model, or style guide). Store everything needed to re-run **image generation only**, without redoing the tarot draw or personality flow.

### What to keep

| Artifact | Source step | Purpose |
|----------|-------------|---------|
| Tarot draw + slot mapping | 1 | Audit trail; optional re-display of reading |
| Questions + card copy JSON | 2 | Show reading in profile; debug |
| User answers | 3 | Re-interpret only if prompts change (rare) |
| **Personality / `avatar_signals` JSON** | 4 | **Personality summary**; input to totems & regen |
| **Totems JSON** (incl. **`image_prompt` per totem**) | 5 | **Regenerate image** under new art direction |
| Ranking / fit copy JSON | 7 | Re-show mascot picker explanations |
| Selected totem id | 8 | Know which concept the user chose |
| Generated image URLs + **`art_direction_version`** | 6 | History; know when a regen is needed |

### Suggested tables

**`users`** (existing + additions):

- `avatar_url` — current public URL (already in schema)
- `avatar_key` — starter catalog id when `source = STARTER`
- `avatar_source` — `starter` \| `spirit_animal` \| null
- `avatar_reading_id` — FK to active reading (spirit animal only)

**`avatar_readings`** (new, one row per completed spirit-animal flow):

- `id`, `user_id`, `created_at`, `completed_at`
- `draw` — `int[]` (5 cards)
- `questions_json` — step 2 output
- `user_answers` — `string[]` (e.g. `["D","B","C","A","E"]`)
- `personality_json` — step 4 output (overview, journey_summary, avatar_signals, …)
- `totems_json` — step 5 output (all 5 totems with `image_prompt`)
- `ranking_json` — step 7 output
- `selected_totem_name` — user's choice
- `art_direction_version` — integer; bump globally when style template changes

**`avatar_renders`** (new, optional but useful):

- `id`, `reading_id`, `totem_name`, `art_direction_version`
- `image_url` — stored asset URL for that render
- `image_prompt` — copy used (from totem at generation time)
- `created_at`

On **art direction bump**, batch job: for each reading with `selected_totem_name`, load `totems_json`, take that totem's `image_prompt`, run step 6 with the **new** global style wrapper, write new `avatar_renders` row, update `users.avatar_url`.

Starter avatars skip `avatar_readings`; `avatar_key` maps to static CDN path, no regen.

### Regeneration flow (admin / deploy)

```text
Bump ART_DIRECTION_VERSION + new image wrapper prompt
    → For each user with avatar_source = SPIRIT_ANIMAL
    → Load personality_json + totems_json + selected_totem_name
    → image_prompt = totems[selected].image_prompt
    → Generate image (appendix D + new wrapper)
    → Upload; insert avatar_renders; update users.avatar_url
```

Personality summary and totem concepts stay stable; only pixels change.

---

## UI / product notes

- **Starter grid:** ~25 pre-made assets; quick path for users who skip the flow.
- **Table display:** avatar chip + seat label; name on hover/long-press for a11y.
- **In-progress sessions:** server-side state machine; expire stale flows; no client-only redraw of tarot.
- **Cost / latency:** 4 LLM calls + 5 image gens per spirit-animal run — consider rate limits, queue, and “generating…” UX.
- **Moderation:** validate LLM JSON; fallback if parse fails; block unsafe image prompts.

---

## Implementation checklist (when scheduled)

- [ ] Wire `users.avatar_url` in store + expose on `User` / `me`
- [ ] Extend `PublicPlayer` with `avatarUrl`, `avatarSource`; game `player(id:)` resolver
- [ ] Migration: `users.avatar_key`, `users.avatar_source`, `users.avatar_reading_id`
- [ ] Migration: `avatar_readings` + `avatar_renders` (JSONB columns for LLM artifacts)
- [ ] Persist step 2–7 JSON + selected totem on flow completion
- [ ] Static starter catalog (~25 PNG/WebP in `frontend/public/avatars/` or CDN)
- [ ] Major Arcana index → name map (0–21)
- [ ] `SpiritAnimalSession` store + GraphQL mutations/queries
- [ ] OpenAI (or configured) client with strict JSON mode / schema validation
- [ ] Image generation integration + object storage for `avatar_url`
- [ ] Frontend wizard: draw → questions → generating → results → pick
- [ ] `TableCard` / room UI: render avatars from `avatarUrl`
- [ ] Admin/batch regen script keyed on `art_direction_version`
- [ ] Tests: tarot draw uniqueness, JSON schema validation, session lifecycle, `player` returns `avatarUrl`

---

## Appendix A — Step 2 prompt (questions)

```
You are generating a tarot-based personality reading for a multiplayer game-lobby website.

The purpose of this reading is to create a personalized mascot companion.

The tarot cards have already been drawn.

Input:

{
  "draw": [17, 3, 21, 8, 14]
}

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

The card determines:
- imagery
- symbolism
- metaphors
- archetypes

The slot determines:
- what chapter of the journey is being revealed

Before writing a question:

1. Identify the card's symbols.
2. Identify what the slot is trying to reveal.
3. Combine both.

A tarot reader should recognize both the card and the slot from the question.

Question Requirements:

- Feel mythic, symbolic, and tarot-like.
- Feel like a crossroads, omen, dream, mystery, legend, journey, or fable.
- Do not feel like a personality test.
- Do not feel therapeutic.
- Maximum 2 sentences.
- Answerable in under 10 seconds.

Answer Requirements:

- Exactly 5 answers.
- All answers should feel attractive and valid.
- No obvious good or bad choices.
- No joke answers.
- Each answer should represent a different expression of the card's energy.
- Each answer should reveal a different motivation, value, tendency, instinct, or worldview.
- Answers should be visually evocative.
- Do NOT reference animals.
- Do NOT reference mascot species.
- Do NOT reference avatar concepts.

Return STRICT JSON.

Schema:

{
  "cards": [
    {
      "slot": "",
      "slot_name": "",
      "card": "",
      "card_meaning_in_general": "",
      "card_meaning_for_slot": "",
      "question": "",
      "answers": [
        {
          "id": "A",
          "label": ""
        }
      ]
    }
  ]
}

The card_meaning_for_slot should be written as part of a tarot reading.

Example:

"In the Compass position, The Emperor speaks of foundations, responsibility, and the desire to create order from uncertainty."

Return exactly 5 cards.
```

---

## Appendix B — Step 4 prompt (interpretation)

```
Cards:

[PASTE CARD JSON]

User Answers:

["D","B","C","A","E"]

You are interpreting a tarot-based personality reading.

Do not generate mascots.

Identify recurring themes that emerge across the entire reading.

Return STRICT JSON.

Schema:

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
    "shadow_traits": [],
    "beacon_themes": []
  }
}

Requirements:

- Focus on recurring patterns.
- Explain why themes emerged.
- Prioritize traits useful for mascot creation.
- Candidate animals are symbolic possibilities, not final recommendations.
- Shadow traits should influence imperfections, flaws, asymmetry, or tension in future mascots.
- Beacon themes should influence mythology, ornaments, and social identity.
- The journey_summary should read like a concise tarot interpretation.
```

---

## Appendix C — Step 5 prompt (totems)

```
Avatar Signals:

[PASTE SIGNAL JSON]

Generate 5 mascot companions.

Each companion should represent a different interpretation of the same person.

Each companion should represent a different type of player someone might enjoy meeting in a game lobby.

The companions should feel discovered rather than assigned.

Requirements:

- Each companion must use a different animal.
- Each companion must use a different silhouette.
- Each companion must use a different color palette.
- Each companion must use a different social archetype.
- No repeated visual motifs.
- No repeated accessories.
- No repeated ornaments.
- No repeated poses.
- No repeated emotional tones.

Art Direction:

- Cute mascot companion.
- Not realistic.
- Not majestic spirit-animal art.
- Not epic fantasy concept art.
- Similar emotional appeal to Pokémon, Little Legends, mascot companions, Animal Crossing villagers, or game familiars.
- Large eyes.
- Strong silhouette.
- Readable at 48px.
- Readable at 96px.
- Designed as a multiplayer profile avatar.

The mascot should communicate:

- Who this player is in a group.
- Why people enjoy having them around.
- How they approach challenges.
- What part of the reading it emphasizes.

Return STRICT JSON.

Schema:

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
}
```

---

## Appendix D — Step 6 prompt (image)

```
Generate a single mascot companion avatar.

Requirements:

- Transparent background
- Square composition
- 1:1 aspect ratio
- Designed to function as a profile avatar
- Readable at 48px
- Readable at 96px
- Large head
- Simplified body
- Strong silhouette
- No environment
- No platform
- No scenery
- No text
- No logo
- No border
- No frame

Visual Style:

- Collectible game mascot
- Little Legend
- Companion creature
- Clean shape language
- Cute and approachable
- Highly recognizable
- Strong visual identity

[PASTE IMAGE_PROMPT]
```

---

## Appendix E — Step 7 prompt (ranking)

```
Journey Signals:

[PASTE SIGNALS]

Mascots:

[PASTE TOTEMS]

Explain how each mascot reflects a different path that emerged from the reading.

Return STRICT JSON.

Schema:

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

Rank the mascots from strongest fit to weakest fit.

The strongest fit should represent the greatest number of recurring themes.

The weaker fits should still be valid interpretations.

Describe them as different paths emerging from the same reading.

Do not describe any mascot as incorrect.
```

---

## Major Arcana index (Rider–Waite, 0–21)

| Index | Card |
|------:|------|
| 0 | The Fool |
| 1 | The Magician |
| 2 | The High Priestess |
| 3 | The Empress |
| 4 | The Emperor |
| 5 | The Hierophant |
| 6 | The Lovers |
| 7 | The Chariot |
| 8 | Strength |
| 9 | The Hermit |
| 10 | Wheel of Fortune |
| 11 | Justice |
| 12 | The Hanged Man |
| 13 | Death |
| 14 | Temperance |
| 15 | The Devil |
| 16 | The Tower |
| 17 | The Star |
| 18 | The Moon |
| 19 | The Sun |
| 20 | Judgement |
| 21 | The World |
