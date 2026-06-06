import { useAuth } from '../auth/AuthProvider'
import {
  DISCARD,
  KING_LABEL,
  LOOK_FOR_GROUP,
  START_GAME,
} from '../../lib/playerCopy'
import {
  displayName,
  enrichTableSeats,
  groupSeatSlotsForDisplay,
  isKing,
  mySeatDisplayName,
  mySeatKeyOnTable,
  seatLabelInSection,
} from '../../lib/tables'
import PlayerAvatar from '../avatars/PlayerAvatar'

function SeatRow({ slot, seatLabel, mySeat, userId, currentUser, busy, onSit }) {
  const taken = Boolean(slot.user)
  const isMine = mySeat === slot.seatKey || slot.user?.id === userId
  const showOccupant = taken || isMine
  const occupantUser = isMine ? currentUser ?? slot.user : slot.user
  const occupantLabel = isMine ? 'You' : displayName(slot.user)

  return (
    <li
      className={`table-seat-row${
        isMine ? ' table-seat-row--mine' : taken ? ' table-seat-row--taken' : ' table-seat-row--open'
      }${seatLabel ? '' : ' table-seat-row--no-label'}`}
    >
      {seatLabel ? <span className="table-seat-row__label">{seatLabel}</span> : null}
      {showOccupant ? (
        <span className="table-seat-row__occupant" title={occupantLabel}>
          <PlayerAvatar user={occupantUser} size="sm" />
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

function SeatSection({ title, slots, mySeat, userId, currentUser, busy, onSit }) {
  return (
    <div className="table-card__team">
      <h4 className="table-card__team-title">{title}</h4>
      <ul className="table-card__seats">
        {slots.map((slot) => (
          <SeatRow
            key={slot.seatKey}
            slot={slot}
            seatLabel={seatLabelInSection(slot, title)}
            mySeat={mySeat}
            userId={userId}
            currentUser={currentUser}
            busy={busy}
            onSit={onSit}
          />
        ))}
      </ul>
    </div>
  )
}

export default function TableCard({ table, busy, onSit, onLeave, onStart, onDiscard }) {
  const { user } = useAuth()
  const enriched = enrichTableSeats(table)
  const mySeat = mySeatKeyOnTable(enriched, user?.id)
  const mySeatLabel = mySeatDisplayName(enriched, user?.id)
  const king = isKing(enriched, user?.id)
  const layout = groupSeatSlotsForDisplay(enriched.seatSlots ?? [])
  const seatedCount = enriched.seats?.length ?? 0

  const seatRowProps = { mySeat, userId: user?.id, currentUser: user, busy, onSit }

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
        ) : null}
      </div>

      {(enriched.lookForGroupOptions ?? []).some((opt) => opt.visible) ? (
        <div className="table-card__lfg">
          {enriched.lookForGroupOptions
            .filter((opt) => opt.visible)
            .map((opt) => (
              <button key={opt.queueId} type="button" className="game-list-button" disabled>
                {LOOK_FOR_GROUP} ({opt.queueName})
              </button>
            ))}
        </div>
      ) : null}

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
