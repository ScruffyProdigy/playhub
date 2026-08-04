# JoinQuest developer setup

Install the JoinQuest agent skill and MCP with the **joinquest** npm package — no `curl | sh`.

## Recommended

```bash
JOINQUEST_API_KEY=lq_dev_... npx joinquest install cursor
```

Or:

```bash
npm create joinquest@latest -- --cursor
```

Preview without writing files:

```bash
npx joinquest install cursor --dry-run
```

## Platforms

`cursor` · `claude` · `claude-desktop` · `copilot` · `roo` · `windsurf` · `cline` · `skill`

Global plugin (Cursor / Claude Code only):

```bash
npx joinquest install cursor --plugin
```

## What gets installed

| Component | Where |
|-----------|--------|
| Agent skill | `.agents/skills/joinquest-integration/` (from npm package) |
| MCP runtime | `npx @joinquest/mcp-integration` (configured in your editor) |

**Complete install surface:** the `joinquest` npm package only. See [packages/joinquest/README.md](../../packages/joinquest/README.md).

## Legacy shell wrappers

These forward to `npx joinquest` for old URLs and docs:

- `install-joinquest-dev.sh`
- `install-joinquest-cursor-plugin.sh`
- `install-joinquest-claude-plugin.sh`

## Lobby scripts

Deploy, test, and database scripts in `scripts/` are for **JoinQuest platform contributors** only — not game developers. See [scripts/README.md](../README.md).
