import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import DeveloperMcpWizard from './DeveloperMcpWizard'
import {
  buildClaudeMcpAddCommand,
  buildInstallDevCommand,
  buildMcpServerConfig,
  createDeveloperApiKey,
  fetchMyDeveloperApiKeys,
} from '../../lib/developers'

vi.mock('../../lib/developers', () => ({
  fetchMyDeveloperApiKeys: vi.fn().mockResolvedValue([]),
  createDeveloperApiKey: vi.fn(),
  buildMcpServerConfig: vi.fn(({ clientId, apiKey }) => ({
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
  })),
  buildClaudeMcpAddCommand: vi.fn(
    ({ apiKey }) => `claude mcp add --env JOINQUEST_API_KEY=${apiKey} ...`,
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
    expect(buildInstallDevCommand).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_test_secret', client: 'cursor' }),
    )
    expect(screen.getByRole('button', { name: /copy install command/i })).toBeInTheDocument()
  })

  it('fills install command when an existing key is pasted', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.type(screen.getByLabelText(/your api key/i), 'lq_dev_saved_key_abc')

    expect(buildInstallDevCommand).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_saved_key_abc', client: 'cursor' }),
    )
    expect(screen.getByText(/export JOINQUEST_API_KEY=lq_dev_saved_key_abc/)).toBeInTheDocument()
  })
})
