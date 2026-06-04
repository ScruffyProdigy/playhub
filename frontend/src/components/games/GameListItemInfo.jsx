import { waitingForGroupLine } from '../../lib/playerCopy'

export default function GameListItemInfo({ game, queueState, queuedCount, error, notice }) {
  const sessionCount = game.activeSessions?.length ?? 0

  return (
    <div className="game-list-copy">
      <h3>{game.name}</h3>
      <p className="game-list-meta">
        {sessionCount} active session{sessionCount === 1 ? '' : 's'}
      </p>
      {queueState === 'waiting' ? (
        <p className="game-list-meta" role="status">
          {waitingForGroupLine(queuedCount)}
        </p>
      ) : null}
      {notice ? (
        <p className="status-message" role="status">
          {notice}
        </p>
      ) : null}
      {error ? (
        <p className="status-message status-message-error" role="status">
          {error}
        </p>
      ) : null}
    </div>
  )
}
