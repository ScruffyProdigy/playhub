export default function GameListItem({ game }) {
  const sessionCount = game.activeSessions?.length ?? 0

  return (
    <li className="game-list-item">
      <div className="game-list-copy">
        <h3>{game.name}</h3>
        <p className="game-list-meta">
          {sessionCount} active session{sessionCount === 1 ? '' : 's'}
        </p>
      </div>
    </li>
  )
}
