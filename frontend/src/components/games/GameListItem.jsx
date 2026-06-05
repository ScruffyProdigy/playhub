import { defaultQueueForGame, joinGroupOptionsForGame } from '../../lib/games'
import { activeMatchBlockedLine } from '../../lib/playerCopy'
import GameListItemInfo from './GameListItemInfo'
import GameQueueActions from './GameQueueActions'
import { useGameQueue } from './useGameQueue'

export default function GameListItem({ game, activeQueue, onQueueChange }) {
  const defaultQueue = defaultQueueForGame(game)
  const joinOptions = joinGroupOptionsForGame(game)

  const isThisQueue =
    defaultQueue?.id && activeQueue?.queueId && activeQueue.queueId === defaultQueue.id
  const queue = useGameQueue(defaultQueue?.id, {
    skipSubscription: Boolean(isThisQueue && activeQueue),
  })
  const showRowActions = !isThisQueue || queue.queueState !== 'matched'
  const blockedActiveMatch =
    activeQueue &&
    defaultQueue?.id &&
    !isThisQueue &&
    queue.queueState === 'idle' &&
    activeQueue.status === 'MATCHED'

  async function handleJoin(queuePath) {
    await queue.handleJoin(queuePath)
    await onQueueChange?.()
  }

  async function handleLeave() {
    await queue.handleLeave()
    await onQueueChange?.()
  }

  const displayError =
    queue.error ||
    (blockedActiveMatch ? activeMatchBlockedLine(activeQueue.gameName) : '')

  return (
    <li className="game-list-item">
      <GameListItemInfo
        game={game}
        queueState={queue.queueState}
        queuedCount={queue.queuedCount}
        error={displayError}
        notice={queue.notice}
      />
      {showRowActions ? (
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
          disabled={!defaultQueue || blockedActiveMatch}
        />
      ) : (
        <p className="game-list-meta game-list-meta--banner">Use the banner above to launch or stop looking.</p>
      )}
    </li>
  )
}
