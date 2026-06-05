import { describe, it, expect } from 'vitest'
import {
  defaultModeForGame,
  joinGroupOptionsForGame,
  joinGroupOptionsForMode,
} from './games'

describe('joinGroupOptionsForMode', () => {
  it('returns single-bucket (fifo) when all seats have empty queue paths', () => {
    const mode = {
      seats: [{ queuePath: null }, { queuePath: '' }],
    }
    expect(joinGroupOptionsForMode(mode)).toEqual({ kind: 'fifo', paths: [] })
  })

  it('returns sorted composition paths', () => {
    const mode = {
      seats: [
        { queuePath: 'Tank' },
        { queuePath: 'DPS' },
        { queuePath: 'DPS' },
        { queuePath: 'Support' },
      ],
    }
    expect(joinGroupOptionsForMode(mode)).toEqual({
      kind: 'composition',
      paths: ['DPS', 'Support', 'Tank'],
    })
  })
})

describe('joinGroupOptionsForGame', () => {
  it('uses the default active mode', () => {
    const game = {
      modes: [
        {
          queues: [{ status: 'inactive' }],
          seats: [{ queuePath: 'DPS' }],
        },
        {
          queues: [{ status: 'active' }],
          seats: [{ queuePath: 'Tank' }, { queuePath: 'DPS' }],
        },
      ],
    }
    expect(defaultModeForGame(game).queues[0].status).toBe('active')
    expect(joinGroupOptionsForGame(game)).toEqual({
      kind: 'composition',
      paths: ['DPS', 'Tank'],
    })
  })
})
