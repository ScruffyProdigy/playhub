import {
  LAUNCH_GAME,
  LEAVE_GAME,
  LEAVE_TABLE_SEAT,
  STOP_LOOKING,
  bannerIntentPlayingHint,
  bannerIntentLaunchPendingHint,
  bannerIntentWaitingHint,
  bannerLiveUpdatesPausedHint,
  bannerTableBackfillHint,
  bannerTableSeatHint,
  bannerTableSeatLine,
  bannerWaitingLine,
} from '../../lib/playerCopy'
import {
  hasFormingTableIntent,
  hasReadyToPlayIntent,
  hasWaitingIntent,
  playingIntentTitle,
  resolveIntentLaunchUrl,
} from '../../lib/intent'

export default function IntentBanner({ activeIntent, activeTableSeat, busy, liveUpdatesConnected = true, onLeave }) {
  if (hasWaitingIntent(activeIntent)) {
    return (
      <aside className="intent-banner" role="region" aria-live="polite" aria-label="Your intent">
        <div className="intent-banner__copy">
          <p className="intent-banner__title">
            {bannerWaitingLine(
              activeIntent.gameName,
              activeIntent.queuedCount,
              activeIntent.queuePathDisplayName,
              activeIntent.formingGaps,
            )}
          </p>
          <p className="intent-banner__hint">
            {liveUpdatesConnected ? bannerIntentWaitingHint() : bannerLiveUpdatesPausedHint()}
          </p>
        </div>
        <div className="intent-banner__actions">
          <button type="button" className="game-list-button" onClick={onLeave} disabled={busy}>
            {busy ? '…' : STOP_LOOKING}
          </button>
        </div>
      </aside>
    )
  }

  if (hasReadyToPlayIntent(activeIntent, activeTableSeat)) {
    const launchUrl = resolveIntentLaunchUrl(activeIntent, activeTableSeat)
    const title = playingIntentTitle(activeIntent, activeTableSeat)

    return (
      <aside className="intent-banner" role="region" aria-live="polite" aria-label="Your intent">
        <div className="intent-banner__copy">
          <p className="intent-banner__title">{title}</p>
          <p className="intent-banner__hint">
            {launchUrl ? bannerIntentPlayingHint() : bannerIntentLaunchPendingHint()}
          </p>
        </div>
        <div className="intent-banner__actions">
          {launchUrl ? (
            <a className="intent-banner__cta game-list-button" href={launchUrl}>
              {LAUNCH_GAME}
            </a>
          ) : null}
          <button
            type="button"
            className="game-list-button game-list-button-secondary"
            onClick={onLeave}
            disabled={busy}
          >
            {busy ? '…' : LEAVE_GAME}
          </button>
        </div>
      </aside>
    )
  }

  if (!hasFormingTableIntent(activeIntent, activeTableSeat)) {
    return null
  }

  return (
    <aside className="intent-banner" role="region" aria-live="polite" aria-label="Your intent">
      <div className="intent-banner__copy">
        <p className="intent-banner__title">
          {bannerTableSeatLine(
            activeTableSeat.gameName,
            activeTableSeat.modeName,
            activeTableSeat.seatDisplayName,
          )}
        </p>
        <p className="intent-banner__hint">
          {activeTableSeat.backfillActive
            ? bannerTableBackfillHint(activeTableSeat.formingGaps)
            : bannerTableSeatHint()}
        </p>
      </div>
      <div className="intent-banner__actions">
        <button type="button" className="game-list-button game-list-button-secondary" onClick={onLeave} disabled={busy}>
          {busy ? '…' : LEAVE_TABLE_SEAT}
        </button>
      </div>
    </aside>
  )
}
