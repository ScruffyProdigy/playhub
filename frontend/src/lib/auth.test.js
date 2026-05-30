import { describe, it, expect, vi, beforeEach } from 'vitest'

vi.mock('./graphql', () => ({
  graphqlRequest: vi.fn(),
}))

import { graphqlRequest } from './graphql'
import { completeMagicLoginOnce } from './auth'

describe('completeMagicLoginOnce', () => {
  beforeEach(() => {
    vi.mocked(graphqlRequest).mockReset()
  })

  it('deduplicates concurrent completions for the same token', async () => {
    vi.mocked(graphqlRequest).mockResolvedValue({
      completeMagic: { id: 'user-1', email: 'player@example.com' },
    })

    const [first, second] = await Promise.all([
      completeMagicLoginOnce('token-123'),
      completeMagicLoginOnce('token-123'),
    ])

    expect(first).toEqual({ id: 'user-1', email: 'player@example.com' })
    expect(second).toEqual(first)
    expect(graphqlRequest).toHaveBeenCalledTimes(1)
  })

  it('rejects an empty token', async () => {
    await expect(completeMagicLoginOnce('   ')).rejects.toThrow('Missing sign-in token')
    expect(graphqlRequest).not.toHaveBeenCalled()
  })
})
