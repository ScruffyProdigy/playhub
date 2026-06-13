import test from 'node:test'
import assert from 'node:assert/strict'
import { buildExampleProvisionPayload } from '../src/example-provision.js'

test('buildExampleProvisionPayload uses first mode seats', () => {
  const payload = buildExampleProvisionPayload({
    game: {
      id: 'a1000000-0000-4000-8000-000000000001',
      modes: [{
        modeKey: 'duel',
        minPlayers: 2,
        seats: [{ seatKey: '1' }, { seatKey: '2' }],
      }],
    },
    credentials: { serviceToken: 'v1.test.token' },
    lobbyIssuer: 'https://joinquest.cc',
    lobbyReturnUrl: 'https://joinquest.cc/return',
    lobbyGraphqlUrl: 'https://joinquest.cc/graphql',
  })

  assert.equal(payload.assignment.gameMode, 'duel')
  assert.equal(payload.assignment.seats.length, 2)
  assert.equal(payload.lobby.serviceToken, 'v1.test.token')
})
