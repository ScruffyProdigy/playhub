import { homedir } from 'node:os'
import { join } from 'node:path'

export function claudeDesktopConfigPath() {
  switch (process.platform) {
    case 'darwin':
      return join(homedir(), 'Library/Application Support/Claude/claude_desktop_config.json')
    case 'win32':
      return join(process.env.APPDATA || join(homedir(), 'AppData/Roaming'), 'Claude/claude_desktop_config.json')
    default:
      return join(process.env.XDG_CONFIG_HOME || join(homedir(), '.config'), 'Claude/claude_desktop_config.json')
  }
}

export function windsurfMcpConfigPath() {
  return join(homedir(), '.codeium/windsurf/mcp_config.json')
}

export function clineMcpConfigPath() {
  switch (process.platform) {
    case 'darwin':
      return join(
        homedir(),
        'Library/Application Support/Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json',
      )
    case 'win32':
      return join(
        process.env.APPDATA || join(homedir(), 'AppData/Roaming'),
        'Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json',
      )
    default:
      return join(
        process.env.XDG_CONFIG_HOME || join(homedir(), '.config'),
        'Code/User/globalStorage/saoudrizwan.claude-dev/settings/cline_mcp_settings.json',
      )
  }
}

export function cursorPluginDir() {
  return join(homedir(), '.cursor/plugins/local/joinquest')
}

export function claudePluginDir() {
  return join(homedir(), '.claude/skills/joinquest-integration')
}
