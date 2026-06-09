import { useAuth } from '../auth/AuthProvider'
import {
  DISCARD,
  formatFormingGapsFromLobbyLine,
  KING_LABEL,
  LOOK_FOR_GROUP,
  START_GAME,
} from '../../lib/playerCopy'
import {
  countSeatedInGroup,
  displayName,
  enrichTableSeats,
  firstOpenSeatKey,
  formatGroupSeatCaption,
  groupSeatSlotsForDisplay,
  isKing,
  isPooledRoleGroup,
  mySeatDisplayName,
  mySeatKeyOnTable,
  queuePathMeta,
  seatLabelInSection,
} from '../../lib/tables'
import PlayerAvatar from '../avatars/PlayerAvatar'

function SeatRow({ slot, seatLabel, mySeat, userId, currentUser, kingUserId, busy, onSit }) {
  const taken = Boolean(slot.user)
  const isMine = mySeat === slot.seatKey || slot.user?.id === userId
  const showOccupant = taken || isMine
  const occupantUser = isMine ? currentUser ?? slot.user : slot.user
  const occupantLabel = isMine ? 'You' : displayName(slot.user)
  const isTableKing = Boolean(kingUserId && occupantUser?.id === kingUserId)

  return (
    <li
      className={`table-seat-row${
        isMine ? ' table-seat-row--mine' : taken ? ' table-seat-row--taken' : ' table-seat-row--open'
      }${seatLabel ? '' : ' table-seat-row--no-label'}`}
    >
      {seatLabel ? <span className="table-seat-row__label">{seatLabel}</span> : null}
      {showOccupant ? (
        <span className="table-seat-row__occupant" title={occupantLabel}>
          <PlayerAvatar user={occupantUser} size="sm" ring={isTableKing ? 'king' : undefined} />
          <span className="table-seat-row__occupant-name">{occupantLabel}</span>
        </span>
      ) : null}
      {!taken ? (
        <button
          type="button"
          className="table-seat-row__sit"
          disabled={busy}
          onClick={() => onSit(slot.seatKey)}
        >
          Sit
        </button>
      ) : null}
    </li>
  )
}

function SeatSection({ title, slots, mode, mySeat, userId, currentUser, kingUserId, busy, onSit }) {
  const queuePath = slots[0]?.queuePath
  const meta = queuePathMeta(mode, queuePath)
  const seatedCount = countSeatedInGroup(slots)
  const pooled = isPooledRoleGroup(slots)
  const mineInGroup = slots.some((slot) => slot.seatKey === mySeat || slot.user?.id === userId)
  const openSeatKey = firstOpenSeatKey(slots)

  if (pooled) {
    const occupants = slots.filter((slot) => slot.user || slot.seatKey === mySeat)
    return (
      <div className="table-card__team">
        <div className="table-card__team-heading">
          <h4 className="table-card__team-title">{title}</h4>
          <p className="table-card__team-caption">{formatGroupSeatCaption(seatedCount, meta)}</p>
        </div>
        <div className="table-card__pooled">
          <div className="table-card__pooled-avatars" aria-label={`${title} seats`}>
            {occupants.length > 0 ? (
              occupants.map((slot) => {
                const isMine = mySeat === slot.seatKey || slot.user?.id === userId
                const occupantUser = isMine ? currentUser ?? slot.user : slot.user
                const occupantLabel = isMine ? 'You' : displayName(slot.user)
                const isTableKing = Boolean(kingUserId && occupantUser?.id === kingUserId)
                const seatBadge = seatLabelInSection(slot, title)
                const titleText = seatBadge ? `${seatBadge} · ${occupantLabel}` : occupantLabel
                return (
                  <span
                    key={slot.seatKey}
                    className={`table-card__pooled-seat${isMine ? ' table-card__pooled-seat--mine' : ''}${
                      isTableKing ? ' table-card__pooled-seat--king' : ''
                    }`}
                    title={titleText}
                  >
                    {seatBadge ? <span className="table-card__pooled-seat-label">{seatBadge}</span> : null}
                    <PlayerAvatar user={occupantUser} size="sm" ring={isTableKing ? 'king' : undefined} />
                    <span className="table-card__pooled-seat-name">{occupantLabel}</span>
                  </span>
                )
              })
            ) : (
              <span className="table-card__pooled-empty">No one seated yet</span>
            )}
          </div>
          {!mineInGroup && openSeatKey ? (
            <button
              type="button"
              className="table-seat-row__sit"
              disabled={busy}
              onClick={() => onSit(openSeatKey)}
            >
              Sit
            </button>
          ) : null}
        </div>
      </div>
    )
  }

  return (
    <div className="table-card__team">
      <div className="table-card__team-heading">
        <h4 className="table-card__team-title">{title}</h4>
        {meta ? <p className="table-card__team-caption">{formatGroupSeatCaption(seatedCount, meta)}</p> : null}
      </div>
      <ul className="table-card__seats">
        {slots.map((slot) => (
          <SeatRow
            key={slot.seatKey}
            slot={slot}
            seatLabel={seatLabelInSection(slot, title)}
            mySeat={mySeat}
            userId={userId}
            currentUser={currentUser}
            kingUserId={kingUserId}
            busy={busy}
            onSit={onSit}
          />
        ))}
      </ul>
    </div>
  )
}

export default function TableCard({ table, busy, onSit, onLeave, onStart, onLookForGroup, onDiscard }) {
  const { user } = useAuth()
  const enriched = enrichTableSeats(table)
  const mySeat = mySeatKeyOnTable(enriched, user?.id)
  const mySeatLabel = mySeatDisplayName(enriched, user?.id)
  const king = isKing(enriched, user?.id)
  const kingUserId = enriched.king?.id
  const layout = groupSeatSlotsForDisplay(enriched.seatSlots ?? [])
  const seatedCount = enriched.seats?.length ?? 0
  const gapsLine = formatFormingGapsFromLobbyLine(enriched.formingGaps)

  const seatRowProps = {
    mode: enriched.mode,
    mySeat,
    userId: user?.id,
    currentUser: user,
    kingUserId,
    busy,
    onSit,
  }

  return (
    <article className="table-card">
      <header className="table-card__header">
        <div>
          <h3 className="table-card__title">
            {enriched.game?.name} · {enriched.mode?.displayName}
          </h3>
          <p className="table-card__meta">
            {seatedCount} seated
            {enriched.king ? ` · ${KING_LABEL}: ${displayName(enriched.king)}` : ''}
            {mySeatLabel ? ` · Your seat: ${mySeatLabel}` : ''}
          </p>
        </div>
      </header>

      <div
        className={`table-card__layout${
          layout.kind === 'teams' || layout.kind === 'roles' ? ' table-card__layout--teams' : ''
        }`}
      >
        {layout.kind === 'teams' ? (
          <>
            {layout.teams.map(([prefix, slots]) => (
              <SeatSection
                key={prefix}
                title={prefix.replace('-', ' ')}
                slots={slots}
                {...seatRowProps}
              />
            ))}
            {layout.ungrouped.length > 0 ? (
              <SeatSection title="Seats" slots={layout.ungrouped} {...seatRowProps} />
            ) : null}
          </>
        ) : null}

        {layout.kind === 'roles' ? (
          <>
            {layout.roles.map(([title, slots]) => (
              <SeatSection key={title} title={title} slots={slots} {...seatRowProps} />
            ))}
            {layout.noPath.length > 0 ? (
              <SeatSection title="Seats" slots={layout.noPath} {...seatRowProps} />
            ) : null}
          </>
        ) : null}

        {layout.kind === 'flat' ? (
          isPooledRoleGroup(layout.slots) ? (
            <SeatSection title="Players" slots={layout.slots} {...seatRowProps} />
          ) : (
            <ul className="table-card__seats table-card__seats--flat">
              {layout.slots.map((slot) => (
                <SeatRow
                  key={slot.seatKey}
                  slot={slot}
                  seatLabel={seatLabelInSection(slot, null)}
                  {...seatRowProps}
                />
              ))}
            </ul>
          )
        ) : null}
      </div>

      {gapsLine ? <p className="table-card__gaps">{gapsLine}</p> : null}

      <div className="table-card__actions">
        {mySeat ? (
          <button type="button" className="game-list-button game-list-button-secondary" disabled={busy} onClick={onLeave}>
            Leave seat
          </button>
        ) : null}
        {king && enriched.canStart ? (
          <button type="button" className="game-list-button" disabled={busy} onClick={onStart}>
            {START_GAME}
          </button>
        ) : null}
        {king
          ? enriched.lookForGroupOptions
              ?.filter((opt) => opt.visible)
              .map((opt) => (
                <button
                  key={opt.queueId}
                  type="button"
                  className="game-list-button"
                  disabled={busy || !opt.enabled || enriched.backfillActive}
                  onClick={() => onLookForGroup?.(opt.queueId)}
                >
                  {LOOK_FOR_GROUP} ({opt.queueName})
                </button>
              ))
          : null}
        {enriched.canDiscard ? (
          <button
            type="button"
            className="game-list-button game-list-button-secondary"
            disabled={busy}
            onClick={onDiscard}
          >
            {DISCARD}
          </button>
        ) : null}
      </div>
    </article>
  )
}
