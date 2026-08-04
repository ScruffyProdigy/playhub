import { join } from 'node:path'
import { copyTree, writeFile } from './fs-util.js'
import { skillPathInProject } from './skill.js'

const MCP_FIRST_INSTRUCTIONS = `# JoinQuest integration

Integrate a multiplayer game with [JoinQuest](https://joinquest.cc).

**Use MCP first:** at the start of any JoinQuest task, call \`joinquest_integration_get_agent_playbook\`.
Copilot and other agents may truncate long instruction files — the MCP playbook is the source of truth.

Verify: "List my JoinQuest games using the MCP tools."
Start: "I want to create a game on JoinQuest."

Requires the **joinquest-integration** MCP server (Node.js 20+, \`npx @joinquest/mcp-integration\`).
`

function requireSkill(cwd) {
  const skill = skillPathInProject(cwd)
  return skill
}

export function installCopilotArtifacts(cwd, { dryRun = false } = {}) {
  const skill = requireSkill(cwd)
  copyTree(skill, join(cwd, '.github/skills/joinquest-integration'), { dryRun })
  writeFile(join(cwd, '.github/copilot-instructions.md'), MCP_FIRST_INSTRUCTIONS, { dryRun })
}

export function installRooArtifacts(cwd, { dryRun = false } = {}) {
  const skill = requireSkill(cwd)
  copyTree(skill, join(cwd, '.roo/rules/joinquest-integration'), { dryRun })
}

export function installClineArtifacts(cwd, { dryRun = false } = {}) {
  const skill = requireSkill(cwd)
  copyTree(skill, join(cwd, '.cline/rules/joinquest-integration'), { dryRun })
  writeFile(join(cwd, '.clinerules'), MCP_FIRST_INSTRUCTIONS, { dryRun })
}

export function installWindsurfArtifacts(cwd, { dryRun = false } = {}) {
  const skill = requireSkill(cwd)
  copyTree(skill, join(cwd, '.windsurf/rules/joinquest-integration'), { dryRun })
}
