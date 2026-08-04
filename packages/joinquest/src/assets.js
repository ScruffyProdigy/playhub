import { existsSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const pkgRoot = join(dirname(fileURLToPath(import.meta.url)), '..')

export function skillSourceDir() {
  const bundled = join(pkgRoot, 'assets/skill')
  if (existsSync(join(bundled, 'SKILL.md'))) {
    return bundled
  }
  const monorepo = join(pkgRoot, '../../.agents/skills/joinquest-integration')
  if (existsSync(join(monorepo, 'SKILL.md'))) {
    return monorepo
  }
  throw new Error('JoinQuest skill assets not found. Reinstall the joinquest package.')
}

export function pluginSourceDir() {
  const bundled = join(pkgRoot, 'assets/plugin')
  if (existsSync(join(bundled, 'mcp.json'))) {
    return bundled
  }
  const monorepo = join(pkgRoot, '../../plugins/joinquest')
  if (existsSync(join(monorepo, 'mcp.json'))) {
    return monorepo
  }
  throw new Error('JoinQuest plugin assets not found. Reinstall the joinquest package.')
}
