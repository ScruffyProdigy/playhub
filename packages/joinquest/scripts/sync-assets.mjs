#!/usr/bin/env node
/** Copy skill + plugin trees into the published package. Run from repo root in prepublishOnly. */
import { cpSync, mkdirSync, rmSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'

const pkgRoot = join(dirname(fileURLToPath(import.meta.url)), '..')
const repoRoot = join(pkgRoot, '../..')
const assetsDir = join(pkgRoot, 'assets')

const copies = [
  {
    from: join(repoRoot, '.agents/skills/joinquest-integration'),
    to: join(assetsDir, 'skill'),
  },
  {
    from: join(repoRoot, 'plugins/joinquest'),
    to: join(assetsDir, 'plugin'),
  },
]

rmSync(assetsDir, { recursive: true, force: true })
mkdirSync(assetsDir, { recursive: true })

for (const { from, to } of copies) {
  cpSync(from, to, { recursive: true })
}

console.log('Synced joinquest CLI assets (skill + plugin).')
