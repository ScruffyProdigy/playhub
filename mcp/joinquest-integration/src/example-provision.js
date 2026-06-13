const CHECK_USER_ONE = 'a0000000-0000-4000-8000-000000000001'
const CHECK_USER_TWO = 'a0000000-0000-4000-8000-000000000002'

export function buildExampleProvisionPayload({ game, credentials, lobbyIssuer, lobbyReturnUrl, lobbyGraphqlUrl }) {
  const mode = game?.modes?.[0]
  if (!mode) {
    throw new Error('Game has no synced modes — connect API and sync manifest first.')
  }

  const seats = mode.seats ?? []
  const seatCount = Math.min(2, seats.length)
  if (seatCount === 0) {
    throw new Error('Game mode has no expanded seats.')
  }

  const externalMatchId = 'example-match-' + game.id.slice(0, 8)
  const assignmentSeats = seats.slice(0, seatCount).map((seat, index) => ({
    seatKey: seat.seatKey,
    lobbyUserId: index === 0 ? CHECK_USER_ONE : CHECK_USER_TWO,
    displayName: index === 0 ? 'Player One' : 'Player Two',
  }))

  return {
    lobbyId: lobbyIssuer,
    lobby: {
      returnUrl: lobbyReturnUrl,
      graphqlUrl: lobbyGraphqlUrl,
      serviceToken: credentials.serviceToken,
    },
    assignment: {
      externalMatchId,
      gameMode: mode.modeKey,
      seats: assignmentSeats,
    },
  }
}
