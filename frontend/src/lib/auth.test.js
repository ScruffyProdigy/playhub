import { afterEach, describe, expect, it, vi } from 'vitest'

vi.mock('./graphql', () => ({
  graphqlRequest: vi.fn(),
}))

import { graphqlRequest } from './graphql'
import { completeSignInWithLinkOnce } from './auth'

describe('completeSignInWithLinkOnce', () => {
  afterEach(() => {
    vi.mocked(graphqlRequest).mockReset()
  })

  it('deduplicates concurrent completions for the same token', async () => {
    vi.mocked(graphqlRequest).mockResolvedValue({
      completeSignInWithLink: { id: 'user-1', email: 'player@example.com' },
    })

    const [first, second] = await Promise.all([
      completeSignInWithLinkOnce('token-123'),
      completeSignInWithLinkOnce('token-123'),
    ])

    expect(first).toEqual({ id: 'user-1', email: 'player@example.com' })
    expect(second).toEqual(first)
    expect(graphqlRequest).toHaveBeenCalledTimes(1)
  })

  it('rejects empty tokens', async () => {
    await expect(completeSignInWithLinkOnce('   ')).rejects.toThrow('Missing sign-in token')
  })
})
