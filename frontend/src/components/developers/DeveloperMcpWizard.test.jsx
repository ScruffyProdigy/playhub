import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import DeveloperMcpWizard from './DeveloperMcpWizard'

vi.mock('../../lib/developers', () => ({
  fetchMyDeveloperApiKeys: vi.fn().mockResolvedValue([]),
  createDeveloperApiKey: vi.fn(),
  buildMcpServerConfig: vi.fn(() => ({
    mcpServers: {
      'joinquest-integration': {
        command: 'npx',
        args: ['-y', '@joinquest/mcp-integration'],
        env: {},
      },
    },
  })),
}))

vi.mock('../../lib/env', () => ({
  getGraphQLUrl: () => 'https://joinquest.cc/graphql',
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
    expect(screen.getByText(/npx -y @joinquest\/mcp-integration/)).toBeInTheDocument()
  })
})
