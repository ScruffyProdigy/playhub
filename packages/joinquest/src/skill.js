import { join } from 'node:path'
import { copyTree } from './fs-util.js'
import { skillSourceDir } from './assets.js'

export function installAgentSkill(cwd, { dryRun = false } = {}) {
  const src = skillSourceDir()
  const dest = join(cwd, '.agents/skills/joinquest-integration')
  copyTree(src, dest, { dryRun })
  return dest
}

export function skillPathInProject(cwd) {
  return join(cwd, '.agents/skills/joinquest-integration')
}
