import { useCallback, useEffect, useState } from 'react'
import { joinGroupOptionsForMode } from '../../lib/games'
import { CREATE_PRIVATE_GAME } from '../../lib/playerCopy'
import { createPrivateTable } from '../../lib/tables'
import { useActiveRoom } from '../rooms/ActiveRoomProvider'
import GameQueueActions from './GameQueueActions'
import { useGameQueue } from './useGameQueue'

function ModeRow({
  game,
  mode,
  activeQueue,
  activeTableSeat,
  onQueueChange,
  onTableChange,
}) {
  const { refresh: refreshRoom, openRoom } = useActiveRoom()
  const defaultQueue = mode.queues?.find((q) => q.status === 'active') ?? null
  const joinOptions = joinGroupOptionsForMode(mode)
  const [tableBusy, setTableBusy] = useState(false)
  const [tableError, setTableError] = useState('')

  const isThisQueue =
    defaultQueue?.id && activeQueue?.queueId && activeQueue.queueId === defaultQueue.id
  const queue = useGameQueue(defaultQueue?.id, {
    skipSubscription: Boolean(isThisQueue && activeQueue),
  })

  const seatedHere =
    activeTableSeat?.gameId === game.id && activeTableSeat?.modeId === mode.id

  async function handleJoin(queuePath) {
    if (activeTableSeat?.tableId) {
      setTableError('Leave your table seat before looking for a group.')
      return
    }
    await queue.handleJoin(queuePath)
    await onQueueChange?.()
  }

  async function handleLeave() {
    await queue.handleLeave()
    await onQueueChange?.()
  }

  async function handleCreatePrivate() {
    if (activeTableSeat?.tableId && !seatedHere) {
      setTableError('Leave your current table seat first.')
      return
    }
    if (activeQueue?.queueId) {
      setTableError('Stop looking for a group before creating a private game.')
      return
    }
    setTableBusy(true)
    setTableError('')
    try {
      await createPrivateTable(game.id, mode.id)
      await refreshRoom()
      openRoom()
      await onTableChange?.()
    } catch (err) {
      setTableError(err.message || 'Could not create private game.')
    } finally {
      setTableBusy(false)
    }
  }

  const blockedByMatch =
    activeQueue &&
    defaultQueue?.id &&
    !isThisQueue &&
    queue.queueState === 'idle' &&
    activeQueue.status === 'MATCHED'

  return (
    <li className="game-mode-row">
      <div className="game-mode-row__copy">
        <h4 className="game-mode-row__title">{mode.displayName}</h4>
        {tableError ? (
          <p className="status-message status-message-error" role="status">
            {tableError}
          </p>
        ) : null}
        {seatedHere ? (
          <p className="game-list-meta" role="status">
            You are seated at a private table for this mode.
          </p>
        ) : null}
      </div>
      <div className="game-mode-row__actions">
        <GameQueueActions
          joinOptions={joinOptions}
          queueState={queue.queueState}
          joinUrl={queue.joinUrl}
          busy={queue.busy}
          selectedQueuePath={
            queue.selectedQueuePath || (isThisQueue ? activeQueue?.queuePath : '') || ''
          }
          onJoin={handleJoin}
          onLeave={handleLeave}
          disabled={!defaultQueue || blockedByMatch || Boolean(activeTableSeat?.tableId && !seatedHere)}
        />
        <button
          type="button"
          className="game-list-button game-list-button-secondary"
          disabled={tableBusy || blockedByMatch}
          onClick={handleCreatePrivate}
        >
          {tableBusy ? '…' : CREATE_PRIVATE_GAME}
        </button>
      </div>
    </li>
  )
}

export default function GameListItem({ game, activeQueue, activeTableSeat, onQueueChange, onTableChange }) {
  const modes = (game.modes ?? []).filter((mode) => mode.status === 'active')

  return (
    <li className="game-list-item game-list-item--modes">
      <div className="game-list-copy">
        <h3>{game.name}</h3>
        <p className="game-list-meta">
          {game.activeSessions?.length ?? 0} active session
          {(game.activeSessions?.length ?? 0) === 1 ? '' : 's'}
        </p>
      </div>
      {modes.length === 0 ? (
        <p className="game-list-meta">No active modes.</p>
      ) : (
        <ul className="game-mode-list">
          {modes.map((mode) => (
            <ModeRow
              key={mode.id ?? mode.modeKey}
              game={game}
              mode={mode}
              activeQueue={activeQueue}
              activeTableSeat={activeTableSeat}
              onQueueChange={onQueueChange}
              onTableChange={onTableChange}
            />
          ))}
        </ul>
      )}
    </li>
  )
}
