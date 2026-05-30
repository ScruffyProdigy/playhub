import GameListItemInfo from './GameListItemInfo'
import GameQueueActions from './GameQueueActions'
import { useGameQueue } from './useGameQueue'

export default function GameListItem({ game }) {
  const queue = useGameQueue(game.id)

  return (
    <li className="game-list-item">
      <GameListItemInfo
        game={game}
        queueState={queue.queueState}
        queuedCount={queue.queuedCount}
        error={queue.error}
      />
      <GameQueueActions
        queueState={queue.queueState}
        joinUrl={queue.joinUrl}
        busy={queue.busy}
        onJoin={queue.handleJoin}
        onLeave={queue.handleLeave}
      />
    </li>
  )
}
