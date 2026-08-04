import { cpSync, existsSync, readFileSync, rmSync, writeFileSync } from 'node:fs'
import { join } from 'node:path'
import { pluginSourceDir } from './assets.js'
import { claudePluginDir, cursorPluginDir } from './paths.js'

const API_KEY_PLACEHOLDER = '<paste-api-key-from-joinquest-dashboard>'

function patchPluginApiKey(pluginDir, apiKey) {
  for (const name of ['mcp.json', '.mcp.json']) {
    const path = join(pluginDir, name)
    if (!existsSync(path)) continue
    const text = readFileSync(path, 'utf8').replaceAll(API_KEY_PLACEHOLDER, apiKey)
    writeFileSync(path, text)
  }
}

export function installCursorPlugin(apiKey, { dryRun = false } = {}) {
  const dest = cursorPluginDir()
  if (!dryRun) {
    const src = pluginSourceDir()
    rmSync(dest, { recursive: true, force: true })
    cpSync(src, dest, { recursive: true })
    patchPluginApiKey(dest, apiKey)
  }
  return dest
}

export function installClaudePlugin(apiKey, { dryRun = false } = {}) {
  const dest = claudePluginDir()
  if (!dryRun) {
    const src = pluginSourceDir()
    rmSync(dest, { recursive: true, force: true })
    cpSync(src, dest, { recursive: true })
    patchPluginApiKey(dest, apiKey)
  }
  return dest
}
