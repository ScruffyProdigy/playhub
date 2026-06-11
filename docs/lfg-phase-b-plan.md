# LFG Phase B — Implementation Plan

**Status:** Shipped (Phase B)  
**Spec:** [seat-templates-and-matchmaking.md](./seat-templates-and-matchmaking.md) (full engine; Phase C adds weighted dequeue)  
**Builds on:** Phase A (`seatTemplate` expansion, `queuePath` join, per-path `sizeForQueue`)

Phase B replaces the one-shot FIFO snapshot with a **persistent forming match** per mode queue, **solo catalog LFG**, and **table backfill** (friends sit first, then king opens the table to the lobby). Phase C adds weighted dequeue and `allocations` by affinity.

**Product decision:** **No frontend party role picker.** Friends use **room → table → sit in role seats → Look for group**. Catalog **Join as …** stays solo-only. Backend **`PartyNodeInput` tree** for tests/API; tables build the tree from seated `seatKey`s.

**Party model:** Tree aligned with `seatTemplate` (`Team` → `SpyMaster` / `Guesser`, etc.). Root role optional/ignored. No `together` flag — structure replaces affinity hints. Table backfill uses **pinned** placement; catalog/API uses **tree** placement with branch permutation (e.g. two `Team` siblings → Team-1 / Team-2).

---

## Goal

1. Hold a **working seat map** per `mode_queue_id` until fire conditions are met.
2. **Solo catalog LFG** — existing **Join as …** / **Look for group**; incremental placement on the forming map.
3. **Table backfill** — seated friends seed the map; catalog fills remaining role gaps; same fire pipeline.
4. **Gap visibility** — show what roles are still needed on the **intent banner** (catalog wait + table backfill) and on **TableCard** (room).
5. **Fire** into the existing session → provision → JWT pipeline.

**Async reconcile:** `joinQueue` enqueues the player and returns `queued: true` immediately. A **forming worker** (see [pubsub.md](./pubsub.md)) runs reconcile on a short delay, fires when the map is ready, provisions the game server, and publishes per-user `queueUpdated` events. GraphQL queries and subscriptions read stored launch URL bases — they do not trigger provision.

**Out of scope for B:** frontend party role picker, weighted dequeue (C), `pins`, variable partial composition, `startMatch(party, startSize)`.

**Deferred (API only):** catalog party tree UI — use tables for coordinated friend groups instead.

---

## Current baseline

| Shipped | Location |
|---------|----------|
| `seatTemplate` expand + `PathSpecs` (`sizeForQueue` per path) | `backend/internal/seattemplate/` |
| Persistent forming match + incremental placement | `backend/internal/store/forming_*.go`, `party.go` |
| `joinQueue(queueId, queuePath, party: PartyNodeInput)` | GraphQL + `forming_join.go` |
| Table backfill (**Look for group**) | `table_backfill.go`, `TableCard.jsx` |
| Gap visibility (`formingGaps`) | intent banner + `TableCard` |
| Tables + king **Start now** | `backend/internal/store/table.go` |
| `affinity_key` column on seats | migration 000014 |
| Party layout tree (`party_tree` JSONB) | migration 000022, `internal/lfg/partytree/` |

---

## Architecture

```text
Catalog solo joinQueue(queuePath)
        or
Table backfill (seed from table_seats, source=table)
        ↓
  solo: implicit party-of-1 on game_queues
        ↓
  get or create FormingMatch for mode_queue_id
        ↓
  place on working map (table pins → solo FIFO per path)
        ↓
  gaps per queuePath → expose on intent banner + TableCard
        ↓
  if ready → FIRE → session + participants (existing)
        ↓
  provision + JWT (existing handoff)
```

**Friend groups:** create table → `sitAtTable(seatKey)` encodes role → king **Look for group** → gaps = `PathSpec` needs minus seated counts by path. No separate role-picker UI.

**One forming match per `mode_queue_id` at a time** (v1).

---

## Milestone 1 — Data model & types

### Migration `000021_forming_matches`

**`parties`**

| Column | Notes |
|--------|--------|
| `id` | UUID PK |
| `leader_user_id` | nullable for solo-as-party-of-1 |
| `mode_queue_id` | queue this party waits in |
| `status` | `waiting`, `placed`, `matched`, `cancelled` |
| `together` | bool — all members same affinity when placed |
| `created_at` | |

**`party_members`**

| Column | Notes |
|--------|--------|
| `party_id`, `user_id` | unique active row per user |
| `queue_path` | member role bucket (`ClueGiver`, `Guesser`, `SpyMaster`, `""`) |
| `sort_order` | stable UI order |

**`forming_matches`**

| Column | Notes |
|--------|--------|
| `id` | UUID PK |
| `mode_queue_id` | partial unique WHERE `status = 'filling'` |
| `mode_id`, `game_id` | denorm |
| `status` | `filling`, `ready`, `fired` |
| `target_spec` | JSON snapshot of `PathSpec[]` at creation |
| `created_at`, `fired_at` | |

**`forming_match_assignments`**

| Column | Notes |
|--------|--------|
| `forming_match_id`, `seat_key` | unique pair |
| `user_id` | nullable = empty cell |
| `party_id` | nullable for solos |
| `queue_path`, `affinity_key` | denorm from seat |
| `source` | `party`, `solo`, `table` |
| `table_id` | optional, for table backfill |

**Extend `game_queues`**

| Column | Notes |
|--------|--------|
| `party_id` | nullable FK |
| `forming_match_id` | set when placed on map |

### `seattemplate` — `affinity_key` derivation

| Template shape | Example seat | `affinity_key` |
|----------------|--------------|----------------|
| `Team: { count: 2, … }` | `Team-1-Guesser-1` | `Team:1` |
| Named leaf cohort | `ClueGiver-Red` | `ClueGiver:Red` |
| No side dimension | `Guesser-1` | *(empty)* |

Stored on `game_mode_seats.affinity_key` at catalog sync.

### Go packages

- `store`: `Party`, `PartyMember`, `FormingMatch`, `FormingAssignment`
- `internal/lfg`: gap vector, placement helpers, fire predicate

---

## Milestone 2 — GraphQL & gaps API

Solo `joinQueue` unchanged at the UI; optional `party: PartyNodeInput` on the mutation for tests only (no player-facing picker).

```graphql
extend type Mutation {
  joinQueue(queueId: ID!, queuePath: String, party: PartyNodeInput): JoinResult!
  startTableBackfill(tableId: ID!, queueId: ID!): JoinResult!
}

type QueuePathGap {
  queuePath: String!
  displayName: String!
  assigned: Int!
  needed: Int!
}

type FormingMatchStatus {
  modeQueueId: ID!
  assignedCount: Int!
  targetCount: Int!
  gaps: [QueuePathGap!]!
}

extend type GameMode {
  forming: FormingMatchStatus
}

extend type Table {
  backfillActive: Boolean!
  formingGaps: [QueuePathGap!]!
}

extend type ActiveIntent {
  formingGaps: [QueuePathGap!]!
}

extend type MyTableSeat {
  backfillActive: Boolean!
  formingGaps: [QueuePathGap!]!
}
```

| Rule | Behavior |
|------|----------|
| Solo catalog join | `party` omitted → implicit party of 1 |
| Friend groups | **Table seats**, not catalog `PartyNodeInput` |
| Gap copy | Human-readable `displayName` from manifest (`Clue Giver`, `Guesser`, …) |

---

## Milestone 3 — Forming engine

Replace `tryFormMatch` in `JoinModeQueue` with:

1. Resolve party (create or load).
2. Validate member `queuePath`s against mode seats.
3. Load or create `forming_match` (`status = filling`).
4. If cannot place: enqueue as waiter only.
5. If can place: assign seats (table pins → solo FIFO per path).
6. Recompute gaps; if ready → `fireFormingMatch()`.
7. Commit; notify affected users.

**Fire condition** (unchanged from Phase A composition):

```text
for each PathSpec: assigned[path] >= spec.PlayersToStart()
```

**B-lite relocation:** only unassigned parties entering the map — do not move already-placed solos (Phase C).

Forming LFG is always enabled (legacy FIFO snapshot removed).

---

## Milestone 4 — Acceptance: Word Hunt (table path)

> Pat and brother create a private table, sit as Clue Giver and Guesser, king clicks **Look for group**. Same match; fires at 2 CG + 4 guessers when catalog fills gaps.

| Step | Assert |
|------|--------|
| Pat sits `ClueGiver-Red`, brother sits `Guesser-1` | Roles from `seat_key`, not a party picker |
| King starts backfill | Forming map seeded; gaps CG +1, Guesser +3 |
| Four solos join catalog | One session; all six correct `seatKey`s |
| Intent banner (seated) | “Your table is waiting for players — need …” |
| TableCard | Same gaps visible in room |

| Test | Assert |
|------|--------|
| `TestTableBackfillWordHunt_PartialGaps` | 2 seated → gaps CG +1, Guesser +3 |
| `TestTableBackfillWordHunt_FireWithStrangers` | seated + 4 catalog solos → one session |
| Regression | `TestJoinModeQueueWordHuntFiresAtPartialCohortSizes` via solo catalog joins |
| Regression | `TestJoinModeQueueSplitPartyWordHunt` (API-only `PartyNodeInput`; optional) |

---

## Milestone 4b — Acceptance: Codenames-style composition

**Template** (correct grammar — `name` not `names`; `count: 2` for two teams; cannot combine `name` + `count` on one node):

```json
{
  "Team": {
    "count": 2,
    "SpyMaster": {},
    "Guesser": { "count": 3 }
  }
}
```

Expands to 8 seats, queue paths `SpyMaster` (2) and `Guesser` (6). Affinity `Team:1` / `Team:2` on all seats under each team block.

| Scenario | How friends do it (Phase B) |
|----------|----------------------------|
| Competing spymasters | Solo catalog **Join as SpyMaster** (or sit at each team's spymaster seat on separate tables — rare) |
| Spymaster + guesser, **same team** | Sit at `Team-1-SpyMaster` + `Team-1-Guesser-*` on same table → backfill |
| Both guessers, **same team** | Sit at two `Team-1-Guesser-*` seats on same table → backfill |

Solo catalog FIFO still works for all scenarios (with known Phase B placement limits — see [Known limitations](#known-limitations-phase-b)).

Phase C adds `allocations` by affinity for catalog-only 3v3 splits without a table.

---

## Milestone 5 — Table LFG backfill + gap visibility (UI)

### Backend

1. King clicks **Look for group** on partial table (`startTableBackfill` or equivalent).
2. Seed forming map from `table_seats` (`source = table`, `table_id` set).
3. Compute gaps from seated counts per `queuePath` vs `PathSpec.PlayersToStart()`.
4. Catalog solos fill remaining cells; on fire → `startTable` → session + provision.
5. Enable `lookForGroupOptions.enabled` in `table_internal.go`.

### Gap visibility — two surfaces, same data

| Surface | Audience | When | Example copy |
|---------|----------|------|----------------|
| **Intent banner** | Seated or waiting player (global sticky) | Catalog **waiting** OR table **backfill active** | Catalog: “Looking for a group in Word Hunt as Clue Giver · need 1 Clue Giver, 3 Guessers” · Table: “Your table is waiting for players — need 1 Clue Giver, 3 Guessers” |
| **TableCard** | Everyone in the room | Table `backfillActive` | “Need 1 Clue Giver, 3 Guessers from the lobby” under seats / Look for group |

**Intent banner states**

| Intent | Title (existing) | Hint / gap line (new) |
|--------|------------------|------------------------|
| Catalog waiting | `bannerWaitingLine` + role | Append or replace player count with **need** list from `formingGaps` |
| Table seated, not backfilling | `bannerTableSeatLine` | Existing: “Use Room below…” |
| Table seated, **backfill active** | Keep seat line | **“Your table is queued to start — need …”** + gaps |

Data: `myActiveIntent.formingGaps` (catalog), `myTableSeat.formingGaps` + `backfillActive` (table), `Table.formingGaps` on room subscription.

### Frontend files (expected)

- `IntentBanner.jsx` + `playerCopy.js` — gap formatting helper (e.g. `formatFormingGaps(gaps)`)
- `TableCard.jsx` — gaps block when `backfillActive`
- `queue.js` / `tables.js` — GraphQL fields for gaps

**Not building:** party role picker, “play with friend” on games list.

**Acceptance:** table with 1 CG + 1 Guesser seated → king backfills → gaps on TableCard + intent banner → +1 CG, +3 guessers from catalog → game starts.

---

## Milestone 6 — Realtime

- `queueUpdated` includes `formingGaps`
- `tableUpdated` includes `formingGaps`, `backfillActive`
- Reuse gap formatter in banner + TableCard on subscription refresh

---

## Milestone 7 — Rollout

| Item | Detail |
|------|--------|
| Forming engine | Always on (no feature flag) |
| Metrics | forming created, fire latency, split vs solo |
| Rollout | staging → Word Hunt first |

---

## Known limitations (Phase B)

Phase B uses **incremental placement** and does **not** relocate players already on the forming map. That can leave a valid match unfilled even when one exists among current waiters.

### Example: `{ "count": 3 }` with A, B solo then C+D party

| Step | Forming map | Waiters |
|------|-------------|---------|
| A joins (solo) | `1`→A, gaps: 2 | — |
| B joins (solo) | `1`→A, `2`→B, gaps: 1 | — |
| C+D join (`together`, size 2) | unchanged (1 gap) | C+D — party needs 2 seats; cannot place |

**Result in Phase B:** nothing fires yet. A third **solo** (E) would complete `{A, B, E}`. The valid roster `{A, C, D}` does not form because B already occupies a seat and will not be bumped.

This is acceptable for Phase B. It is the same family of issue as the [third pair waits for game 2](./seat-templates-and-matchmaking.md#a-three-teams-of-two-in-3v3) scenario — early placements can block better packings.

### Mitigation (Phase C dequeue)

When we add **weighted dequeue** (for skill parity, win-rate tracking, and smarter pulling from the waiter list), **party size** is already in the target score function — larger groups are harder to place later and should be **slightly favored** over long-waiting solos when re-solving who lands on the map.

Concretely for the example above: the match still does not fire *immediately* when C+D arrive, but **given enough time** C+D accumulate dequeue priority over solo B. On the next placement pass (or when relocation is enabled), the engine can prefer `{A, C, D}` — bumping B back to the waiter queue — rather than holding B on the map until a fourth solo appears.

Phase C factors (see [Dequeue policy](./seat-templates-and-matchmaking.md#dequeue-policy-not-strict-fifo)):

| Factor | Role here |
|--------|-----------|
| **Party size** | C+D (size 2) beat solo B over time |
| **Complementarity / fit** | Prefer assignments that fill the map and satisfy `together` |
| **Skill / MMR** | Same dequeue pass can bias toward even-skill lobbies |
| **Wait time** | Tie-break within tier — solos are not permanently locked in |

**Phase B:** document and accept the stall. **Phase C:** dequeue + optional relocation address it without requiring assign-at-fire for every fifo mode.

---

## Phase C+ deferrals

| Item | Phase |
|------|-------|
| Weighted dequeue (party size, skill/MMR, complementarity, age) | C |
| Relocate already-placed players when dequeue re-solves map | C |
| `allocations` by affinity (3v3 team splits) | C |
| `pins` on `seatKey` | C+ |
| Variable partial role counts | E |

See [Known limitations (Phase B)](./lfg-phase-b-plan.md#known-limitations-phase-b) for the `{ count: 3 }` / A+B vs C+D stall and how Phase C addresses it.

---

## Build order

| Step | Deliverable | Status |
|------|-------------|--------|
| 1 | Migration + types + `affinity_key` derivation | ✅ |
| 2 | Forming match CRUD + gap computation + unit tests | ✅ |
| 3 | `joinQueue` → forming engine; fire → session | ✅ |
| 4 | Table backfill backend + Word Hunt acceptance tests | ⬜ |
| 5 | Gap API + intent banner + TableCard gap UI | ⬜ |
| 6 | Subscriptions (`formingGaps` on queue + table) | ⬜ |
| 7 | Leave/disband edge cases, rollout | ⬜ |

**Removed from plan:** frontend party role picker (tables cover friend groups).

---

## Success criteria

1. **Word Hunt (table):** Pat + brother seated at CG + Guesser; king backfills; catalog fills; all six same match.
2. **Codenames-style (table):** Same-team seats via `Team-1-*` sitting; backfill fills remainder.
3. **Solo catalog** FIFO unchanged for duel and per-path composition.
4. **Gap visibility:** intent banner (catalog wait + table backfill) and TableCard show the same need list.
5. No regression in provision/JWT/handoff tests.
6. `formingGaps` on API types used by frontend (not only `GameMode.forming`).
