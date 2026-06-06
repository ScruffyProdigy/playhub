import {
  LAUNCH_GAME,
  LEAVE_MATCH,
  LEAVE_TABLE_SEAT,
  STOP_LOOKING,
  bannerMatchedLine,
  bannerTableSeatHint,
  bannerTableSeatLine,
  bannerWaitingLine,
} from '../../lib/playerCopy'

export default function ActiveQueueBanner({ activeQueue, activeTableSeat, busy, onLeave }) {
  if (activeQueue?.queueId) {
    const isMatched = activeQueue.status === 'MATCHED'

    return (
      <aside className="active-queue-banner" role="region" aria-live="polite" aria-label="Your group">
        <div className="active-queue-banner__copy">
          <p className="active-queue-banner__title">
            {isMatched
              ? bannerMatchedLine(activeQueue.gameName)
              : bannerWaitingLine(
                  activeQueue.gameName,
                  activeQueue.queuedCount,
                  activeQueue.queuePathDisplayName,
                )}
          </p>
          {isMatched ? (
            <p className="active-queue-banner__hint">Launch when you are ready — your seat is reserved.</p>
          ) : (
            <p className="active-queue-banner__hint">We will notify you here when your group is ready.</p>
          )}
        </div>
        <div className="active-queue-banner__actions">
          {isMatched && activeQueue.joinUrl ? (
            <a className="active-queue-banner__cta game-list-button" href={activeQueue.joinUrl}>
              {LAUNCH_GAME}
            </a>
          ) : null}
          <button
            type="button"
            className={`game-list-button${isMatched ? ' game-list-button-secondary' : ''}`}
            onClick={onLeave}
            disabled={busy}
          >
            {busy ? '…' : isMatched ? LEAVE_MATCH : STOP_LOOKING}
          </button>
        </div>
      </aside>
    )
  }

  if (!activeTableSeat?.tableId) {
    return null
  }

  return (
    <aside className="active-queue-banner" role="region" aria-live="polite" aria-label="Your table seat">
      <div className="active-queue-banner__copy">
        <p className="active-queue-banner__title">
          {bannerTableSeatLine(
            activeTableSeat.gameName,
            activeTableSeat.modeName,
            activeTableSeat.seatDisplayName,
          )}
        </p>
        <p className="active-queue-banner__hint">{bannerTableSeatHint()}</p>
      </div>
      <div className="active-queue-banner__actions">
        <button type="button" className="game-list-button game-list-button-secondary" onClick={onLeave} disabled={busy}>
          {busy ? '…' : LEAVE_TABLE_SEAT}
        </button>
      </div>
    </aside>
  )
}
