import {
  LAUNCH_GAME,
  LEAVE_MATCH,
  LOOKING_FOR_GROUP,
  LOOK_FOR_GROUP,
  STOP_LOOKING,
  joinAsLabel,
  waitingAsRoleLine,
} from '../../lib/playerCopy'

export default function GameQueueActions({
  joinOptions,
  queueState,
  joinUrl,
  busy,
  selectedQueuePath = '',
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

  if (joinOptions?.kind === 'composition') {
    return (
      <section className="join-group-panel" aria-label={LOOK_FOR_GROUP}>
        <p className="join-group-panel__label">{LOOK_FOR_GROUP}</p>
        {queueState === 'waiting' ? (
          <div className="join-group-panel__actions">
            <p className="join-group-panel__status">
              {selectedQueuePath ? waitingAsRoleLine(selectedQueuePath) : LOOKING_FOR_GROUP}
            </p>
            <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
              {STOP_LOOKING}
            </button>
          </div>
        ) : (
          <div className="join-group-panel__actions">
            {joinOptions.paths.map((path) => (
              <button
                key={path}
                type="button"
                className="game-list-button"
                onClick={() => onJoin(path)}
                disabled={busy || disabled}
              >
                {busy ? LOOKING_FOR_GROUP : joinAsLabel(path)}
              </button>
            ))}
          </div>
        )}
      </section>
    )
  }

  return (
    <div className="game-list-actions">
      {queueState === 'waiting' ? (
        <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
          {STOP_LOOKING}
        </button>
      ) : (
        <button type="button" className="game-list-button" onClick={() => onJoin()} disabled={busy || disabled}>
          {busy ? LOOKING_FOR_GROUP : LOOK_FOR_GROUP}
        </button>
      )}
    </div>
  )
}
