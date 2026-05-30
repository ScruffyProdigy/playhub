export default function GameQueueActions({
  queueState,
  joinUrl,
  busy,
  onJoin,
  onLeave,
}) {
  if (queueState === 'matched' && joinUrl) {
    return (
      <div className="game-list-actions">
        <a className="game-list-button" href={joinUrl}>
          Launch game
        </a>
        <button type="button" className="game-list-button game-list-button-secondary" onClick={onLeave} disabled={busy}>
          Leave match
        </button>
      </div>
    )
  }

  return (
    <div className="game-list-actions">
      {queueState === 'waiting' ? (
        <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
          Leave queue
        </button>
      ) : (
        <button type="button" className="game-list-button" onClick={onJoin} disabled={busy}>
          {busy ? 'Joining…' : 'Join queue'}
        </button>
      )}
    </div>
  )
}
