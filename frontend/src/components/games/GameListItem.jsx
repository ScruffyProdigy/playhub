import { useCallback, useEffect, useState } from 'react'
import { joinGroupOptionsForMode } from '../../lib/games'
import { hasPlayingIntent, hasWaitingIntent } from '../../lib/intent'
import { CREATE_PRIVATE_GAME } from '../../lib/playerCopy'
import { createPrivateTable } from '../../lib/tables'
import { useActiveRoom } from '../rooms/ActiveRoomProvider'
import GameQueueActions from './GameQueueActions'
import { useGameQueue } from './useGameQueue'

function ModeRow({
  game,
  mode,
  activeIntent,
  activeTableSeat,
  onQueueChange,
  onQueueJoined,
  onTableChange,
}) {
  const { refresh: refreshRoom, openRoom } = useActiveRoom()
  const defaultQueue = mode.queues?.find((q) => q.status === 'active') ?? null
  const joinOptions = joinGroupOptionsForMode(mode)
  const [tableBusy, setTableBusy] = useState(false)
  const [tableError, setTableError] = useState('')

  const isThisQueue =
    defaultQueue?.id && activeIntent?.queueId && activeIntent.queueId === defaultQueue.id
  const queue = useGameQueue(defaultQueue?.id, {
    skipSubscription: Boolean(isThisQueue && activeIntent),
  })

  const seatedHere =
    activeTableSeat?.gameId === game.id && activeTableSeat?.modeId === mode.id
  const inActiveGame = activeIntent?.status === 'MATCHED'

  async function handleJoin(queuePath) {
    if (inActiveGame) {
      setTableError('Launch or finish your current game before looking for a group.')
      return
    }
    if (activeTableSeat?.tableId) {
      setTableError('Leave your table seat before looking for a group.')
      return
    }
    const pathLabel =
      joinOptions?.kind === 'composition'
        ? joinOptions.paths
            .map((path) => (typeof path === 'string' ? { queuePath: path, displayName: path } : path))
            .find((entry) => entry.queuePath === queuePath)?.displayName
        : null
    const result = await queue.handleJoin(queuePath)
    if (result) {
      onQueueJoined?.(defaultQueue?.id, result, {
        gameId: game.id,
        gameName: game.name,
        modeName: mode.displayName,
        queuePathDisplayName: pathLabel ?? null,
      })
    }
  }

  async function handleLeave() {
    await queue.handleLeave()
    await onQueueChange?.()
  }

  async function handleCreatePrivate() {
    if (inActiveGame) {
      setTableError('Launch or finish your current game before creating a private game.')
      return
    }
    if (activeTableSeat?.tableId && !seatedHere) {
      setTableError('Leave your current table seat first.')
      return
    }
    if (activeIntent?.queueId) {
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
    inActiveGame ||
    (activeIntent &&
      defaultQueue?.id &&
      !isThisQueue &&
      queue.queueState === 'idle' &&
      activeIntent.status === 'MATCHED')

  const resolvedQueueState = isThisQueue
    ? hasWaitingIntent(activeIntent)
      ? 'waiting'
      : hasPlayingIntent(activeIntent)
        ? 'matched'
        : queue.queueState
    : queue.queueState
  const resolvedJoinUrl =
    isThisQueue && activeIntent?.joinUrl ? activeIntent.joinUrl : queue.joinUrl

  return (
    <li className="game-mode-row">
      <div className="game-mode-row__copy">
        <h4 className="game-mode-row__title">{mode.displayName}</h4>
        {tableError ? (
          <p className="status-message status-message-error" role="status">
            {tableError}
          </p>
        ) : null}
        {queue.error ? (
          <p className="status-message status-message-error" role="status">
            {queue.error}
          </p>
        ) : null}
        {queue.notice ? (
          <p className="status-message" role="status">
            {queue.notice}
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
          queueState={resolvedQueueState}
          joinUrl={resolvedJoinUrl}
          busy={queue.busy}
          selectedQueuePath={
            queue.selectedQueuePath || (isThisQueue ? activeIntent?.queuePath : '') || ''
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

export default function GameListItem({ game, activeIntent, activeTableSeat, onQueueChange, onQueueJoined, onTableChange }) {
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
              activeIntent={activeIntent}
              activeTableSeat={activeTableSeat}
              onQueueChange={onQueueChange}
              onQueueJoined={onQueueJoined}
              onTableChange={onTableChange}
            />
          ))}
        </ul>
      )}
    </li>
  )
}
