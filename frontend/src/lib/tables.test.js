import { describe, it, expect } from 'vitest'
import {
  enrichTableSeats,
  groupSeatSlotsByTeam,
  groupSeatSlotsForDisplay,
  mergeTableRecord,
  mySeatDisplayName,
  mySeatKeyOnTable,
  seatLabelInSection,
  sectionTitleForSeat,
} from './tables'

describe('groupSeatSlotsByTeam', () => {
  it('puts fifo seats in ungrouped only', () => {
    const slots = [
      { seatKey: '1', displayName: '1' },
      { seatKey: '2', displayName: '2' },
    ]
    const { teams, ungrouped } = groupSeatSlotsByTeam(slots)
    expect(teams).toHaveLength(0)
    expect(ungrouped).toHaveLength(2)
  })

  it('groups team-prefixed seats separately', () => {
    const slots = [
      { seatKey: 'Team-1-DPS-1', displayName: 'DPS · 1', teamPrefix: 'Team-1' },
      { seatKey: 'Team-2-DPS-1', displayName: 'DPS · 1', teamPrefix: 'Team-2' },
    ]
    const { teams, ungrouped } = groupSeatSlotsByTeam(slots)
    expect(teams).toHaveLength(2)
    expect(ungrouped).toHaveLength(0)
  })
})

describe('groupSeatSlotsForDisplay', () => {
  it('groups composition roles into sections', () => {
    const layout = groupSeatSlotsForDisplay([
      { seatKey: 'ClueGiver-Red', displayName: 'Clue Giver · Red', queuePath: 'ClueGiver' },
      { seatKey: 'ClueGiver-Blue', displayName: 'Clue Giver · Blue', queuePath: 'ClueGiver' },
      { seatKey: 'Guesser-1', displayName: 'Guesser · 1', queuePath: 'Guesser' },
    ])
    expect(layout.kind).toBe('roles')
    expect(layout.roles).toHaveLength(2)
    expect(layout.roles[0][0]).toBe('Clue Giver')
  })
})

describe('seatLabelInSection', () => {
  it('shortens clue giver seats to color only under section header', () => {
    expect(
      seatLabelInSection({ displayName: 'Clue Giver · Red' }, 'Clue Giver'),
    ).toBe('Red')
  })

  it('hides numbered guesser seats under section header', () => {
    expect(seatLabelInSection({ displayName: 'Guesser · 1' }, 'Guesser')).toBeNull()
    expect(seatLabelInSection({ displayName: 'Guesser · 2' }, 'Guesser')).toBeNull()
  })
})

describe('sectionTitleForSeat', () => {
  const wordHuntSlots = [
    { seatKey: 'x', displayName: 'Clue Giver · Red', queuePath: 'ClueGiver' },
    { seatKey: 'y', displayName: 'Guesser · 1', queuePath: 'Guesser' },
  ]

  it('finds role section title for a seat key', () => {
    expect(sectionTitleForSeat('x', wordHuntSlots)).toBe('Clue Giver')
    expect(sectionTitleForSeat('y', wordHuntSlots)).toBe('Guesser')
  })
})

describe('mergeTableRecord', () => {
  it('replaces seats when update includes an empty list', () => {
    const prev = {
      id: 't1',
      seats: [{ seatKey: '1', user: { id: 'u1' } }],
      seatSlots: [{ seatKey: '1', displayName: '1', user: { id: 'u1' } }],
    }
    const merged = mergeTableRecord(prev, { id: 't1', seats: [], seatSlots: [{ seatKey: '1', displayName: '1' }] })
    expect(merged.seats).toEqual([])
    expect(merged.seatSlots[0].user).toBeUndefined()
  })
})

describe('enrichTableSeats', () => {
  it('fills slot users from table.seats when missing on slots', () => {
    const enriched = enrichTableSeats({
      seats: [{ seatKey: '1', user: { id: 'u1', displayName: 'Pat' } }],
      seatSlots: [{ seatKey: '1', displayName: 'Seat 1' }],
    })
    expect(enriched.seatSlots[0].user?.displayName).toBe('Pat')
  })
})

describe('mySeatKeyOnTable', () => {
  it('finds seat from seats array', () => {
    const key = mySeatKeyOnTable(
      { seats: [{ seatKey: 'a', user: { id: 'u1' } }] },
      'u1',
    )
    expect(key).toBe('a')
  })

  it('falls back to seatSlots when seats missing', () => {
    const key = mySeatKeyOnTable(
      { seatSlots: [{ seatKey: 'b', user: { id: 'u2' } }] },
      'u2',
    )
    expect(key).toBe('b')
  })
})

describe('mySeatDisplayName', () => {
  it('returns shortened seat label for role sections', () => {
    const name = mySeatDisplayName(
      {
        seats: [{ seatKey: 'x', user: { id: 'u1' } }],
        seatSlots: [
          { seatKey: 'x', displayName: 'Clue Giver · Red', queuePath: 'ClueGiver' },
          { seatKey: 'y', displayName: 'Guesser · 1', queuePath: 'Guesser' },
        ],
      },
      'u1',
    )
    expect(name).toBe('Clue Giver · Red')
  })

  it('returns role name only for numbered guesser seats', () => {
    const name = mySeatDisplayName(
      {
        seats: [{ seatKey: 'y', user: { id: 'u1' } }],
        seatSlots: [
          { seatKey: 'x', displayName: 'Clue Giver · Red', queuePath: 'ClueGiver' },
          { seatKey: 'y', displayName: 'Guesser · 2', queuePath: 'Guesser' },
        ],
      },
      'u1',
    )
    expect(name).toBe('Guesser')
  })
})
