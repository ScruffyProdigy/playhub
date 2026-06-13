import assert from 'node:assert/strict'
import { test } from 'node:test'
import { graphqlRequest, MUTATIONS } from '../src/graphql.js'

test('registerMyGame mutation includes credentials and connection status', () => {
  assert.match(MUTATIONS.registerMyGame, /registerMyGame/)
  assert.match(MUTATIONS.registerMyGame, /serviceToken/)
  assert.match(MUTATIONS.registerMyGame, /connectError/)
})

test('graphqlRequest sends registerMyGame mutation', async () => {
  const originalFetch = globalThis.fetch
  let capturedBody

  globalThis.fetch = async (_url, init) => {
    capturedBody = JSON.parse(init.body)
    return {
      ok: true,
      json: async () => ({
        data: {
          registerMyGame: {
            connected: true,
            connectError: null,
            webhookSecret: 'whsec_test',
            serviceToken: 'svc_test',
            game: {
              id: 'game-1',
              slug: 'my-game',
              name: 'My Game',
              visibility: 'PRIVATE_TESTING',
              apiBaseUrl: 'https://mygame.example.com',
            },
          },
        },
      }),
    }
  }

  try {
    const data = await graphqlRequest(
      'https://joinquest.cc/graphql',
      { authHeader: 'Bearer lq_dev_test' },
      MUTATIONS.registerMyGame,
      {
        input: {
          name: 'My Game',
          slug: 'my-game',
          shortDescription: 'A fun game',
          apiBaseUrl: 'https://mygame.example.com',
          contactEmail: 'dev@example.com',
        },
      },
    )

    assert.equal(capturedBody.variables.input.slug, 'my-game')
    assert.equal(data.registerMyGame.game.id, 'game-1')
    assert.equal(data.registerMyGame.serviceToken, 'svc_test')
  } finally {
    globalThis.fetch = originalFetch
  }
})
