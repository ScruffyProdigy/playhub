import { IconCatalog, IconMessage } from '../icons/ShareIcons'

export default function AppDock({ unreadCount, roomOpen, onOpenCatalog, onOpenRoom }) {
  return (
    <nav className="app-dock" aria-label="Main navigation">
      <button
        type="button"
        className={`app-dock__tab ${!roomOpen ? 'app-dock__tab--active' : ''}`}
        onClick={onOpenCatalog}
        aria-current={!roomOpen ? 'page' : undefined}
      >
        <IconCatalog />
        <span>Catalog</span>
      </button>
      <button
        type="button"
        className={`app-dock__tab ${roomOpen ? 'app-dock__tab--active' : ''}`}
        onClick={onOpenRoom}
        aria-current={roomOpen ? 'page' : undefined}
      >
        <span className="app-dock__icon-wrap">
          <IconMessage />
          {unreadCount > 0 ? (
            <span className="app-dock__badge" aria-label={`${unreadCount} unread messages`}>
              {unreadCount > 99 ? '99+' : unreadCount}
            </span>
          ) : null}
        </span>
        <span>Room</span>
      </button>
    </nav>
  )
}
