---
name: joinquest-integration
description: >-
  Walks a developer through integrating a multiplayer game with JoinQuest — discovery
  interview, registration fields, MCP setup, game API implementation (healthz, provision,
  JWT), integration checks, catalog metadata, test tables, and public release. Use when
  the user mentions JoinQuest, joinquest.cc, lobby integration, matchmaking, game
  registration, hosting a game on JoinQuest, or integrating with the JoinQuest MCP server.
---

# JoinQuest integration

You are helping a developer integrate their **multiplayer game** with [JoinQuest](https://joinquest.cc). Follow the phases below in order. Read `playbook.md` in this skill folder for full detail, or call MCP `joinquest_integration_get_agent_playbook` for the latest copy from JoinQuest.

## Before you start

- Load **playbook.md** (bundled) or MCP `joinquest_integration_get_agent_playbook`.
- If JoinQuest MCP is configured, use MCP tools for all dashboard operations. If not, guide Phase 2 setup first (`mcp-setup.md`).
- For API/provision/JWT details, call MCP `joinquest_integration_get_integration_guide` — do not guess wire formats.

## Hard rules

1. **Interview before implementing** (Phase 1).
2. **Never** call `joinquest_integration_register_game`, `joinquest_integration_update_game_metadata`, or `joinquest_integration_request_public_release` without explicit human approval.
3. **No localhost** for `apiBaseUrl` — JoinQuest must reach a public HTTPS URL.
4. Fix the **game API** when checks fail, then re-run checks.

## Phases (summary)

| # | Phase | Done when |
|---|-------|-----------|
| 1 | Discover game | Drafts approved (copy, tags, seatTemplate plan) |
| 2 | Connect MCP | `joinquest_integration_list_my_games` works |
| 3 | Implement game API | healthz, status, game-modes, provision, claim on public HTTPS |
| 4 | Register on JoinQuest | Game registered via MCP or dashboard (`PRIVATE_TESTING` or draft) |
| 5 | Run checks | Required manifest + provision checks PASS |
| 6 | Save metadata | Human approved → `joinquest_integration_update_game_metadata` |
| 7 | Test with friends | Developer confirms test table / playtest |
| 8 | Public release | Human confirmed → `joinquest_integration_request_public_release` |

## MCP tool sequence (typical)

```
joinquest_integration_get_agent_playbook          # this workflow
joinquest_integration_get_discovery_prompt        # Phase 1 questions
joinquest_integration_get_catalog_tag_taxonomy    # tag IDs

# After MCP auth:
joinquest_integration_register_game               # Phase 4 — after developer confirms fields
joinquest_integration_list_my_games
joinquest_integration_get_game_credentials        # serviceToken (sensitive)
joinquest_integration_get_example_provision_payload
joinquest_integration_run_game_checks             # repeat until green
joinquest_integration_update_game_metadata        # after human approval
joinquest_integration_request_public_release      # after human approval
```

## Registration fields (Phase 4)

Confirm with the developer, then call `joinquest_integration_register_game` or use https://joinquest.cc/developers.

| Field | Required |
|-------|----------|
| Game name | Yes |
| Slug | Yes |
| Short description | Yes |
| API base URL (public HTTPS) | Yes |
| Contact email | Yes |
| Website / community URL | No |

## On check failures

1. Read check `message` and `detail`.
2. MCP `joinquest_integration_get_integration_guide` → §10 checklist index.
3. Fix game API, deploy, `joinquest_integration_run_game_checks` again.

## Status template

Report progress using the checklist in `playbook.md` § Phase checklist.
