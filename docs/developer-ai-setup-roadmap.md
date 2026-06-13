# Developer AI setup roadmap

Track simplifying JoinQuest integration for game developers using Cursor, Claude Code, and other agents. **Skill + MCP** are the two layers; the goal is one mental model and minimal steps.

## Target experience

| Audience | Done when |
|----------|-----------|
| **Cursor** | Dashboard → copy config → Cmd+Q → Agent chat → verify MCP → “I want to create a game on JoinQuest” |
| **Claude Code** | Dashboard → copy `claude mcp add …` → new session → approve MCP → same prompt |
| **Everyone** | Agent skill installed in game repo automatically (not a separate GitHub copy) |

Production only for game devs: `JOINQUEST_API_KEY` in MCP env (no URL). Lobby contributors: see [development.md](./development.md) for local override.

---

## Tier 1 — Highest impact

| # | Item | Status | Notes |
|---|------|--------|-------|
| 1 | **Publish `@joinquest/mcp-integration` to npm** | Republish 0.1.1 | 0.1.0 missing default bin; Cursor wrapper path fix in 0.1.1 |
| 2 | **Dashboard copy-ready commands** | Done | One-click Copy Cursor JSON + Copy Claude command after key gen |
| 3 | **Unified setup script** (skill + MCP hints) | Done | `install-joinquest-dev.sh` — skill + `--cursor` / `--claude` / `--all` |
| 10 | **Skill install automation** | Done | Bundled in `install-joinquest-dev.sh` |
| 4 | **Deploy doc/wizard/MCP fixes to joinquest.cc** | Done | Deployed commit 816bb8f to joinquest.cc |

## Tier 2 — Remove secrets & friction

| # | Item | Status | Notes |
|---|------|--------|-------|
| 5 | **`joinquest login` / device auth for MCP** | Pending | No API key in `.cursor/mcp.json` |
| 6 | **Single `setup.sh --cursor \| --claude`** | Pending | One URL in docs and dashboard |
| 7 | **Cursor MCP reliability** | Ongoing | `joinquest-integration-mcp-cursor` wrapper; consider HTTP MCP later |

## Tier 3 — Polish & reach

| # | Item | Status | Notes |
|---|------|--------|-------|
| 8 | **npx-first docs everywhere** | In progress | After npm publish, install script = fallback only |
| 9 | **Claude Desktop / Copilot** | Partial | `--claude-desktop` install flag; Copilot still manual |
| 10 | **Skill install automation** | Done | Bundled in `install-joinquest-dev.sh` |
| 11 | **Dashboard API key management UX** | Pending | Clear list/revoke/regenerate for MCP |

## Explicitly not now

- Per-IDE extensions
- Localhost in game-developer funnel
- Assuming Cursor “Tools & MCP” settings UI exists

---

## MCP config reference (npx, post-publish)

**Cursor** — needs cursor bin via npx (Electron env workaround):

```json
{
  "mcpServers": {
    "joinquest-integration": {
      "type": "stdio",
      "command": "npx",
      "args": ["--yes", "--package", "@joinquest/mcp-integration", "joinquest-integration-mcp-cursor"],
      "env": { "JOINQUEST_API_KEY": "lq_dev_..." }
    }
  }
}
```

**Claude Code:**

```bash
claude mcp add --scope project --transport stdio \
  --env JOINQUEST_API_KEY=lq_dev_... \
  joinquest-integration -- npx -y @joinquest/mcp-integration
```

**Other stdio clients:**

```json
"command": "npx",
"args": ["-y", "@joinquest/mcp-integration"]
```

---

## Changelog

- 2026-06-13 — Roadmap created; MCP package npx-first configs landed; publish pending npm auth.
- 2026-06-13 — `@joinquest/mcp-integration@0.1.0` published to npm.
