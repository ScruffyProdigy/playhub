# JoinQuest developer setup scripts

Scripts game developers run to install the JoinQuest agent skill and MCP in a **game repo** (or as a Cursor/Claude plugin).

**Lobby contributors:** these are not for you — see [../README.md](../README.md) and [../../docs/lobby-maintenance.md](../../docs/lobby-maintenance.md).

## Complete file list (nothing else is loaded)

When you `curl | sh` the installer, this is the **entire** shell surface — no deploy scripts, no lobby internals:

| File | Role |
|------|------|
| [install-joinquest-dev.sh](../install-joinquest-dev.sh) | Entry point — parse flags, orchestrate install |
| [lib/joinquest-skill.sh](../lib/joinquest-skill.sh) | Download/copy `.agents/skills/joinquest-integration/` into your repo |
| [lib/joinquest-mcp-env.sh](../lib/joinquest-mcp-env.sh) | Read `JOINQUEST_API_KEY` |
| [lib/joinquest-mcp-config.sh](../lib/joinquest-mcp-config.sh) | Write MCP JSON (Cursor, Copilot, Roo, etc.) |
| [lib/joinquest-platform.sh](../lib/joinquest-platform.sh) | Copilot / Roo / Cline / Windsurf rules files |

**Also downloaded (not shell):** a GitHub tarball of this repo; only `.agents/skills/joinquest-integration/` is kept (SKILL.md, playbook, mcp-setup).

**At agent runtime (not install time):** `npx @joinquest/mcp-integration` when your editor starts the MCP server.

**Not used:** `deploy-*.sh`, `dev.sh`, `lobby-*.sh`, or any other file under `scripts/`.

## Two install modes

| Mode | Script | Where files go |
|------|--------|----------------|
| **Game repo** | `install-joinquest-dev.sh --<platform>` | Your project (`.cursor/mcp.json`, `.agents/skills/`, etc.) |
| **Global plugin** | `install-joinquest-cursor-plugin.sh` or `install-joinquest-claude-plugin.sh` | `~/.cursor/plugins/` or `~/.claude/skills/` — single self-contained script, no `lib/` chain |

Platform flags for `install-joinquest-dev.sh`: `--cursor`, `--claude`, `--claude-desktop`, `--copilot`, `--roo`, `--windsurf`, `--cline`, `--skill-only`, `--all`.

## Review before running

```bash
# See what would happen (no downloads, no writes)
curl -fsSL https://raw.githubusercontent.com/scruffyprodigy/playhub/main/scripts/install-joinquest-dev.sh | sh -s -- --dry-run --copilot

# Download all five shell files, then read them
mkdir -p joinquest-setup-lib
curl -fsSL …/install-joinquest-dev.sh -o install-joinquest-dev.sh
curl -fsSL …/lib/joinquest-skill.sh -o joinquest-setup-lib/joinquest-skill.sh
# … (see developer dashboard “review first” copy for the full block)
less install-joinquest-dev.sh joinquest-setup-lib/*.sh
```

## Plugin install scripts (separate, no lib chain)

| Script | Purpose |
|--------|---------|
| [install-joinquest-cursor-plugin.sh](../install-joinquest-cursor-plugin.sh) | Cursor plugin → `~/.cursor/plugins/local/joinquest` |
| [install-joinquest-claude-plugin.sh](../install-joinquest-claude-plugin.sh) | Claude Code plugin → `~/.claude/skills/joinquest-integration` |
