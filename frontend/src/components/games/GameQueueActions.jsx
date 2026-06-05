import {
  LAUNCH_GAME,
  LEAVE_MATCH,
  LOOKING_FOR_GROUP,
  LOOK_FOR_GROUP,
  STOP_LOOKING,
  joinAsLabel,
  waitingAsRoleLine,
} from '../../lib/playerCopy'

function pathOption(path) {
  if (typeof path === 'string') {
    return { queuePath: path, displayName: path }
  }
  return {
    queuePath: path.queuePath,
    displayName: path.displayName || path.queuePath,
  }
}

function JoinPathButton({ path, busy, disabled, onJoin }) {
  const { queuePath, displayName } = pathOption(path)
  return (
    <button
      type="button"
      className="game-list-button"
      onClick={() => onJoin(queuePath)}
      disabled={busy || disabled}
    >
      {busy ? LOOKING_FOR_GROUP : joinAsLabel(displayName)}
    </button>
  )
}

function JoinGroupPanel({ children }) {
  return (
    <section className="join-group-panel" aria-label={LOOK_FOR_GROUP}>
      <p className="join-group-panel__label">{LOOK_FOR_GROUP}</p>
      <div className="join-group-panel__actions">{children}</div>
    </section>
  )
}

function FifoQueueActions({ queueState, busy, disabled, onJoin, onLeave }) {
  if (queueState === 'waiting') {
    return (
      <div className="game-list-actions game-list-actions--stack">
        <p className="queue-status-line">{LOOKING_FOR_GROUP}</p>
        <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
          {STOP_LOOKING}
        </button>
      </div>
    )
  }

  return (
    <div className="game-list-actions">
      <button
        type="button"
        className="game-list-button"
        onClick={() => onJoin()}
        disabled={busy || disabled}
      >
        {busy ? LOOKING_FOR_GROUP : LOOK_FOR_GROUP}
      </button>
    </div>
  )
}

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

  const compositionPaths =
    joinOptions?.kind === 'composition'
      ? joinOptions.paths.filter((path) => (pathOption(path).queuePath?.trim() ?? '') !== '')
      : []

  const selectedDisplayName = compositionPaths
    .map(pathOption)
    .find(({ queuePath }) => queuePath === selectedQueuePath)?.displayName

  if (compositionPaths.length > 0) {
    const alternatePaths = compositionPaths.filter(
      (path) => pathOption(path).queuePath !== selectedQueuePath,
    )

    return (
      <JoinGroupPanel>
        {queueState === 'waiting' ? (
          <>
            <p className="queue-status-line">
              {selectedQueuePath
                ? waitingAsRoleLine(selectedDisplayName || selectedQueuePath)
                : LOOKING_FOR_GROUP}
            </p>
            {alternatePaths.map((path) => (
              <JoinPathButton
                key={pathOption(path).queuePath}
                path={path}
                busy={busy}
                disabled={disabled}
                onJoin={onJoin}
              />
            ))}
            <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
              {STOP_LOOKING}
            </button>
          </>
        ) : (
          compositionPaths.map((path) => (
            <JoinPathButton
              key={pathOption(path).queuePath}
              path={path}
              busy={busy}
              disabled={disabled}
              onJoin={onJoin}
            />
          ))
        )}
      </JoinGroupPanel>
    )
  }

  return <FifoQueueActions queueState={queueState} busy={busy} disabled={disabled} onJoin={onJoin} onLeave={onLeave} />
}
