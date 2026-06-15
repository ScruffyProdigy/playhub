import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import DeveloperMcpWizard from './DeveloperMcpWizard'
import {
  buildClaudeMcpAddCommand,
  buildInstallClaudePluginCommand,
  buildInstallCursorPluginCommand,
  buildInstallDevCommand,
  buildMcpServerConfig,
  createDeveloperApiKey,
  fetchMyDeveloperApiKeys,
} from '../../lib/developers'

vi.mock('../../lib/developers', () => ({
  fetchMyDeveloperApiKeys: vi.fn().mockResolvedValue([]),
  createDeveloperApiKey: vi.fn(),
  buildMcpServerConfig: vi.fn(({ clientId, apiKey }) =>
    clientId === 'copilot'
      ? {
          servers: {
            'joinquest-integration': {
              type: 'stdio',
              command: 'npx',
              args: ['-y', '@joinquest/mcp-integration'],
              env: { JOINQUEST_API_KEY: apiKey },
            },
          },
        }
      : {
          mcpServers: {
            'joinquest-integration': {
              type: 'stdio',
              command: 'npx',
              args:
                clientId === 'cursor'
                  ? ['--yes', '--package', '@joinquest/mcp-integration', 'joinquest-integration-mcp-cursor']
                  : ['-y', '@joinquest/mcp-integration'],
              env: { JOINQUEST_API_KEY: apiKey },
            },
          },
        },
  ),
  buildClaudeMcpAddCommand: vi.fn(
    ({ apiKey }) => `claude mcp add --env JOINQUEST_API_KEY=${apiKey} ...`,
  ),
  buildInstallCursorPluginCommand: vi.fn(
    ({ apiKey }) =>
      `export JOINQUEST_API_KEY=${apiKey || 'lq_dev_PASTE_YOUR_KEY'}\ncurl -fsSL .../install-joinquest-cursor-plugin.sh | sh`,
  ),
  buildInstallClaudePluginCommand: vi.fn(
    ({ apiKey }) =>
      `export JOINQUEST_API_KEY=${apiKey || 'lq_dev_PASTE_YOUR_KEY'}\ncurl -fsSL .../install-joinquest-claude-plugin.sh | sh`,
  ),
  buildInstallDevCommand: vi.fn(
    ({ apiKey, client }) =>
      client === 'skill-only'
        ? 'curl -fsSL .../install-joinquest-dev.sh | sh -s -- --skill-only'
        : `export JOINQUEST_API_KEY=${apiKey || 'lq_dev_PASTE_YOUR_KEY'}\ncurl -fsSL .../install-joinquest-dev.sh | sh -s -- --${client}`,
  ),
  buildInstallDevInspectCommand: vi.fn(
    ({ apiKey }) => `curl -fsSL ... -o install-joinquest-dev.sh\nless install-joinquest-dev.sh`,
  ),
  INSTALL_DEV_SCRIPT_GITHUB: 'https://github.com/example/install-joinquest-dev.sh',
  INSTALL_CURSOR_PLUGIN_SCRIPT_GITHUB: 'https://github.com/example/install-joinquest-cursor-plugin.sh',
  INSTALL_CLAUDE_PLUGIN_SCRIPT_GITHUB: 'https://github.com/example/install-joinquest-claude-plugin.sh',
  CURSOR_PLUGIN_GITHUB: 'https://github.com/example/plugins/joinquest',
}))

describe('DeveloperMcpWizard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders expanded setup without crashing', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard />)

    await user.click(screen.getByRole('button', { name: /show setup/i }))

    expect(screen.getByRole('tablist', { name: /which editor do you use/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /generate new api key/i })).toBeInTheDocument()
    expect(screen.getByLabelText(/your api key/i)).toBeInTheDocument()
    expect(screen.getAllByText(/@joinquest\/mcp-integration/).length).toBeGreaterThan(0)
  })

  it('shows one-click copy buttons after generating an API key', async () => {
    createDeveloperApiKey.mockResolvedValue({
      secret: 'lq_dev_test_secret',
      apiKey: { id: '1', keyPrefix: 'lq_dev_test', lastUsedAt: null },
    })

    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.click(screen.getByRole('button', { name: /generate new api key/i }))

    expect(screen.getByLabelText(/your api key/i)).toHaveValue('lq_dev_test_secret')
    expect(buildMcpServerConfig).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_test_secret', clientId: 'cursor' }),
    )
    expect(buildInstallCursorPluginCommand).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_test_secret' }),
    )
    expect(screen.getByRole('button', { name: /copy plugin install command/i })).toBeInTheDocument()
  })

  it('fills plugin install command when an existing key is pasted', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.type(screen.getByLabelText(/your api key/i), 'lq_dev_saved_key_abc')

    expect(buildInstallCursorPluginCommand).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_saved_key_abc' }),
    )
    expect(screen.getByText(/install-joinquest-cursor-plugin\.sh/)).toBeInTheDocument()
  })

  it('shows Claude Code plugin install after switching tabs', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.click(screen.getByRole('tab', { name: 'Claude Code' }))

    expect(buildInstallClaudePluginCommand).toHaveBeenCalled()
    expect(screen.getByRole('button', { name: /copy plugin install command/i })).toBeInTheDocument()
  })

  it('shows Copilot install command after switching tabs', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.click(screen.getByRole('tab', { name: 'GitHub Copilot' }))

    expect(buildInstallDevCommand).toHaveBeenCalledWith(
      expect.objectContaining({ client: 'copilot' }),
    )
    expect(screen.getByText(/--copilot/)).toBeInTheDocument()
  })

  it('explains ChatGPT is unsupported', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.click(screen.getByRole('tab', { name: 'ChatGPT' }))

    expect(screen.getByText(/ChatGPT is not supported yet/i)).toBeInTheDocument()
    expect(screen.queryByRole('button', { name: /generate new api key/i })).not.toBeInTheDocument()
  })
})
