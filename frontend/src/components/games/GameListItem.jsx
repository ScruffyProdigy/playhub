import { gameCatalogHeroUrl, gamePagePath, gameTagChips } from '../../lib/gameCard'
import CatalogGameLink from './CatalogGameLink'
import GameModesPanel from './GameModesPanel'

export default function GameListItem({
  game,
  activeIntent,
  activeTableSeat,
  onQueueChange,
  onQueueJoined,
  onTableChange,
}) {
  const tags = gameTagChips(game.tags)
  const detailPath = gamePagePath(game)

  return (
    <li className="game-list-item game-list-item--hero">
      <div className="game-list-item__hero-wrap">
        {detailPath ? (
          <CatalogGameLink slug={game.slug} className="game-list-item__hero-link" href={detailPath}>
            <img
              className="game-list-item__hero"
              src={gameCatalogHeroUrl(game)}
              alt=""
              width={400}
              height={160}
              loading="lazy"
            />
          </CatalogGameLink>
        ) : (
          <img
            className="game-list-item__hero"
            src={gameCatalogHeroUrl(game)}
            alt=""
            width={400}
            height={160}
            loading="lazy"
          />
        )}
      </div>
      <div className="game-list-item__body">
        {detailPath ? (
          <h3>
            <CatalogGameLink slug={game.slug} className="game-list-item__title-link" href={detailPath}>
              {game.name}
            </CatalogGameLink>
          </h3>
        ) : (
          <h3>{game.name}</h3>
        )}
        {tags.length > 0 ? (
          <ul className="game-list-item__tags" aria-label="Game tags">
            {tags.map((tag) => (
              <li key={tag} className="game-list-item__tag">
                {tag}
              </li>
            ))}
          </ul>
        ) : null}
        {game.shortDescription ? (
          <p className="game-list-item__description">{game.shortDescription}</p>
        ) : null}
        <GameModesPanel
          game={game}
          activeIntent={activeIntent}
          activeTableSeat={activeTableSeat}
          onQueueChange={onQueueChange}
          onQueueJoined={onQueueJoined}
          onTableChange={onTableChange}
          heading="Modes"
        />
      </div>
    </li>
  )
}
