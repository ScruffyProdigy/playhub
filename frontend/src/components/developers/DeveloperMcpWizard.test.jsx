import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import DeveloperMcpWizard from './DeveloperMcpWizard'
import {
  buildClaudeMcpAddCommand,
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
}))

describe('DeveloperMcpWizard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders expanded setup without crashing', async () => {
    const user = userEvent.setup()
    render(<DeveloperMcpWizard />)

    await user.click(screen.getByRole('button', { name: /show setup/i }))

    expect(screen.getByRole('button', { name: /generate api key/i })).toBeInTheDocument()
    expect(screen.getAllByText(/@joinquest\/mcp-integration/).length).toBeGreaterThan(0)
  })

  it('shows one-click copy buttons after generating an API key', async () => {
    createDeveloperApiKey.mockResolvedValue({
      secret: 'lq_dev_test_secret',
      apiKey: { id: '1', keyPrefix: 'lq_dev_test', lastUsedAt: null },
    })

    const user = userEvent.setup()
    render(<DeveloperMcpWizard defaultExpanded />)

    await user.click(screen.getByRole('button', { name: /generate api key/i }))

    const quickCopy = document.querySelector('.developer-mcp__quick-copy')
    expect(quickCopy).not.toBeNull()
    expect(within(quickCopy).getByRole('button', { name: /copy cursor json/i })).toBeInTheDocument()
    expect(within(quickCopy).getByRole('button', { name: /copy claude command/i })).toBeInTheDocument()
    expect(buildMcpServerConfig).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_test_secret', clientId: 'cursor' }),
    )
    expect(buildClaudeMcpAddCommand).toHaveBeenCalledWith(
      expect.objectContaining({ apiKey: 'lq_dev_test_secret' }),
    )
  })
})
