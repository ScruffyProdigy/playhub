import {
  LAUNCH_GAME,
  LEAVE_MATCH,
  LOOKING_FOR_GROUP,
  LOOK_FOR_GROUP,
  STOP_LOOKING,
} from '../../lib/playerCopy'

export default function GameQueueActions({
  queueState,
  joinUrl,
  busy,
  onJoin,
  onLeave,
  disabled = false,
}) {
  if (queueState === 'matched' && joinUrl) {
    return (
      <div className="game-list-actions">
        <a className="game-list-button" href={joinUrl}>
          {LAUNCH_GAME}
        </a>
        <button type="button" className="game-list-button game-list-button-secondary" onClick={onLeave} disabled={busy}>
          {LEAVE_MATCH}
        </button>
      </div>
    )
  }

  return (
    <div className="game-list-actions">
      {queueState === 'waiting' ? (
        <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
          {STOP_LOOKING}
        </button>
      ) : (
        <button type="button" className="game-list-button" onClick={onJoin} disabled={busy || disabled}>
          {busy ? LOOKING_FOR_GROUP : LOOK_FOR_GROUP}
        </button>
      )}
    </div>
  )
}
