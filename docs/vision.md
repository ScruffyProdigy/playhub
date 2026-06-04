# Product vision

**JoinQuest** (this repository) is the platform that lets a web developer **ship the multiplayer game they actually care about** instead of losing momentum on everything around it.

---

## The problem

You have an idea for a multiplayer game. You can see it going viral, growing a community, maybe even making real money. To get there you first need:

- Accounts and sign-in
- A lobby and matchmaking
- Queues, parties, and “find a game” UX
- A storefront and entitlements
- Secure payments and fraud basics
- Ops, scaling, and cross-origin auth between your site and the game client

That plumbing is necessary, but it is rarely what excited you about the project. Weeks go into infrastructure; the fun part stalls. Interest fades. **The game never ships.**

---

## What we are building

JoinQuest is **shared infrastructure for indie and small-studio web games**:

| You focus on | We handle (now or on the roadmap) |
|--------------|-----------------------------------|
| Game rules, feel, art, your brand | Player accounts and session auth |
| Your own website and URL | Catalog, queues, matchmaking, parties |
| Game server and client | Provision handoff, seat assignment, JWT link-out |
| Optional: your own shop later | Discovery, traffic, digital goods, payments |

**Integrate with us** and players can discover your title, queue up, get matched into seats, and land on your game with a signed identity — without you rebuilding a lobby product first.

You **keep independence**: your game still lives on **your** domain and codebase. You can add your own login, billing, or analytics anytime. Most teams won’t need to — the point is to do the boring parts **once**, well, for everyone.

---

## How integration works (mental model)

```text
Players  →  JoinQuest (accounts, catalog, queues, LFG, commerce)
                    │
                    │  provision + per-seat JWT
                    ▼
             Your game (your play_url, your rules)
```

1. **Register** your game API and publish a **seat manifest** (what seats exist; we decide who fills them).
2. **Players** use JoinQuest to browse, queue, and party up.
3. When a match is ready, we **push** the roster to your server and send each player to your URL with a **short-lived seat token**.
4. Your game **claims** the assigned `seatKey` and runs the match. You never run matchmaking unless you choose to.

Wire details: [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md). Layout and LFG: [`seat-templates-and-matchmaking.md`](./seat-templates-and-matchmaking.md).

---

## Who this repo is for

| Reader | Start here |
|--------|------------|
| **Visitor / potential game partner** | This doc → [`README.md`](../README.md) |
| **Game developer integrating** | [`lobby-protocol-handoff.md`](./lobby-protocol-handoff.md), [`game-catalog-architecture.md`](./game-catalog-architecture.md) |
| **JoinQuest contributor** | [`development.md`](./development.md), [`architecture.md`](./architecture.md), [`contributing.md`](./contributing.md) |
| **Agent / implementer** | [`seat-templates-and-matchmaking.md`](./seat-templates-and-matchmaking.md), [`game-catalog-architecture.md`](./game-catalog-architecture.md) |

---

## Principles

1. **Games own fun; we own plumbing** — Matchmaking, identity across origins, and catalog UX stay on our side unless a title explicitly opts out.
2. **Contract clarity** — Games receive a **final seat map** (`seatKey` + user id), not queue algorithms. Less to misunderstand, easier to test.
3. **Fail closed, recover gracefully** — Bad manifests don’t list; banned users trigger re-matchmake, not silent broken matches.
4. **Same platform, many games** — One account, one inventory layer, many `play_url`s — the “arcade” experience for players and a distribution channel for authors.

---

## Status

Core **catalog, auth, FIFO queues, provision, and JWT handoff** are implemented. **Template-based LFG**, richer parties, and full **commerce** are specified and in progress. See feature lists in [`README.md`](../README.md) and implementation notes in the architecture docs.
